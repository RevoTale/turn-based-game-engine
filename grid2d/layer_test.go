package grid2d

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSparseLayer_LazyAllocationAndStringSupport(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(3, 2)
	require.NoError(t, err)

	layer, err := NewSparseLayer[string](grid)
	require.NoError(t, err)

	assert.Nil(t, layer.values)
	assert.Equal(t, 0, layer.Len())

	err = layer.Set(Position{X: 1, Y: 0}, "burning")
	require.NoError(t, err)

	assert.NotNil(t, layer.values)
	assert.Equal(t, 1, layer.Len())

	value, ok, err := layer.Get(Position{X: 1, Y: 0})
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "burning", value)
}

func TestSparseLayer_BoundsAndDelete(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(2, 2)
	require.NoError(t, err)

	layer, err := NewSparseLayer[bool](grid)
	require.NoError(t, err)

	err = layer.Set(Position{X: 2, Y: 0}, true)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrOutOfBounds)

	err = layer.Set(Position{X: 1, Y: 1}, true)
	require.NoError(t, err)

	err = layer.Delete(Position{X: 1, Y: 1})
	require.NoError(t, err)
	assert.Equal(t, 0, layer.Len())
}
