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

func TestRegistry_NonEnsureOperationsDoNotCreateLayerSpace(t *testing.T) {
	t.Parallel()

	registry := NewRegistry[uint8, uint8, bool]()
	_, err := registry.CreateGrid(testBoardGridID, 3, 2)
	require.NoError(t, err)
	assert.Equal(t, 0, registry.LayerSpaceCount())

	layer, ok, err := registry.Get(testBoardGridID, testHitsLayerKey)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, layer)
	assert.Equal(t, 0, registry.LayerSpaceCount())

	err = registry.Set(testBoardGridID, testHitsLayerKey, grid2d.Position{X: 1, Y: 1}, true)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrLayerNotFound)
	assert.Equal(t, 0, registry.LayerSpaceCount())

	err = registry.DeleteValue(testBoardGridID, testHitsLayerKey, grid2d.Position{X: 1, Y: 1})
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrLayerNotFound)
	assert.Equal(t, 0, registry.LayerSpaceCount())

	deleted, err := registry.DeleteLayer(testBoardGridID, testHitsLayerKey)
	require.NoError(t, err)
	assert.False(t, deleted)
	assert.Equal(t, 0, registry.LayerSpaceCount())
}

func TestRegistry_EnsureAndEnsureSet(t *testing.T) {
	t.Parallel()

	registry := NewRegistry[uint8, uint8, string]()
	_, err := registry.CreateGrid(testBoardGridID, 3, 2)
	require.NoError(t, err)

	first, err := registry.Ensure(testBoardGridID, testStateLayerKey)
	require.NoError(t, err)
	assert.Equal(t, 1, registry.LayerSpaceCount())

	second, err := registry.Ensure(testBoardGridID, testStateLayerKey)
	require.NoError(t, err)
	assert.Same(t, first, second)
	assert.Equal(t, 1, registry.LayerSpaceCount())

	err = registry.EnsureSet(testBoardGridID, testStateLayerKey, grid2d.Position{X: 2, Y: 1}, "occupied")
	require.NoError(t, err)

	value, ok, err := first.Get(grid2d.Position{X: 2, Y: 1})
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "occupied", value)
}

func TestRegistry_GridErrorsAndNilSafety(t *testing.T) {
	t.Parallel()

	registry := NewRegistry[uint8, uint8, bool]()
	_, err := registry.CreateGrid(testBoardGridID, 2, 2)
	require.NoError(t, err)

	_, _, err = registry.Get(testMissingGridID, testHitsLayerKey)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrGridNotFound)

	_, err = registry.Ensure(testMissingGridID, testHitsLayerKey)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrGridNotFound)

	err = registry.Set(testMissingGridID, testHitsLayerKey, grid2d.Position{X: 0, Y: 0}, true)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrGridNotFound)

	_, err = registry.DeleteLayer(testMissingGridID, testHitsLayerKey)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrGridNotFound)

	var nilRegistry *Registry[uint8, uint8, bool]

	_, err = nilRegistry.Ensure(testBoardGridID, testHitsLayerKey)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrNilLayerRegistry)

	_, _, err = nilRegistry.Get(testBoardGridID, testHitsLayerKey)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrNilLayerRegistry)

	err = nilRegistry.Set(testBoardGridID, testHitsLayerKey, grid2d.Position{X: 0, Y: 0}, true)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrNilLayerRegistry)

	err = nilRegistry.DeleteValue(testBoardGridID, testHitsLayerKey, grid2d.Position{X: 0, Y: 0})
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrNilLayerRegistry)

	deleted, err := nilRegistry.DeleteLayer(testBoardGridID, testHitsLayerKey)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrNilLayerRegistry)
	assert.False(t, deleted)

	assert.False(t, nilRegistry.DeleteGrid(testBoardGridID))
	assert.Equal(t, 0, nilRegistry.LayerSpaceCount())
}

func TestRegistry_DeleteGridRemovesGridAndSpace(t *testing.T) {
	t.Parallel()

	registry := NewRegistry[uint8, uint8, bool]()
	_, err := registry.CreateGrid(testBoardGridID, 2, 2)
	require.NoError(t, err)

	_, err = registry.Ensure(testBoardGridID, testHitsLayerKey)
	require.NoError(t, err)
	assert.Equal(t, 1, registry.LayerSpaceCount())

	assert.True(t, registry.DeleteGrid(testBoardGridID))
	assert.Equal(t, 0, registry.LayerSpaceCount())
	assert.False(t, registry.DeleteGrid(testBoardGridID))

	_, err = registry.Ensure(testBoardGridID, testHitsLayerKey)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrGridNotFound)
}
