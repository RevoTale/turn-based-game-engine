package multi

import (
	"testing"

	"github.com/RevoTale/turn-based-game-engine/grid2d"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testBoardGridID   uint8 = 1
	testMissingGridID uint8 = 99
	testHitsLayerKey  uint8 = 1
	testStateLayerKey uint8 = 2
)

func TestRegistry_LayerOperations(t *testing.T) {
	t.Parallel()

	registry := NewRegistry[uint8, uint8, string]()
	_, err := registry.CreateGrid(testBoardGridID, 3, 2)
	require.NoError(t, err)
	assert.Equal(t, 0, registry.LayerCount(testBoardGridID))

	layer, err := registry.Layer(testBoardGridID, testStateLayerKey)
	require.NoError(t, err)
	assert.Equal(t, 1, registry.LayerCount(testBoardGridID))

	require.NoError(t, layer.Set(grid2d.Position{X: 2, Y: 1}, "occupied"))
	value, ok, err := layer.Get(grid2d.Position{X: 2, Y: 1})
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "occupied", value)

	sameLayer, err := registry.Layer(testBoardGridID, testStateLayerKey)
	require.NoError(t, err)
	assert.Same(t, layer, sameLayer)
}

func TestRegistry_LayerIfExists_DoesNotCreateLayer(t *testing.T) {
	t.Parallel()

	registry := NewRegistry[uint8, uint8, bool]()
	_, err := registry.CreateGrid(testBoardGridID, 3, 2)
	require.NoError(t, err)
	assert.Equal(t, 0, registry.LayerCount(testBoardGridID))

	layer, ok, err := registry.LayerIfExists(testBoardGridID, testStateLayerKey)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, layer)
	assert.Equal(t, 0, registry.LayerCount(testBoardGridID))
}

func TestRegistry_GridErrorsAndNilSafety(t *testing.T) {
	t.Parallel()

	registry := NewRegistry[uint8, uint8, bool]()
	_, err := registry.CreateGrid(testBoardGridID, 2, 2)
	require.NoError(t, err)

	_, err = registry.Layer(testMissingGridID, testStateLayerKey)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrGridNotFound)

	_, ok, err := registry.LayerIfExists(testMissingGridID, testStateLayerKey)
	require.NoError(t, err)
	assert.False(t, ok)

	var nilRegistry *Registry[uint8, uint8, bool]
	_, err = nilRegistry.Layer(testBoardGridID, testStateLayerKey)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrNilGridSet)

	_, _, err = nilRegistry.LayerIfExists(testBoardGridID, testStateLayerKey)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrNilGridSet)
}

func TestRegistry_DeleteGridRemovesGridAndLayers(t *testing.T) {
	t.Parallel()

	registry := NewRegistry[uint8, uint8, bool]()
	_, err := registry.CreateGrid(testBoardGridID, 2, 2)
	require.NoError(t, err)

	_, err = registry.Layer(testBoardGridID, testHitsLayerKey)
	require.NoError(t, err)
	assert.Equal(t, 1, registry.LayerCount(testBoardGridID))

	assert.True(t, registry.DeleteGrid(testBoardGridID))
	assert.Equal(t, 0, registry.LayerCount(testBoardGridID))
	assert.False(t, registry.DeleteGrid(testBoardGridID))

	_, err = registry.Layer(testBoardGridID, testHitsLayerKey)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrGridNotFound)
}
