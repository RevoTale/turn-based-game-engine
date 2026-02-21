package grid2d

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type cellPair[T comparable] struct {
	pos   Position
	value T
}

// LockedSparseLayer wraps SparseLayer with one global read/write lock.
//
// Use this when you need thread-safe layer access with simple lock semantics.
type LockedSparseLayer[T comparable] struct {
	mu    sync.RWMutex
	layer *SparseLayer[T]
}

// NewLockedSparseLayer creates a globally locked sparse layer.
func NewLockedSparseLayer[T comparable](grid *Grid) (*LockedSparseLayer[T], error) {
	layer, err := NewSparseLayer[T](grid)
	if err != nil {
		return nil, err
	}
	return &LockedSparseLayer[T]{
		layer: layer,
	}, nil
}

// Grid returns the grid used by this layer.
func (l *LockedSparseLayer[T]) Grid() *Grid {
	if l == nil {
		return nil
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.layer == nil {
		return nil
	}
	return l.layer.grid
}

// Len returns how many cells currently have stored values.
func (l *LockedSparseLayer[T]) Len() int {
	if l == nil {
		return 0
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.layer == nil {
		return 0
	}
	return l.layer.Len()
}

// Set writes value at coordinates.
func (l *LockedSparseLayer[T]) Set(pos Position, value T) error {
	if l == nil {
		return ErrNilGrid
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.layer == nil {
		return ErrNilGrid
	}
	return l.layer.Set(pos, value)
}

// SetIndex writes value by CellIndex.
func (l *LockedSparseLayer[T]) SetIndex(index CellIndex, value T) error {
	if l == nil {
		return ErrNilGrid
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.layer == nil {
		return ErrNilGrid
	}
	return l.layer.SetIndex(index, value)
}

// Get reads value by coordinates.
func (l *LockedSparseLayer[T]) Get(pos Position) (T, bool, error) {
	if l == nil {
		var zero T
		return zero, false, ErrNilGrid
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.layer == nil {
		var zero T
		return zero, false, ErrNilGrid
	}
	return l.layer.Get(pos)
}

// GetIndex reads value by CellIndex.
func (l *LockedSparseLayer[T]) GetIndex(index CellIndex) (T, bool, error) {
	if l == nil {
		var zero T
		return zero, false, ErrNilGrid
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.layer == nil {
		var zero T
		return zero, false, ErrNilGrid
	}
	return l.layer.GetIndex(index)
}

// Delete removes value by coordinates.
func (l *LockedSparseLayer[T]) Delete(pos Position) error {
	if l == nil {
		return ErrNilGrid
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.layer == nil {
		return ErrNilGrid
	}
	return l.layer.Delete(pos)
}

// DeleteIndex removes value by CellIndex.
func (l *LockedSparseLayer[T]) DeleteIndex(index CellIndex) error {
	if l == nil {
		return ErrNilGrid
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.layer == nil {
		return ErrNilGrid
	}
	return l.layer.DeleteIndex(index)
}

// Clear removes all stored values from the layer.
func (l *LockedSparseLayer[T]) Clear() {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.layer == nil {
		return
	}
	l.layer.Clear()
}

// ForEach calls visit for each stored cell value.
//
// It takes a snapshot under read lock and then invokes visit without lock.
func (l *LockedSparseLayer[T]) ForEach(visit func(pos Position, value T) bool) {
	if l == nil || visit == nil {
		return
	}

	pairs := l.snapshot()
	for _, pair := range pairs {
		if !visit(pair.pos, pair.value) {
			return
		}
	}
}

func (l *LockedSparseLayer[T]) snapshot() []cellPair[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.layer == nil || l.layer.grid == nil || l.layer.values == nil {
		return nil
	}

	out := make([]cellPair[T], 0, len(l.layer.values))
	for index, value := range l.layer.values {
		pos, ok := l.layer.grid.Position(index)
		if !ok {
			continue
		}
		out = append(out, cellPair[T]{
			pos:   pos,
			value: value,
		})
	}
	return out
}

type rowBucket[T comparable] struct {
	mu     sync.RWMutex
	values map[int]T
}

// RowLockedSparseLayer stores one sparse map per row with independent locks.
//
// Use this when concurrent access is typically distributed across rows.
type RowLockedSparseLayer[T comparable] struct {
	grid  *Grid
	rows  []rowBucket[T]
	count atomic.Int64
}

// NewRowLockedSparseLayer creates a per-row locked sparse layer.
func NewRowLockedSparseLayer[T comparable](grid *Grid) (*RowLockedSparseLayer[T], error) {
	if grid == nil {
		return nil, ErrNilGrid
	}
	return &RowLockedSparseLayer[T]{
		grid: grid,
		rows: make([]rowBucket[T], grid.Height()),
	}, nil
}

// Grid returns the grid used by this layer.
func (l *RowLockedSparseLayer[T]) Grid() *Grid {
	if l == nil {
		return nil
	}
	return l.grid
}

// Len returns how many cells currently have stored values.
func (l *RowLockedSparseLayer[T]) Len() int {
	if l == nil {
		return 0
	}
	return int(l.count.Load())
}

// Set writes value at coordinates.
func (l *RowLockedSparseLayer[T]) Set(pos Position, value T) error {
	if err := l.validatePos(pos); err != nil {
		return err
	}

	row := &l.rows[pos.Y]
	row.mu.Lock()
	defer row.mu.Unlock()

	if row.values == nil {
		row.values = make(map[int]T)
	}
	if _, exists := row.values[pos.X]; !exists {
		l.count.Add(1)
	}
	row.values[pos.X] = value
	return nil
}

// SetIndex writes value by CellIndex.
func (l *RowLockedSparseLayer[T]) SetIndex(index CellIndex, value T) error {
	pos, err := l.position(index)
	if err != nil {
		return err
	}
	return l.Set(pos, value)
}

// Get reads value by coordinates.
func (l *RowLockedSparseLayer[T]) Get(pos Position) (T, bool, error) {
	if err := l.validatePos(pos); err != nil {
		var zero T
		return zero, false, err
	}

	row := &l.rows[pos.Y]
	row.mu.RLock()
	defer row.mu.RUnlock()
	if row.values == nil {
		var zero T
		return zero, false, nil
	}
	value, ok := row.values[pos.X]
	return value, ok, nil
}

// GetIndex reads value by CellIndex.
func (l *RowLockedSparseLayer[T]) GetIndex(index CellIndex) (T, bool, error) {
	pos, err := l.position(index)
	if err != nil {
		var zero T
		return zero, false, err
	}
	return l.Get(pos)
}

// Delete removes value by coordinates.
func (l *RowLockedSparseLayer[T]) Delete(pos Position) error {
	if err := l.validatePos(pos); err != nil {
		return err
	}

	row := &l.rows[pos.Y]
	row.mu.Lock()
	defer row.mu.Unlock()
	if row.values == nil {
		return nil
	}
	if _, exists := row.values[pos.X]; exists {
		delete(row.values, pos.X)
		l.count.Add(-1)
	}
	return nil
}

// DeleteIndex removes value by CellIndex.
func (l *RowLockedSparseLayer[T]) DeleteIndex(index CellIndex) error {
	pos, err := l.position(index)
	if err != nil {
		return err
	}
	return l.Delete(pos)
}

// Clear removes all stored values from the layer.
func (l *RowLockedSparseLayer[T]) Clear() {
	if l == nil {
		return
	}

	for i := range l.rows {
		l.rows[i].mu.Lock()
	}
	for i := len(l.rows) - 1; i >= 0; i-- {
		row := &l.rows[i]
		if row.values != nil {
			clear(row.values)
		}
		row.mu.Unlock()
	}
	l.count.Store(0)
}

// ForEach calls visit for each stored cell value.
//
// It iterates row snapshots to avoid holding row locks while visit executes.
func (l *RowLockedSparseLayer[T]) ForEach(visit func(pos Position, value T) bool) {
	if l == nil || l.grid == nil || visit == nil {
		return
	}

	for y := range l.rows {
		cells := l.snapshotRow(y)
		for _, pair := range cells {
			if !visit(pair.pos, pair.value) {
				return
			}
		}
	}
}

func (l *RowLockedSparseLayer[T]) snapshotRow(y int) []cellPair[T] {
	row := &l.rows[y]
	row.mu.RLock()
	defer row.mu.RUnlock()
	if row.values == nil {
		return nil
	}

	out := make([]cellPair[T], 0, len(row.values))
	for x, value := range row.values {
		out = append(out, cellPair[T]{
			pos:   Position{X: x, Y: y},
			value: value,
		})
	}
	return out
}

func (l *RowLockedSparseLayer[T]) validatePos(pos Position) error {
	if l == nil || l.grid == nil {
		return ErrNilGrid
	}
	if !l.grid.InBounds(pos) {
		return fmt.Errorf("%w: (%d,%d)", ErrOutOfBounds, pos.X, pos.Y)
	}
	return nil
}

func (l *RowLockedSparseLayer[T]) position(index CellIndex) (Position, error) {
	if l == nil || l.grid == nil {
		return Position{}, ErrNilGrid
	}
	pos, ok := l.grid.Position(index)
	if !ok {
		return Position{}, fmt.Errorf("%w: %d", ErrOutOfBounds, index)
	}
	return pos, nil
}
