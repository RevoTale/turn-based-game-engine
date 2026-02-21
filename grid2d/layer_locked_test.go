package grid2d

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLockedSparseLayer_BasicOperations(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(3, 3)
	require.NoError(t, err)

	layer, err := NewLockedSparseLayer[int](grid)
	require.NoError(t, err)

	err = layer.Set(Position{X: 1, Y: 1}, 7)
	require.NoError(t, err)

	value, ok, err := layer.Get(Position{X: 1, Y: 1})
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, 7, value)

	err = layer.Delete(Position{X: 1, Y: 1})
	require.NoError(t, err)
	assert.Equal(t, 0, layer.Len())
}

func TestRowLockedSparseLayer_BasicOperations(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(4, 4)
	require.NoError(t, err)

	layer, err := NewRowLockedSparseLayer[string](grid)
	require.NoError(t, err)

	err = layer.Set(Position{X: 2, Y: 3}, "hit")
	require.NoError(t, err)

	value, ok, err := layer.Get(Position{X: 2, Y: 3})
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "hit", value)
	assert.Equal(t, 1, layer.Len())

	err = layer.Delete(Position{X: 2, Y: 3})
	require.NoError(t, err)
	assert.Equal(t, 0, layer.Len())
}

func TestRowLockedSparseLayer_ConcurrentRowWrites(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(16, 16)
	require.NoError(t, err)

	layer, err := NewRowLockedSparseLayer[int](grid)
	require.NoError(t, err)

	var wg sync.WaitGroup
	for y := 0; y < grid.Height(); y++ {
		y := y
		wg.Add(1)
		go func() {
			defer wg.Done()
			for x := 0; x < grid.Width(); x++ {
				setErr := layer.Set(Position{X: x, Y: y}, y*100+x)
				require.NoError(t, setErr)
			}
		}()
	}
	wg.Wait()

	assert.Equal(t, grid.CellCount(), layer.Len())

	for y := 0; y < grid.Height(); y++ {
		for x := 0; x < grid.Width(); x++ {
			value, ok, getErr := layer.Get(Position{X: x, Y: y})
			require.NoError(t, getErr)
			require.True(t, ok)
			assert.Equal(t, y*100+x, value)
		}
	}
}

func TestLockedSparseLayer_Bounds(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(2, 2)
	require.NoError(t, err)

	layer, err := NewLockedSparseLayer[bool](grid)
	require.NoError(t, err)

	err = layer.Set(Position{X: 5, Y: 0}, true)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrOutOfBounds)
}

func TestRowLockedSparseLayer_Bounds(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(2, 2)
	require.NoError(t, err)

	layer, err := NewRowLockedSparseLayer[bool](grid)
	require.NoError(t, err)

	err = layer.Set(Position{X: 0, Y: 3}, true)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrOutOfBounds)
}
