package grid2d

import "fmt"

// SparseLayer stores only written cell values.
type SparseLayer[T comparable] struct {
	grid   *Grid
	values map[CellIndex]T
}

func NewSparseLayer[T comparable](grid *Grid) (*SparseLayer[T], error) {
	if grid == nil {
		return nil, ErrNilGrid
	}
	return &SparseLayer[T]{grid: grid}, nil
}

func (l *SparseLayer[T]) Grid() *Grid {
	if l == nil {
		return nil
	}
	return l.grid
}

func (l *SparseLayer[T]) Len() int {
	if l == nil || l.values == nil {
		return 0
	}
	return len(l.values)
}

func (l *SparseLayer[T]) Set(pos Position, value T) error {
	index, err := l.index(pos)
	if err != nil {
		return err
	}
	return l.SetIndex(index, value)
}

func (l *SparseLayer[T]) SetIndex(index CellIndex, value T) error {
	if err := l.validateIndex(index); err != nil {
		return err
	}
	if l.values == nil {
		l.values = make(map[CellIndex]T)
	}
	l.values[index] = value
	return nil
}

func (l *SparseLayer[T]) Get(pos Position) (T, bool, error) {
	index, err := l.index(pos)
	if err != nil {
		var zero T
		return zero, false, err
	}
	return l.GetIndex(index)
}

func (l *SparseLayer[T]) GetIndex(index CellIndex) (T, bool, error) {
	if err := l.validateIndex(index); err != nil {
		var zero T
		return zero, false, err
	}
	if l.values == nil {
		var zero T
		return zero, false, nil
	}
	value, ok := l.values[index]
	return value, ok, nil
}

func (l *SparseLayer[T]) Delete(pos Position) error {
	index, err := l.index(pos)
	if err != nil {
		return err
	}
	return l.DeleteIndex(index)
}

func (l *SparseLayer[T]) DeleteIndex(index CellIndex) error {
	if err := l.validateIndex(index); err != nil {
		return err
	}
	if l.values == nil {
		return nil
	}
	delete(l.values, index)
	return nil
}

func (l *SparseLayer[T]) Clear() {
	if l == nil || l.values == nil {
		return
	}
	clear(l.values)
}

// ForEach iterates written values only without creating copies.
func (l *SparseLayer[T]) ForEach(visit func(pos Position, value T) bool) {
	if l == nil || l.grid == nil || l.values == nil || visit == nil {
		return
	}

	for index, value := range l.values {
		pos, ok := l.grid.Position(index)
		if !ok {
			continue
		}
		if !visit(pos, value) {
			return
		}
	}
}

func (l *SparseLayer[T]) index(pos Position) (CellIndex, error) {
	if l == nil || l.grid == nil {
		return 0, ErrNilGrid
	}
	return l.grid.Index(pos)
}

func (l *SparseLayer[T]) validateIndex(index CellIndex) error {
	if l == nil || l.grid == nil {
		return ErrNilGrid
	}
	i := int(index)
	if i < 0 || i >= l.grid.CellCount() {
		return fmt.Errorf("%w: %d", ErrOutOfBounds, index)
	}
	return nil
}

// LayerSpace holds typed sparse layers for one grid.
type LayerSpace[K comparable, T comparable] struct {
	grid   *Grid
	layers map[K]*SparseLayer[T]
}

func NewLayerSpace[K comparable, T comparable](grid *Grid) (*LayerSpace[K, T], error) {
	if grid == nil {
		return nil, ErrNilGrid
	}
	return &LayerSpace[K, T]{
		grid:   grid,
		layers: make(map[K]*SparseLayer[T]),
	}, nil
}

func (s *LayerSpace[K, T]) Grid() *Grid {
	if s == nil {
		return nil
	}
	return s.grid
}

func (s *LayerSpace[K, T]) Create(key K) (*SparseLayer[T], error) {
	if s == nil || s.grid == nil {
		return nil, ErrNilGrid
	}
	if _, exists := s.layers[key]; exists {
		return nil, fmt.Errorf("%w: %v", ErrLayerExists, key)
	}
	layer, err := NewSparseLayer[T](s.grid)
	if err != nil {
		return nil, err
	}
	s.layers[key] = layer
	return layer, nil
}

func (s *LayerSpace[K, T]) Get(key K) (*SparseLayer[T], bool) {
	if s == nil {
		return nil, false
	}
	layer, ok := s.layers[key]
	return layer, ok
}

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

func (s *LayerSpace[K, T]) Count() int {
	if s == nil {
		return 0
	}
	return len(s.layers)
}

// ForEach iterates layer references without extra allocation.
func (s *LayerSpace[K, T]) ForEach(visit func(key K, layer *SparseLayer[T]) bool) {
	if s == nil || visit == nil {
		return
	}
	for key, layer := range s.layers {
		if !visit(key, layer) {
			return
		}
	}
}
