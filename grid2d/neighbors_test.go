package grid2d

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForEachNeighbor_Moore(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(5, 5)
	require.NoError(t, err)

	var got []Position
	err = grid.ForEachNeighbor(Position{X: 2, Y: 2}, NeighborhoodMoore, func(pos Position) bool {
		got = append(got, pos)
		return true
	})
	require.NoError(t, err)

	assert.ElementsMatch(t, []Position{
		{X: 2, Y: 1},
		{X: 3, Y: 1},
		{X: 3, Y: 2},
		{X: 3, Y: 3},
		{X: 2, Y: 3},
		{X: 1, Y: 3},
		{X: 1, Y: 2},
		{X: 1, Y: 1},
	}, got)
}

func TestForEachNeighbor_VonNeumannAtEdge(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(3, 3)
	require.NoError(t, err)

	var got []Position
	err = grid.ForEachNeighbor(Position{X: 0, Y: 0}, NeighborhoodVonNeumann, func(pos Position) bool {
		got = append(got, pos)
		return true
	})
	require.NoError(t, err)

	assert.ElementsMatch(t, []Position{
		{X: 1, Y: 0},
		{X: 0, Y: 1},
	}, got)
}

func TestForEachNeighbor_ValidationAndEarlyStop(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(2, 2)
	require.NoError(t, err)

	err = grid.ForEachNeighbor(Position{X: 9, Y: 9}, NeighborhoodMoore, func(pos Position) bool {
		return true
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrOutOfBounds)

	var got []Position
	err = grid.ForEachNeighbor(Position{X: 0, Y: 0}, NeighborhoodMoore, func(pos Position) bool {
		got = append(got, pos)
		return false
	})
	require.NoError(t, err)
	assert.Len(t, got, 1)

	err = grid.ForEachNeighbor(Position{X: 0, Y: 0}, Neighborhood(255), func(pos Position) bool {
		return true
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidNeighborhood)
}
