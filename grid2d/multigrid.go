package grid2d

import "fmt"

// GridSet keeps independent grids addressed by id.
type GridSet[G comparable] struct {
	grids map[G]*Grid
}

func NewGridSet[G comparable]() *GridSet[G] {
	return &GridSet[G]{grids: make(map[G]*Grid)}
}

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

func (s *GridSet[G]) Get(id G) (*Grid, bool) {
	if s == nil {
		return nil, false
	}
	grid, ok := s.grids[id]
	return grid, ok
}

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

func (s *GridSet[G]) Count() int {
	if s == nil {
		return 0
	}
	return len(s.grids)
}

// ForEach iterates grids without making copies.
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

// MultiLayerSpace manages typed layer spaces for many grids.
type MultiLayerSpace[G comparable, K comparable, T comparable] struct {
	gridSet *GridSet[G]
	spaces  map[G]*LayerSpace[K, T]
}

func NewMultiLayerSpace[G comparable, K comparable, T comparable](gridSet *GridSet[G]) (*MultiLayerSpace[G, K, T], error) {
	if gridSet == nil {
		return nil, ErrNilGridSet
	}
	return &MultiLayerSpace[G, K, T]{
		gridSet: gridSet,
		spaces:  make(map[G]*LayerSpace[K, T]),
	}, nil
}

// Space returns the typed layer space for a grid, creating it lazily.
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

func (m *MultiLayerSpace[G, K, T]) ForgetGrid(gridID G) {
	if m == nil {
		return
	}
	delete(m.spaces, gridID)
}

func (m *MultiLayerSpace[G, K, T]) DeleteGrid(gridID G) bool {
	if m == nil || m.gridSet == nil {
		return false
	}
	delete(m.spaces, gridID)
	return m.gridSet.Delete(gridID)
}

func (m *MultiLayerSpace[G, K, T]) Count() int {
	if m == nil {
		return 0
	}
	return len(m.spaces)
}
