package engruntime

import (
	"context"
	"sync"
	"time"
)

// controls subscriber behavior.
type SubscribeConfig struct {
	// Buffer sets channel capacity.
	// - 0 means unbuffered
	// - values < 0 are normalized to 1
	Buffer int
}

// controls stream delivery policy.
type PublishConfig struct {
	// DeliveryTimeout bounds per-subscriber delivery wait.
	// 0 means non-blocking send attempt.
	DeliveryTimeout time.Duration
	// RemoveUndelivered removes subscribers that failed delivery (slow/closed).
	RemoveUndelivered bool
}

// summarizes one publish attempt.
type PublishResult struct {
	Subscribers int
	Delivered   int
	Removed     int
}

type streamTarget[S comparable, E any] struct {
	subscriber S
	ch         chan E
}

// TopicStream is a typed, topic-based pub/sub stream keyed by subscriber id.
//
// Delivery semantics are at-most-once:
// - no retries
// - no reconnect replay
// - slow/closed subscribers can be dropped (policy-driven)
type TopicStream[K comparable, S comparable, E any] struct {
	mu     sync.RWMutex
	topics map[K]map[S]chan E
}

func NewTopicStream[K comparable, S comparable, E any]() *TopicStream[K, S, E] {
	return &TopicStream[K, S, E]{
		topics: make(map[K]map[S]chan E),
	}
}

// Subscribe registers or replaces subscriber in a topic and auto-cleans it on ctx.Done.
func (s *TopicStream[K, S, E]) Subscribe(ctx context.Context, topic K, subscriber S, cfg SubscribeConfig) <-chan E {
	if s == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	buffer := cfg.Buffer
	if buffer < 0 {
		buffer = 1
	}

	ch := make(chan E, buffer)

	s.mu.Lock()
	topicSubs, ok := s.topics[topic]
	if !ok {
		topicSubs = make(map[S]chan E)
		s.topics[topic] = topicSubs
	}

	if old, exists := topicSubs[subscriber]; exists {
		closeNoPanic(old)
	}
	topicSubs[subscriber] = ch
	s.mu.Unlock()

	go func() {
		<-ctx.Done()
		s.removeSubscriber(topic, subscriber, true)
	}()

	return ch
}

// Unsubscribe removes and closes a subscriber channel.
func (s *TopicStream[K, S, E]) Unsubscribe(topic K, subscriber S) bool {
	return s.removeSubscriber(topic, subscriber, true)
}

// Publish delivers event to all subscribers in topic according to cfg.
func (s *TopicStream[K, S, E]) Publish(topic K, event E, cfg PublishConfig) PublishResult {
	if s == nil {
		return PublishResult{}
	}

	targets := s.snapshot(topic)
	if len(targets) == 0 {
		return PublishResult{}
	}

	result := PublishResult{
		Subscribers: len(targets),
	}

	var toRemove []S
	for _, target := range targets {
		delivered, remove := deliverWithPolicy(target.ch, event, cfg)
		if delivered {
			result.Delivered++
			continue
		}
		if remove {
			toRemove = append(toRemove, target.subscriber)
		}
	}

	for _, subscriber := range toRemove {
		if s.removeSubscriber(topic, subscriber, true) {
			result.Removed++
		}
	}

	return result
}

// TopicSubscriberCount returns number of active subscribers for a topic.
func (s *TopicStream[K, S, E]) TopicSubscriberCount(topic K) int {
	if s == nil {
		return 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.topics[topic])
}

func (s *TopicStream[K, S, E]) snapshot(topic K) []streamTarget[S, E] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	topicSubs, ok := s.topics[topic]
	if !ok || len(topicSubs) == 0 {
		return nil
	}

	out := make([]streamTarget[S, E], 0, len(topicSubs))
	for subscriber, ch := range topicSubs {
		out = append(out, streamTarget[S, E]{
			subscriber: subscriber,
			ch:         ch,
		})
	}
	return out
}

func (s *TopicStream[K, S, E]) removeSubscriber(topic K, subscriber S, closeChannel bool) bool {
	if s == nil {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	topicSubs, ok := s.topics[topic]
	if !ok {
		return false
	}

	ch, ok := topicSubs[subscriber]
	if !ok {
		return false
	}
	delete(topicSubs, subscriber)
	if len(topicSubs) == 0 {
		delete(s.topics, topic)
	}

	if closeChannel {
		closeNoPanic(ch)
	}
	return true
}

func deliverWithPolicy[E any](ch chan E, event E, cfg PublishConfig) (delivered bool, remove bool) {
	if ch == nil {
		return false, cfg.RemoveUndelivered
	}

	// Closed channels should be treated as removable subscribers.
	defer func() {
		if r := recover(); r != nil {
			delivered = false
			remove = true
		}
	}()

	if cfg.DeliveryTimeout <= 0 {
		select {
		case ch <- event:
			return true, false
		default:
			return false, cfg.RemoveUndelivered
		}
	}

	timer := time.NewTimer(cfg.DeliveryTimeout)
	defer func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}()

	select {
	case ch <- event:
		return true, false
	case <-timer.C:
		return false, cfg.RemoveUndelivered
	}
}

func closeNoPanic[E any](ch chan E) {
	if ch == nil {
		return
	}
	defer func() {
		_ = recover()
	}()
	close(ch)
}
