package grid2d

import "fmt"

// SparseLayer stores values only for cells that were written.
//
// Use it to avoid allocating full board-sized arrays for sparse game state.
type SparseLayer[T comparable] struct {
	grid   *Grid
	values map[CellIndex]T
}

// NewSparseLayer creates an empty layer for one grid.
//
// Internal storage map is created on first write.
// Returns ErrNilGrid when grid is nil.
func NewSparseLayer[T comparable](grid *Grid) (*SparseLayer[T], error) {
	if grid == nil {
		return nil, ErrNilGrid
	}
	return &SparseLayer[T]{grid: grid}, nil
}

// Grid returns the grid used by this layer.
func (l *SparseLayer[T]) Grid() *Grid {
	if l == nil {
		return nil
	}
	return l.grid
}

// Len returns how many cells currently have stored values.
func (l *SparseLayer[T]) Len() int {
	if l == nil || l.values == nil {
		return 0
	}
	return len(l.values)
}

// Set writes value at coordinates.
//
// Returns ErrNilGrid for nil layer/grid and ErrOutOfBounds for invalid position.
func (l *SparseLayer[T]) Set(pos Position, value T) error {
	index, err := l.index(pos)
	if err != nil {
		return err
	}
	return l.SetIndex(index, value)
}

// SetIndex writes value by CellIndex.
//
// Returns ErrNilGrid for nil layer/grid and ErrOutOfBounds for invalid index.
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

// Get reads value by coordinates.
//
// It returns (zero, false, nil) when cell has no value.
// Returns ErrNilGrid for nil layer/grid and ErrOutOfBounds for invalid position.
func (l *SparseLayer[T]) Get(pos Position) (T, bool, error) {
	index, err := l.index(pos)
	if err != nil {
		var zero T
		return zero, false, err
	}
	return l.GetIndex(index)
}

// GetIndex reads value by CellIndex.
//
// It returns (zero, false, nil) when cell has no value.
// Returns ErrNilGrid for nil layer/grid and ErrOutOfBounds for invalid index.
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

// Delete removes value by coordinates.
//
// Returns ErrNilGrid for nil layer/grid and ErrOutOfBounds for invalid position.
func (l *SparseLayer[T]) Delete(pos Position) error {
	index, err := l.index(pos)
	if err != nil {
		return err
	}
	return l.DeleteIndex(index)
}

// DeleteIndex removes value by CellIndex.
//
// Deleting a missing value is a no-op.
// Returns ErrNilGrid for nil layer/grid and ErrOutOfBounds for invalid index.
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

// Clear removes all stored values from the layer.
func (l *SparseLayer[T]) Clear() {
	if l == nil || l.values == nil {
		return
	}
	clear(l.values)
}

// ForEach calls visit for each stored cell value.
//
// It only visits written cells and does not create a copy.
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
