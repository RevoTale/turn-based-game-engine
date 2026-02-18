package engruntime

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomDelayPolicy_BoundsNormalization(t *testing.T) {
	t.Parallel()

	policy := NewRandomDelayPolicy[string](RandomDelayPolicyConfig{
		Min: -2 * time.Second,
		Max: -1 * time.Second,
	})
	min, max := policy.Bounds()
	assert.Equal(t, time.Duration(0), min)
	assert.Equal(t, time.Duration(0), max)

	policy = NewRandomDelayPolicy[string](RandomDelayPolicyConfig{
		Min: 3 * time.Second,
		Max: time.Second,
	})
	min, max = policy.Bounds()
	assert.Equal(t, 3*time.Second, min)
	assert.Equal(t, 3*time.Second, max)
}

func TestRandomDelayPolicy_NextWithinBounds(t *testing.T) {
	t.Parallel()

	policy := NewRandomDelayPolicy[string](RandomDelayPolicyConfig{
		Min: time.Millisecond,
		Max: 3 * time.Millisecond,
	})

	min, max := policy.Bounds()
	for i := 0; i < 1000; i++ {
		d := policy.Next("room-1")
		assert.GreaterOrEqual(t, d, min)
		assert.LessOrEqual(t, d, max)
	}
}

func TestRandomDelayPolicy_ConstantBounds(t *testing.T) {
	t.Parallel()

	policy := NewRandomDelayPolicy[int](RandomDelayPolicyConfig{
		Min: 2 * time.Second,
		Max: 2 * time.Second,
	})
	assert.Equal(t, 2*time.Second, policy.Next(1))
	assert.Equal(t, 2*time.Second, policy.Next(2))
}

func TestRandomDelayPolicy_NilPolicy(t *testing.T) {
	t.Parallel()

	var policy *RandomDelayPolicy[string]
	assert.Equal(t, time.Duration(0), policy.Next("key"))
	min, max := policy.Bounds()
	assert.Equal(t, time.Duration(0), min)
	assert.Equal(t, time.Duration(0), max)
}

func TestRandomDelayPolicy_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	policy := NewRandomDelayPolicy[int](RandomDelayPolicyConfig{
		Min: time.Millisecond,
		Max: 5 * time.Millisecond,
	})
	min, max := policy.Bounds()

	var wg sync.WaitGroup
	for worker := 0; worker < 32; worker++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				d := policy.Next(id)
				require.GreaterOrEqual(t, d, min)
				require.LessOrEqual(t, d, max)
			}
		}(worker)
	}
	wg.Wait()
}
