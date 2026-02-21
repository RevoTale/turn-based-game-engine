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

func TestRegistry_SpaceAndLayerOperations(t *testing.T) {
	t.Parallel()

	registry := NewRegistry[uint8, uint8, string]()
	_, err := registry.CreateGrid(testBoardGridID, 3, 2)
	require.NoError(t, err)
	assert.Equal(t, 0, registry.LayerSpaceCount())

	space, err := registry.Space(testBoardGridID)
	require.NoError(t, err)
	assert.Equal(t, 1, registry.LayerSpaceCount())

	layer, err := space.Ensure(testStateLayerKey)
	require.NoError(t, err)
	require.NoError(t, layer.Set(grid2d.Position{X: 2, Y: 1}, "occupied"))

	value, ok, err := layer.Get(grid2d.Position{X: 2, Y: 1})
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "occupied", value)
}

func TestRegistry_SpaceIfExists_DoesNotCreateSpace(t *testing.T) {
	t.Parallel()

	registry := NewRegistry[uint8, uint8, bool]()
	_, err := registry.CreateGrid(testBoardGridID, 3, 2)
	require.NoError(t, err)
	assert.Equal(t, 0, registry.LayerSpaceCount())

	space, ok, err := registry.SpaceIfExists(testBoardGridID)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, space)
	assert.Equal(t, 0, registry.LayerSpaceCount())
}

func TestRegistry_GridErrorsAndNilSafety(t *testing.T) {
	t.Parallel()

	registry := NewRegistry[uint8, uint8, bool]()
	_, err := registry.CreateGrid(testBoardGridID, 2, 2)
	require.NoError(t, err)

	_, err = registry.Space(testMissingGridID)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrGridNotFound)

	_, ok, err := registry.SpaceIfExists(testMissingGridID)
	require.NoError(t, err)
	assert.False(t, ok)

	var nilRegistry *Registry[uint8, uint8, bool]
	_, err = nilRegistry.Space(testBoardGridID)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrNilGridSet)

	_, _, err = nilRegistry.SpaceIfExists(testBoardGridID)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrNilGridSet)
}

func TestRegistry_DeleteGridRemovesGridAndSpace(t *testing.T) {
	t.Parallel()

	registry := NewRegistry[uint8, uint8, bool]()
	_, err := registry.CreateGrid(testBoardGridID, 2, 2)
	require.NoError(t, err)

	space, err := registry.Space(testBoardGridID)
	require.NoError(t, err)
	_, err = space.Ensure(testHitsLayerKey)
	require.NoError(t, err)
	assert.Equal(t, 1, registry.LayerSpaceCount())

	assert.True(t, registry.DeleteGrid(testBoardGridID))
	assert.Equal(t, 0, registry.LayerSpaceCount())
	assert.False(t, registry.DeleteGrid(testBoardGridID))

	_, err = registry.Space(testBoardGridID)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrGridNotFound)
}
