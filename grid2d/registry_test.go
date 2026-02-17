package grid2d

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLayerRegistry_RequiresGridSet(t *testing.T) {
	t.Parallel()

	registry, err := NewLayerRegistry[string, string, bool](nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilGridSet)
	assert.Nil(t, registry)
}

func TestLayerRegistry_NonEnsureOperationsDoNotCreateLayerSpace(t *testing.T) {
	t.Parallel()

	grids := NewGridSet[string]()
	_, err := grids.Create("board", 3, 2)
	require.NoError(t, err)

	registry, err := NewLayerRegistry[string, string, bool](grids)
	require.NoError(t, err)
	assert.Equal(t, 0, registry.LayerSpaceCount())

	layer, ok, err := registry.Get("board", "hits")
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, layer)
	assert.Equal(t, 0, registry.LayerSpaceCount())

	err = registry.Set("board", "hits", Position{X: 1, Y: 1}, true)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrLayerNotFound)
	assert.Equal(t, 0, registry.LayerSpaceCount())

	err = registry.DeleteValue("board", "hits", Position{X: 1, Y: 1})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrLayerNotFound)
	assert.Equal(t, 0, registry.LayerSpaceCount())

	deleted, err := registry.DeleteLayer("board", "hits")
	require.NoError(t, err)
	assert.False(t, deleted)
	assert.Equal(t, 0, registry.LayerSpaceCount())
}

func TestLayerRegistry_EnsureAndEnsureSet(t *testing.T) {
	t.Parallel()

	grids := NewGridSet[string]()
	_, err := grids.Create("board", 3, 2)
	require.NoError(t, err)

	registry, err := NewLayerRegistry[string, string, string](grids)
	require.NoError(t, err)

	first, err := registry.Ensure("board", "status")
	require.NoError(t, err)
	assert.Equal(t, 1, registry.LayerSpaceCount())

	second, err := registry.Ensure("board", "status")
	require.NoError(t, err)
	assert.Same(t, first, second)
	assert.Equal(t, 1, registry.LayerSpaceCount())

	err = registry.EnsureSet("board", "status", Position{X: 2, Y: 1}, "occupied")
	require.NoError(t, err)

	value, ok, err := first.Get(Position{X: 2, Y: 1})
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "occupied", value)
}

func TestLayerRegistry_GridErrorsAndNilSafety(t *testing.T) {
	t.Parallel()

	grids := NewGridSet[string]()
	_, err := grids.Create("board", 2, 2)
	require.NoError(t, err)

	registry, err := NewLayerRegistry[string, string, bool](grids)
	require.NoError(t, err)

	_, _, err = registry.Get("missing", "hits")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrGridNotFound)

	_, err = registry.Ensure("missing", "hits")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrGridNotFound)

	err = registry.Set("missing", "hits", Position{X: 0, Y: 0}, true)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrGridNotFound)

	_, err = registry.DeleteLayer("missing", "hits")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrGridNotFound)

	var nilRegistry *LayerRegistry[string, string, bool]

	_, err = nilRegistry.Ensure("board", "hits")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilLayerRegistry)

	_, _, err = nilRegistry.Get("board", "hits")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilLayerRegistry)

	err = nilRegistry.Set("board", "hits", Position{X: 0, Y: 0}, true)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilLayerRegistry)

	err = nilRegistry.DeleteValue("board", "hits", Position{X: 0, Y: 0})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilLayerRegistry)

	deleted, err := nilRegistry.DeleteLayer("board", "hits")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilLayerRegistry)
	assert.False(t, deleted)

	assert.False(t, nilRegistry.DeleteGrid("board"))
	assert.Equal(t, 0, nilRegistry.LayerSpaceCount())
}

func TestLayerRegistry_DeleteGridRemovesGridAndSpace(t *testing.T) {
	t.Parallel()

	grids := NewGridSet[string]()
	_, err := grids.Create("board", 2, 2)
	require.NoError(t, err)

	registry, err := NewLayerRegistry[string, string, bool](grids)
	require.NoError(t, err)

	_, err = registry.Ensure("board", "hits")
	require.NoError(t, err)
	assert.Equal(t, 1, registry.LayerSpaceCount())

	assert.True(t, registry.DeleteGrid("board"))
	assert.Equal(t, 0, registry.LayerSpaceCount())
	assert.False(t, registry.DeleteGrid("board"))

	_, err = registry.Ensure("board", "hits")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrGridNotFound)
}
