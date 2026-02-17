package grid2d

import "fmt"

// ConnectedComponent returns the connected region that matches `match`.
//
// The search starts from `start` and explores neighbors according to
// `neighborhood`. Positions for which `match` returns false are treated as
// blocked cells.
//
// Returns:
// - ErrNilGrid for nil receiver
// - ErrOutOfBounds when start is outside grid bounds
// - ErrInvalidNeighborhood when neighborhood is unknown
// - ErrNilMatcher when match is nil
func (g *Grid) ConnectedComponent(start Position, neighborhood Neighborhood, match func(pos Position) bool) ([]Position, error) {
	if g == nil {
		return nil, ErrNilGrid
	}
	if !g.InBounds(start) {
		return nil, fmt.Errorf("%w: (%d,%d)", ErrOutOfBounds, start.X, start.Y)
	}
	if match == nil {
		return nil, ErrNilMatcher
	}
	if _, err := directionsForNeighborhood(neighborhood); err != nil {
		return nil, err
	}

	if !match(start) {
		return nil, nil
	}

	component := make([]Position, 0)
	stack := []Position{start}
	seen := make(map[CellIndex]struct{}, 16)

	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]

		currentIndex, indexErr := g.Index(current)
		if indexErr != nil {
			continue
		}
		if _, alreadySeen := seen[currentIndex]; alreadySeen {
			continue
		}
		seen[currentIndex] = struct{}{}

		if !match(current) {
			continue
		}

		component = append(component, current)
		err := g.ForEachNeighbor(current, neighborhood, func(next Position) bool {
			nextIndex, nextIndexErr := g.Index(next)
			if nextIndexErr != nil {
				return true
			}
			if _, alreadySeen := seen[nextIndex]; alreadySeen {
				return true
			}
			stack = append(stack, next)
			return true
		})
		if err != nil {
			return nil, err
		}
	}

	return component, nil
}
