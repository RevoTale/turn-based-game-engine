package grid2d

import "fmt"

// GridSet stores many grids by id.
type GridSet[G comparable] struct {
	grids map[G]*Grid
}

// NewGridSet creates an empty GridSet.
func NewGridSet[G comparable]() *GridSet[G] {
	return &GridSet[G]{grids: make(map[G]*Grid)}
}

// Create builds a new grid and stores it under id.
//
// Returns ErrNilGridSet for nil receiver, ErrGridExists when id already exists,
// and validation errors from NewGrid for invalid dimensions.
func (s *GridSet[G]) Create(id G, width, height int) (*Grid, error) {
	if s == nil {
		return nil, ErrNilGridSet
	}
	if _, exists := s.grids[id]; exists {
		return nil, fmt.Errorf("%w: %v", ErrGridExists, id)
	}
	grid, err := NewGrid(width, height)
	if err != nil {
		return nil, err
	}
	s.grids[id] = grid
	return grid, nil
}

// Add stores an existing grid under id.
//
// Returns ErrNilGridSet for nil receiver, ErrNilGrid when grid is nil,
// and ErrGridExists when id already exists.
func (s *GridSet[G]) Add(id G, grid *Grid) error {
	if s == nil {
		return ErrNilGridSet
	}
	if grid == nil {
		return ErrNilGrid
	}
	if _, exists := s.grids[id]; exists {
		return fmt.Errorf("%w: %v", ErrGridExists, id)
	}
	s.grids[id] = grid
	return nil
}

// Get returns grid by id.
//
// Second return value is false when id is not found.
func (s *GridSet[G]) Get(id G) (*Grid, bool) {
	if s == nil {
		return nil, false
	}
	grid, ok := s.grids[id]
	return grid, ok
}

// Delete removes grid by id.
//
// Returns true when a grid was removed.
func (s *GridSet[G]) Delete(id G) bool {
	if s == nil {
		return false
	}
	if _, ok := s.grids[id]; !ok {
		return false
	}
	delete(s.grids, id)
	return true
}

// Count returns how many grids are in the set.
func (s *GridSet[G]) Count() int {
	if s == nil {
		return 0
	}
	return len(s.grids)
}

// ForEach calls visit for each grid.
//
// Returning false from visit stops iteration.
func (s *GridSet[G]) ForEach(visit func(id G, grid *Grid) bool) {
	if s == nil || visit == nil {
		return
	}
	for id, grid := range s.grids {
		if !visit(id, grid) {
			return
		}
	}
}

// MultiLayerSpace stores per-grid LayerSpace values for many grids.
type MultiLayerSpace[G comparable, K comparable, T comparable] struct {
	gridSet *GridSet[G]
	spaces  map[G]*LayerSpace[K, T]
}

// NewMultiLayerSpace creates an empty per-grid layer manager.
//
// Returns ErrNilGridSet when gridSet is nil.
func NewMultiLayerSpace[G comparable, K comparable, T comparable](gridSet *GridSet[G]) (*MultiLayerSpace[G, K, T], error) {
	if gridSet == nil {
		return nil, ErrNilGridSet
	}
	return &MultiLayerSpace[G, K, T]{
		gridSet: gridSet,
		spaces:  make(map[G]*LayerSpace[K, T]),
	}, nil
}

// Space returns layer space for gridID.
//
// If the layer space does not exist yet, Space creates it.
func (m *MultiLayerSpace[G, K, T]) Space(gridID G) (*LayerSpace[K, T], error) {
	if m == nil || m.gridSet == nil {
		return nil, ErrNilGridSet
	}
	if space, ok := m.spaces[gridID]; ok {
		return space, nil
	}

	grid, ok := m.gridSet.Get(gridID)
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrGridNotFound, gridID)
	}

	space, err := NewLayerSpace[K, T](grid)
	if err != nil {
		return nil, err
	}
	m.spaces[gridID] = space
	return space, nil
}

// SpaceIfExists returns layer space only if it already exists.
//
// It never creates a new layer space.
func (m *MultiLayerSpace[G, K, T]) SpaceIfExists(gridID G) (*LayerSpace[K, T], bool, error) {
	if m == nil || m.gridSet == nil {
		return nil, false, ErrNilGridSet
	}
	if _, ok := m.gridSet.Get(gridID); !ok {
		return nil, false, fmt.Errorf("%w: %v", ErrGridNotFound, gridID)
	}
	space, ok := m.spaces[gridID]
	return space, ok, nil
}

// ForgetGrid removes layer data for gridID but keeps the grid itself.
func (m *MultiLayerSpace[G, K, T]) ForgetGrid(gridID G) {
	if m == nil {
		return
	}
	delete(m.spaces, gridID)
}

// DeleteGrid removes layer data and then removes the grid from GridSet.
//
// Returns true when the grid was removed from GridSet.
func (m *MultiLayerSpace[G, K, T]) DeleteGrid(gridID G) bool {
	if m == nil || m.gridSet == nil {
		return false
	}
	delete(m.spaces, gridID)
	return m.gridSet.Delete(gridID)
}

// Count returns how many layer spaces currently exist.
//
// This is different from total grid count because spaces are created on demand.
func (m *MultiLayerSpace[G, K, T]) Count() int {
	if m == nil {
		return 0
	}
	return len(m.spaces)
}
