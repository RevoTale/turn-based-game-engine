package grid2d

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGrid_Validation(t *testing.T) {
	t.Parallel()

	_, err := NewGrid(0, 1)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidGridSize)
	assert.Equal(t, CodeInvalidGridSize, CodeOf(err))

	_, err = NewGrid(2, 0)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidGridSize)
}

func TestGrid_IndexAndPosition_NonSquare(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(3, 2)
	require.NoError(t, err)
	assert.Equal(t, 6, grid.CellCount())

	index, err := grid.Index(Position{X: 2, Y: 1})
	require.NoError(t, err)
	assert.Equal(t, CellIndex(5), index)

	pos, ok := grid.Position(index)
	require.True(t, ok)
	assert.Equal(t, Position{X: 2, Y: 1}, pos)
}

func TestGrid_OutOfBounds(t *testing.T) {
	t.Parallel()

	grid, err := NewGrid(4, 2)
	require.NoError(t, err)

	_, err = grid.Index(Position{X: 4, Y: 1})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrOutOfBounds)
	assert.Equal(t, CodeOutOfBounds, CodeOf(err))

	_, err = grid.Index(Position{X: 2, Y: 2})
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrOutOfBounds))

	_, ok := grid.Position(CellIndex(-1))
	assert.False(t, ok)

	_, ok = grid.Position(CellIndex(8))
	assert.False(t, ok)
}
