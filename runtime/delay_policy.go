package engruntime

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// RandomDelayPolicyConfig configures random delay generation bounds.
type RandomDelayPolicyConfig struct {
	// Min is the minimum delay.
	Min time.Duration
	// Max is the maximum delay.
	Max time.Duration
	// Now supplies current time used for RNG seeding.
	Now func() time.Time
}

// RandomDelayPolicy generates pseudo-random delays per key.
//
// Each key gets its own RNG instance to reduce contention across keys while
// remaining safe for concurrent use.
type RandomDelayPolicy[K comparable] struct {
	min   time.Duration
	max   time.Duration
	now   func() time.Time
	seed  uint64
	rngBy sync.Map // map[K]*keyedRNG
}

type keyedRNG struct {
	mu  sync.Mutex
	rng *rand.Rand
}

// NewRandomDelayPolicy creates a key-scoped random delay policy.
func NewRandomDelayPolicy[K comparable](cfg RandomDelayPolicyConfig) *RandomDelayPolicy[K] {
	min := cfg.Min
	max := cfg.Max
	if min < 0 {
		min = 0
	}
	if max < 0 {
		max = 0
	}
	if max < min {
		max = min
	}

	now := cfg.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}

	return &RandomDelayPolicy[K]{
		min: min,
		max: max,
		now: now,
	}
}

// Next returns the next delay for key, respecting configured bounds.
func (p *RandomDelayPolicy[K]) Next(key K) time.Duration {
	if p == nil {
		return 0
	}

	if p.max <= p.min {
		return p.min
	}

	entry := p.rngForKey(key)
	entry.mu.Lock()
	defer entry.mu.Unlock()

	span := int64(p.max - p.min)
	return p.min + inclusiveOffset(entry.rng, span)
}

// Bounds returns normalized minimum and maximum delay bounds.
func (p *RandomDelayPolicy[K]) Bounds() (time.Duration, time.Duration) {
	if p == nil {
		return 0, 0
	}
	return p.min, p.max
}

func (p *RandomDelayPolicy[K]) rngForKey(key K) *keyedRNG {
	if existing, ok := p.rngBy.Load(key); ok {
		return existing.(*keyedRNG)
	}

	seed := p.now().UnixNano() + int64(atomic.AddUint64(&p.seed, 1))
	entry := &keyedRNG{
		rng: rand.New(rand.NewSource(seed)),
	}
	actual, _ := p.rngBy.LoadOrStore(key, entry)
	return actual.(*keyedRNG)
}

func inclusiveOffset(rng *rand.Rand, span int64) time.Duration {
	if rng == nil || span <= 0 {
		return 0
	}
	if span == math.MaxInt64 {
		return time.Duration(rng.Int63n(span))
	}
	return time.Duration(rng.Int63n(span + 1))
}
