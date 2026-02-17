package grid2d

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectedComponent_VonNeumann(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(4, 4)
	require.NoError(t, err)

	blocked := map[Position]bool{
		{X: 0, Y: 0}: true,
		{X: 1, Y: 0}: true,
		{X: 1, Y: 1}: true,
		{X: 3, Y: 3}: true,
	}

	component, err := grid.ConnectedComponent(Position{X: 0, Y: 0}, NeighborhoodVonNeumann, func(pos Position) bool {
		return blocked[pos]
	})
	require.NoError(t, err)
	assert.ElementsMatch(t, []Position{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 1, Y: 1},
	}, component)
}

func TestConnectedComponent_ValidationAndNoMatch(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(2, 2)
	require.NoError(t, err)

	_, err = grid.ConnectedComponent(Position{X: 9, Y: 9}, NeighborhoodMoore, func(pos Position) bool { return true })
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrOutOfBounds)

	_, err = grid.ConnectedComponent(Position{X: 0, Y: 0}, Neighborhood(255), func(pos Position) bool { return true })
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidNeighborhood)

	_, err = grid.ConnectedComponent(Position{X: 0, Y: 0}, NeighborhoodMoore, nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNilMatcher)

	component, err := grid.ConnectedComponent(Position{X: 0, Y: 0}, NeighborhoodMoore, func(pos Position) bool { return false })
	require.NoError(t, err)
	assert.Nil(t, component)
}
