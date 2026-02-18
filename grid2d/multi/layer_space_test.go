package multi

import (
	"testing"

	"github.com/RevoTale/turn-based-game-engine/grid2d"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLayerSpace_CreateAndDuplicate(t *testing.T) {
	t.Parallel()

	grid, err := grid2d.NewGrid(5, 3)
	require.NoError(t, err)

	const terrainLayerKey uint8 = 1
	space, err := NewLayerSpace[uint8, uint16](grid)
	require.NoError(t, err)

	layer, err := space.Create(terrainLayerKey)
	require.NoError(t, err)
	require.NotNil(t, layer)
	assert.Equal(t, 1, space.Count())

	_, err = space.Create(terrainLayerKey)
	require.Error(t, err)
	assert.ErrorIs(t, err, grid2d.ErrLayerExists)
	assert.Equal(t, grid2d.CodeLayerExists, grid2d.CodeOf(err))
}
