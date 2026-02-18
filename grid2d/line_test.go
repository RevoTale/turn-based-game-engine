package grid2d

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectLineAxis(t *testing.T) {
	t.Parallel()

	axis, err := DetectLineAxis([]Position{
		{X: 1, Y: 4},
		{X: 3, Y: 4},
		{X: 2, Y: 4},
	})
	require.NoError(t, err)
	assert.Equal(t, LineAxisHorizontal, axis)

	axis, err = DetectLineAxis([]Position{
		{X: 7, Y: 2},
		{X: 7, Y: 1},
		{X: 7, Y: 3},
	})
	require.NoError(t, err)
	assert.Equal(t, LineAxisVertical, axis)

	_, err = DetectLineAxis(nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrEmptyPositions)

	_, err = DetectLineAxis([]Position{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotStraightLine)
}

func TestSortPositionsByAxis(t *testing.T) {
	t.Parallel()

	sorted, err := SortPositionsByAxis([]Position{
		{X: 4, Y: 1},
		{X: 2, Y: 1},
		{X: 3, Y: 1},
	}, LineAxisHorizontal)
	require.NoError(t, err)
	assert.Equal(t, []Position{
		{X: 2, Y: 1},
		{X: 3, Y: 1},
		{X: 4, Y: 1},
	}, sorted)

	sorted, err = SortPositionsByAxis([]Position{
		{X: 2, Y: 4},
		{X: 2, Y: 2},
		{X: 2, Y: 3},
	}, LineAxisVertical)
	require.NoError(t, err)
	assert.Equal(t, []Position{
		{X: 2, Y: 2},
		{X: 2, Y: 3},
		{X: 2, Y: 4},
	}, sorted)

	_, err = SortPositionsByAxis([]Position{{X: 1, Y: 1}}, LineAxis(255))
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidLineAxis)
}

func TestIsConsecutiveByAxis(t *testing.T) {
	t.Parallel()

	assert.True(t, IsConsecutiveByAxis([]Position{
		{X: 1, Y: 5},
		{X: 2, Y: 5},
		{X: 3, Y: 5},
	}, LineAxisHorizontal))

	assert.False(t, IsConsecutiveByAxis([]Position{
		{X: 1, Y: 5},
		{X: 3, Y: 5},
	}, LineAxisHorizontal))

	assert.True(t, IsConsecutiveByAxis([]Position{
		{X: 6, Y: 1},
		{X: 6, Y: 2},
		{X: 6, Y: 3},
	}, LineAxisVertical))

	assert.False(t, IsConsecutiveByAxis([]Position{
		{X: 6, Y: 1},
		{X: 6, Y: 3},
	}, LineAxisVertical))
}

func TestLineExtremes(t *testing.T) {
	t.Parallel()

	first, last, axis, err := LineExtremes([]Position{
		{X: 5, Y: 2},
		{X: 3, Y: 2},
		{X: 4, Y: 2},
	})
	require.NoError(t, err)
	assert.Equal(t, LineAxisHorizontal, axis)
	assert.Equal(t, Position{X: 3, Y: 2}, first)
	assert.Equal(t, Position{X: 5, Y: 2}, last)

	_, _, _, err = LineExtremes(nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrEmptyPositions)
}

func TestLineExtensions(t *testing.T) {
	t.Parallel()

	before, after, err := LineExtensions(
		Position{X: 2, Y: 7},
		Position{X: 4, Y: 7},
		LineAxisHorizontal,
	)
	require.NoError(t, err)
	assert.Equal(t, Position{X: 1, Y: 7}, before)
	assert.Equal(t, Position{X: 5, Y: 7}, after)

	before, after, err = LineExtensions(
		Position{X: 9, Y: 3},
		Position{X: 9, Y: 6},
		LineAxisVertical,
	)
	require.NoError(t, err)
	assert.Equal(t, Position{X: 9, Y: 2}, before)
	assert.Equal(t, Position{X: 9, Y: 7}, after)

	_, _, err = LineExtensions(Position{}, Position{}, LineAxis(255))
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidLineAxis)
}
