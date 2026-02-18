package grid2d

import "sort"

// LineAxis identifies the axis used by a straight line on the grid.
type LineAxis uint8

const (
	// LineAxisHorizontal means all positions share Y and advance over X.
	LineAxisHorizontal LineAxis = iota
	// LineAxisVertical means all positions share X and advance over Y.
	LineAxisVertical
)

// DetectLineAxis reports whether positions form a horizontal or vertical line.
//
// It returns ErrEmptyPositions when positions is empty and ErrNotStraightLine
// when positions are neither all in one row nor all in one column.
func DetectLineAxis(positions []Position) (LineAxis, error) {
	if len(positions) == 0 {
		return 0, ErrEmptyPositions
	}

	first := positions[0]
	sameX := true
	sameY := true
	for _, pos := range positions[1:] {
		if pos.X != first.X {
			sameX = false
		}
		if pos.Y != first.Y {
			sameY = false
		}
	}

	if !sameX && !sameY {
		return 0, ErrNotStraightLine
	}

	// Single-cell lines satisfy both conditions. Prefer horizontal as stable
	// default to keep deterministic sorting behavior.
	if sameY {
		return LineAxisHorizontal, nil
	}
	return LineAxisVertical, nil
}

// SortPositionsByAxis returns a sorted copy of positions for the given axis.
//
// Horizontal sorting is by X then Y. Vertical sorting is by Y then X.
// Returns ErrInvalidLineAxis when axis is unknown.
func SortPositionsByAxis(positions []Position, axis LineAxis) ([]Position, error) {
	switch axis {
	case LineAxisHorizontal, LineAxisVertical:
	default:
		return nil, ErrInvalidLineAxis
	}

	out := make([]Position, len(positions))
	copy(out, positions)

	if axis == LineAxisHorizontal {
		sort.Slice(out, func(i, j int) bool {
			if out[i].X == out[j].X {
				return out[i].Y < out[j].Y
			}
			return out[i].X < out[j].X
		})
		return out, nil
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Y == out[j].Y {
			return out[i].X < out[j].X
		}
		return out[i].Y < out[j].Y
	})
	return out, nil
}

// IsConsecutiveByAxis reports whether sorted positions form a contiguous line.
//
// Positions must be sorted according to axis before calling this function.
// It returns true for zero or one position.
func IsConsecutiveByAxis(positions []Position, axis LineAxis) bool {
	if len(positions) < 2 {
		return true
	}

	switch axis {
	case LineAxisHorizontal:
		for i := 1; i < len(positions); i++ {
			prev := positions[i-1]
			curr := positions[i]
			if curr.Y != prev.Y || curr.X != prev.X+1 {
				return false
			}
		}
		return true
	case LineAxisVertical:
		for i := 1; i < len(positions); i++ {
			prev := positions[i-1]
			curr := positions[i]
			if curr.X != prev.X || curr.Y != prev.Y+1 {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// LineExtremes returns edge positions for a straight horizontal or vertical line.
//
// The returned values are the first and last positions along line axis order.
// It returns ErrEmptyPositions when positions is empty and ErrNotStraightLine
// when positions are not in one row or one column.
func LineExtremes(positions []Position) (first Position, last Position, axis LineAxis, err error) {
	axis, err = DetectLineAxis(positions)
	if err != nil {
		return Position{}, Position{}, 0, err
	}

	sorted, err := SortPositionsByAxis(positions, axis)
	if err != nil {
		return Position{}, Position{}, 0, err
	}
	return sorted[0], sorted[len(sorted)-1], axis, nil
}

// LineExtensions returns one-step extension candidates for a line segment.
//
// `first` and `last` are expected to be extremes along the line axis.
// Returns ErrInvalidLineAxis when axis is unknown.
func LineExtensions(first Position, last Position, axis LineAxis) (before Position, after Position, err error) {
	switch axis {
	case LineAxisHorizontal:
		return Position{X: first.X - 1, Y: first.Y}, Position{X: last.X + 1, Y: last.Y}, nil
	case LineAxisVertical:
		return Position{X: first.X, Y: first.Y - 1}, Position{X: last.X, Y: last.Y + 1}, nil
	default:
		return Position{}, Position{}, ErrInvalidLineAxis
	}
}
