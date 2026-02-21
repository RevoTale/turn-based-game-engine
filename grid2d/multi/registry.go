package multi

import (
	"fmt"

	"github.com/RevoTale/turn-based-game-engine/grid2d"
)

// Registry stores many grids and per-grid typed layer spaces.
type Registry[G IntegerKey, K IntegerKey, T comparable] struct {
	grids  map[G]*grid2d.Grid
	spaces map[G]*LayerSpace[K, T]
}

// NewRegistry creates an empty registry.
func NewRegistry[G IntegerKey, K IntegerKey, T comparable]() *Registry[G, K, T] {
	return &Registry[G, K, T]{
		grids:  make(map[G]*grid2d.Grid),
		spaces: make(map[G]*LayerSpace[K, T]),
	}
}

// CreateGrid builds a new grid and stores it under id.
//
// Returns grid2d.ErrNilGridSet for nil receiver, grid2d.ErrGridExists when id
// already exists, and validation errors from grid2d.NewGrid for invalid
// dimensions.
func (r *Registry[G, K, T]) CreateGrid(id G, width, height int) (*grid2d.Grid, error) {
	if r == nil {
		return nil, grid2d.ErrNilGridSet
	}
	if _, exists := r.grids[id]; exists {
		return nil, fmt.Errorf("%w: %v", grid2d.ErrGridExists, id)
	}
	grid, err := grid2d.NewGrid(width, height)
	if err != nil {
		return nil, err
	}
	r.grids[id] = grid
	return grid, nil
}

// AddGrid stores an existing grid under id.
//
// Returns grid2d.ErrNilGridSet for nil receiver, grid2d.ErrNilGrid when grid is
// nil, and grid2d.ErrGridExists when id already exists.
func (r *Registry[G, K, T]) AddGrid(id G, grid *grid2d.Grid) error {
	if r == nil {
		return grid2d.ErrNilGridSet
	}
	if grid == nil {
		return grid2d.ErrNilGrid
	}
	if _, exists := r.grids[id]; exists {
		return fmt.Errorf("%w: %v", grid2d.ErrGridExists, id)
	}
	r.grids[id] = grid
	return nil
}

// GetGrid returns grid by id.
//
// Second return value is false when id is not found.
func (r *Registry[G, K, T]) GetGrid(id G) (*grid2d.Grid, bool) {
	if r == nil {
		return nil, false
	}
	grid, ok := r.grids[id]
	return grid, ok
}

// DeleteGrid removes layer space and then removes the grid.
//
// Returns true when the grid was removed.
func (r *Registry[G, K, T]) DeleteGrid(id G) bool {
	if r == nil {
		return false
	}
	delete(r.spaces, id)
	if _, ok := r.grids[id]; !ok {
		return false
	}
	delete(r.grids, id)
	return true
}

// GridCount returns how many grids are in the registry.
func (r *Registry[G, K, T]) GridCount() int {
	if r == nil {
		return 0
	}
	return len(r.grids)
}

// ForEachGrid calls visit for each grid.
//
// Returning false from visit stops iteration.
func (r *Registry[G, K, T]) ForEachGrid(visit func(id G, grid *grid2d.Grid) bool) {
	if r == nil || visit == nil {
		return
	}
	for id, grid := range r.grids {
		if !visit(id, grid) {
			return
		}
	}
}

// Space returns layer space for gridID.
//
// If the layer space does not exist yet, Space creates it.
func (r *Registry[G, K, T]) Space(gridID G) (*LayerSpace[K, T], error) {
	if r == nil {
		return nil, grid2d.ErrNilGridSet
	}
	if space, ok := r.spaces[gridID]; ok {
		return space, nil
	}
	grid, ok := r.grids[gridID]
	if !ok {
		return nil, fmt.Errorf("%w: %v", grid2d.ErrGridNotFound, gridID)
	}
	space, err := NewLayerSpace[K, T](grid)
	if err != nil {
		return nil, err
	}
	r.spaces[gridID] = space
	return space, nil
}

// SpaceIfExists returns layer space only if it already exists.
//
// It never creates a new layer space.
// Missing grid ids return (nil, false, nil).
func (r *Registry[G, K, T]) SpaceIfExists(gridID G) (*LayerSpace[K, T], bool, error) {
	if r == nil {
		return nil, false, grid2d.ErrNilGridSet
	}
	if _, ok := r.grids[gridID]; !ok {
		return nil, false, nil
	}
	space, ok := r.spaces[gridID]
	return space, ok, nil
}

// ForgetGrid removes layer data for gridID but keeps the grid itself.
func (r *Registry[G, K, T]) ForgetGrid(gridID G) {
	if r == nil {
		return
	}
	delete(r.spaces, gridID)
}

// LayerSpaceCount returns how many layer spaces currently exist.
func (r *Registry[G, K, T]) LayerSpaceCount() int {
	if r == nil {
		return 0
	}
	return len(r.spaces)
}
