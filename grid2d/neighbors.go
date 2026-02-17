package grid2d

import "fmt"

// Neighborhood defines which adjacent cells are visited.
type Neighborhood uint8

const (
	// NeighborhoodVonNeumann visits orthogonal neighbors only (N, E, S, W).
	NeighborhoodVonNeumann Neighborhood = iota
	// NeighborhoodMoore visits all 8 surrounding neighbors.
	NeighborhoodMoore
)

var vonNeumannDirections = [...]Position{
	{X: 0, Y: -1},
	{X: 1, Y: 0},
	{X: 0, Y: 1},
	{X: -1, Y: 0},
}

var mooreDirections = [...]Position{
	{X: 0, Y: -1},  // North
	{X: 1, Y: -1},  // Northeast
	{X: 1, Y: 0},   // East
	{X: 1, Y: 1},   // Southeast
	{X: 0, Y: 1},   // South
	{X: -1, Y: 1},  // Southwest
	{X: -1, Y: 0},  // West
	{X: -1, Y: -1}, // Northwest
}

func directionsForNeighborhood(neighborhood Neighborhood) ([]Position, error) {
	switch neighborhood {
	case NeighborhoodVonNeumann:
		return vonNeumannDirections[:], nil
	case NeighborhoodMoore:
		return mooreDirections[:], nil
	default:
		return nil, ErrInvalidNeighborhood
	}
}

// ForEachNeighbor visits valid neighbor cells around center.
//
// It returns ErrNilGrid for a nil receiver and ErrOutOfBounds when center is
// outside grid bounds. Unknown neighborhood values return
// ErrInvalidNeighborhood.
//
// Returning false from visit stops iteration early.
func (g *Grid) ForEachNeighbor(center Position, neighborhood Neighborhood, visit func(pos Position) bool) error {
	if g == nil {
		return ErrNilGrid
	}
	if !g.InBounds(center) {
		return fmt.Errorf("%w: (%d,%d)", ErrOutOfBounds, center.X, center.Y)
	}
	if visit == nil {
		return nil
	}

	directions, err := directionsForNeighborhood(neighborhood)
	if err != nil {
		return err
	}

	for _, d := range directions {
		next := Position{X: center.X + d.X, Y: center.Y + d.Y}
		if !g.InBounds(next) {
			continue
		}
		if !visit(next) {
			return nil
		}
	}
	return nil
}
