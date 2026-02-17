package grid2d

import "fmt"

// Position identifies a cell location in 2D coordinates.
//
// X grows to the right, Y grows downward.
type Position struct {
	X int
	Y int
}

// CellIndex is a single-number cell address inside a grid.
type CellIndex int

// Grid stores board size and helps convert between coordinates and indexes.
type Grid struct {
	width     int
	height    int
	cellCount int
}

// NewGrid creates a grid with the given width and height.
//
// Width and height must be positive.
// Returns ErrGridTooLarge when width*height cannot fit in int.
func NewGrid(width, height int) (*Grid, error) {
	if width <= 0 || height <= 0 {
		return nil, ErrInvalidGridSize
	}

	maxInt := int(^uint(0) >> 1)
	if width > maxInt/height {
		return nil, ErrGridTooLarge
	}

	return &Grid{
		width:     width,
		height:    height,
		cellCount: width * height,
	}, nil
}

// Width returns the grid width.
//
// It returns 0 for a nil receiver.
func (g *Grid) Width() int {
	if g == nil {
		return 0
	}
	return g.width
}

// Height returns the grid height.
//
// It returns 0 for a nil receiver.
func (g *Grid) Height() int {
	if g == nil {
		return 0
	}
	return g.height
}

// CellCount returns total number of cells (width * height).
//
// It returns 0 for a nil receiver.
func (g *Grid) CellCount() int {
	if g == nil {
		return 0
	}
	return g.cellCount
}

// InBounds reports whether pos is inside the grid.
//
// It returns false for a nil receiver.
func (g *Grid) InBounds(pos Position) bool {
	if g == nil {
		return false
	}
	return pos.X >= 0 && pos.X < g.width && pos.Y >= 0 && pos.Y < g.height
}

// Index converts coordinates to CellIndex.
//
// This is useful when storing cell values in a map keyed by CellIndex.
//
// Returns ErrNilGrid when receiver is nil.
// Returns ErrOutOfBounds when position is outside grid bounds.
func (g *Grid) Index(pos Position) (CellIndex, error) {
	if g == nil {
		return 0, ErrNilGrid
	}
	if !g.InBounds(pos) {
		return 0, fmt.Errorf("%w: (%d,%d)", ErrOutOfBounds, pos.X, pos.Y)
	}
	return CellIndex(pos.Y*g.width + pos.X), nil
}

// Position converts CellIndex back to coordinates.
//
// Second return value is false when index is invalid or receiver is nil.
func (g *Grid) Position(index CellIndex) (Position, bool) {
	if g == nil {
		return Position{}, false
	}

	i := int(index)
	if i < 0 || i >= g.cellCount {
		return Position{}, false
	}
	return Position{
		X: i % g.width,
		Y: i / g.width,
	}, true
}
