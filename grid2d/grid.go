package grid2d

import "fmt"

type Position struct {
	X int
	Y int
}

type CellIndex int

// Grid describes immutable 2D bounds.
type Grid struct {
	width     int
	height    int
	cellCount int
}

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

func (g *Grid) Width() int {
	if g == nil {
		return 0
	}
	return g.width
}

func (g *Grid) Height() int {
	if g == nil {
		return 0
	}
	return g.height
}

func (g *Grid) CellCount() int {
	if g == nil {
		return 0
	}
	return g.cellCount
}

func (g *Grid) InBounds(pos Position) bool {
	if g == nil {
		return false
	}
	return pos.X >= 0 && pos.X < g.width && pos.Y >= 0 && pos.Y < g.height
}

func (g *Grid) Index(pos Position) (CellIndex, error) {
	if g == nil {
		return 0, ErrNilGrid
	}
	if !g.InBounds(pos) {
		return 0, fmt.Errorf("%w: (%d,%d)", ErrOutOfBounds, pos.X, pos.Y)
	}
	return CellIndex(pos.Y*g.width + pos.X), nil
}

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
