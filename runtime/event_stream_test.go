package engruntime

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTopicStream_PublishDropsSlowSubscriber(t *testing.T) {
	t.Parallel()

	stream := NewTopicStream[string, string, string]()
	topic := "room-1"

	slowCtx, cancelSlow := context.WithCancel(context.Background())
	defer cancelSlow()
	_ = stream.Subscribe(slowCtx, topic, "slow", SubscribeConfig{Buffer: 0})

	fastCtx, cancelFast := context.WithCancel(context.Background())
	defer cancelFast()
	fast := stream.Subscribe(fastCtx, topic, "fast", SubscribeConfig{Buffer: 1})

	result := stream.Publish(topic, "hello", PublishConfig{
		DeliveryTimeout:   0,
		RemoveUndelivered: true,
	})

	assert.Equal(t, 2, result.Subscribers)
	assert.Equal(t, 1, result.Delivered)
	assert.Equal(t, 1, result.Removed)
	assert.Equal(t, 1, stream.TopicSubscriberCount(topic))

	select {
	case got := <-fast:
		assert.Equal(t, "hello", got)
	default:
		t.Fatal("expected fast subscriber to receive event")
	}
}

func TestTopicStream_SubscribeReplacesExistingSubscriberID(t *testing.T) {
	t.Parallel()

	stream := NewTopicStream[string, string, int]()
	topic := "room-2"

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	oldCh := stream.Subscribe(ctx1, topic, "same-id", SubscribeConfig{Buffer: 1})

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	newCh := stream.Subscribe(ctx2, topic, "same-id", SubscribeConfig{Buffer: 1})

	assert.Equal(t, 1, stream.TopicSubscriberCount(topic))

	_, oldOpen := <-oldCh
	assert.False(t, oldOpen)

	result := stream.Publish(topic, 10, PublishConfig{RemoveUndelivered: true})
	assert.Equal(t, 1, result.Delivered)

	select {
	case got := <-newCh:
		assert.Equal(t, 10, got)
	default:
		t.Fatal("expected event on new subscription")
	}
}

func TestTopicStream_ContextCancelClosesSubscription(t *testing.T) {
	t.Parallel()

	stream := NewTopicStream[string, string, string]()
	topic := "room-3"

	ctx, cancel := context.WithCancel(context.Background())
	ch := stream.Subscribe(ctx, topic, "client-1", SubscribeConfig{Buffer: 1})
	require.NotNil(t, ch)
	assert.Equal(t, 1, stream.TopicSubscriberCount(topic))

	cancel()

	require.Eventually(t, func() bool {
		return stream.TopicSubscriberCount(topic) == 0
	}, 200*time.Millisecond, 10*time.Millisecond)

	_, open := <-ch
	assert.False(t, open)
}

func TestTopicStream_PublishTimeoutDropsStalledSubscriber(t *testing.T) {
	t.Parallel()

	stream := NewTopicStream[string, string, int]()
	topic := "room-4"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = stream.Subscribe(ctx, topic, "slow", SubscribeConfig{Buffer: 0})

	result := stream.Publish(topic, 10, PublishConfig{
		DeliveryTimeout:   5 * time.Millisecond,
		RemoveUndelivered: true,
	})

	assert.Equal(t, 1, result.Subscribers)
	assert.Equal(t, 0, result.Delivered)
	assert.Equal(t, 1, result.Removed)
	assert.Equal(t, 0, stream.TopicSubscriberCount(topic))
}
