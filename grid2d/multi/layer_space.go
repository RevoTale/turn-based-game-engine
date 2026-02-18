package multi

import (
	"fmt"

	"github.com/RevoTale/turn-based-game-engine/grid2d"
)

// LayerSpace stores many named layers for one grid.
type LayerSpace[K IntegerKey, T comparable] struct {
	grid   *grid2d.Grid
	layers map[K]*grid2d.SparseLayer[T]
}

// NewLayerSpace creates an empty layer collection for one grid.
//
// Returns grid2d.ErrNilGrid when grid is nil.
func NewLayerSpace[K IntegerKey, T comparable](grid *grid2d.Grid) (*LayerSpace[K, T], error) {
	if grid == nil {
		return nil, grid2d.ErrNilGrid
	}
	return &LayerSpace[K, T]{
		grid:   grid,
		layers: make(map[K]*grid2d.SparseLayer[T]),
	}, nil
}

// Grid returns the shared grid used by all layers in this space.
func (s *LayerSpace[K, T]) Grid() *grid2d.Grid {
	if s == nil {
		return nil
	}
	return s.grid
}

// Create adds a new empty layer with key.
//
// Returns grid2d.ErrNilGrid for nil space/grid and grid2d.ErrLayerExists when
// key already exists.
func (s *LayerSpace[K, T]) Create(key K) (*grid2d.SparseLayer[T], error) {
	if s == nil || s.grid == nil {
		return nil, grid2d.ErrNilGrid
	}
	if _, exists := s.layers[key]; exists {
		return nil, fmt.Errorf("%w: %v", grid2d.ErrLayerExists, key)
	}
	layer, err := grid2d.NewSparseLayer[T](s.grid)
	if err != nil {
		return nil, err
	}
	s.layers[key] = layer
	return layer, nil
}

// Get returns layer by key.
//
// Second return value is false when key is missing.
func (s *LayerSpace[K, T]) Get(key K) (*grid2d.SparseLayer[T], bool) {
	if s == nil {
		return nil, false
	}
	layer, ok := s.layers[key]
	return layer, ok
}

// Delete removes layer by key.
//
// Returns true when a layer was removed.
func (s *LayerSpace[K, T]) Delete(key K) bool {
	if s == nil {
		return false
	}
	if _, ok := s.layers[key]; !ok {
		return false
	}
	delete(s.layers, key)
	return true
}

// Count returns how many layers are in this space.
func (s *LayerSpace[K, T]) Count() int {
	if s == nil {
		return 0
	}
	return len(s.layers)
}

// ForEach calls visit for each layer.
//
// Returning false from visit stops iteration.
func (s *LayerSpace[K, T]) ForEach(visit func(key K, layer *grid2d.SparseLayer[T]) bool) {
	if s == nil || visit == nil {
		return
	}
	for key, layer := range s.layers {
		if !visit(key, layer) {
			return
		}
	}
}
