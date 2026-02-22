package multi

import (
	"fmt"

	"github.com/RevoTale/turn-based-game-engine/grid2d"
)

// Registry stores many grids and per-grid named sparse layers.
type Registry[G IntegerKey, K IntegerKey, T comparable] struct {
	grids  map[G]*grid2d.Grid
	layers map[G]map[K]*grid2d.SparseLayer[T]
}

// NewRegistry creates an empty registry.
func NewRegistry[G IntegerKey, K IntegerKey, T comparable]() *Registry[G, K, T] {
	return &Registry[G, K, T]{
		grids:  make(map[G]*grid2d.Grid),
		layers: make(map[G]map[K]*grid2d.SparseLayer[T]),
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

// DeleteGrid removes all layers and then removes the grid.
//
// Returns true when the grid was removed.
func (r *Registry[G, K, T]) DeleteGrid(id G) bool {
	if r == nil {
		return false
	}
	delete(r.layers, id)
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

// Layer returns named layer for gridID and creates it when missing.
func (r *Registry[G, K, T]) Layer(gridID G, layerKey K) (*grid2d.SparseLayer[T], error) {
	if r == nil {
		return nil, grid2d.ErrNilGridSet
	}
	grid, ok := r.grids[gridID]
	if !ok {
		return nil, fmt.Errorf("%w: %v", grid2d.ErrGridNotFound, gridID)
	}
	layers := r.layers[gridID]
	if layers == nil {
		layers = make(map[K]*grid2d.SparseLayer[T])
		r.layers[gridID] = layers
	}
	if layer, exists := layers[layerKey]; exists {
		return layer, nil
	}
	layer, err := grid2d.NewSparseLayer[T](grid)
	if err != nil {
		return nil, err
	}
	layers[layerKey] = layer
	return layer, nil
}

// LayerIfExists returns named layer only if already created.
//
// It never creates a layer.
// Missing grid or missing layer returns (nil, false, nil).
func (r *Registry[G, K, T]) LayerIfExists(gridID G, layerKey K) (*grid2d.SparseLayer[T], bool, error) {
	if r == nil {
		return nil, false, grid2d.ErrNilGridSet
	}
	if _, ok := r.grids[gridID]; !ok {
		return nil, false, nil
	}
	layers := r.layers[gridID]
	if layers == nil {
		return nil, false, nil
	}
	layer, ok := layers[layerKey]
	return layer, ok, nil
}

// DeleteLayer removes one named layer from a grid.
func (r *Registry[G, K, T]) DeleteLayer(gridID G, layerKey K) bool {
	if r == nil {
		return false
	}
	layers := r.layers[gridID]
	if layers == nil {
		return false
	}
	if _, ok := layers[layerKey]; !ok {
		return false
	}
	delete(layers, layerKey)
	if len(layers) == 0 {
		delete(r.layers, gridID)
	}
	return true
}

// LayerCount returns how many layers currently exist for gridID.
func (r *Registry[G, K, T]) LayerCount(gridID G) int {
	if r == nil {
		return 0
	}
	return len(r.layers[gridID])
}

// ForgetGridLayers removes all layers for one grid but keeps the grid itself.
func (r *Registry[G, K, T]) ForgetGridLayers(gridID G) {
	if r == nil {
		return
	}
	delete(r.layers, gridID)
}
