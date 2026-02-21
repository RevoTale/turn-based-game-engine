package state

import "sync"

// Store holds one authoritative typed state and transactional version counter.
type Store[S any] struct {
	mu      sync.Mutex
	state   S
	version uint64
}

// Tx is one in-progress transactional mutation scope.
type Tx[S any] struct {
	state     *S
	version   uint64
	rollbacks []func(*S)
}

// New creates a store initialized with the provided state.
func New[S any](initial S) *Store[S] {
	return &Store[S]{
		state: initial,
	}
}

// Version returns the current committed version.
func (s *Store[S]) Version() (uint64, error) {
	if s == nil {
		return 0, ErrNilStore
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	return s.version, nil
}

// View executes one read callback under lock.
//
// The callback receives a pointer for performance. The callback must treat it
// as read-only.
func (s *Store[S]) View(read func(state *S, version uint64)) error {
	if s == nil {
		return ErrNilStore
	}
	if read == nil {
		return ErrNilTransaction
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	read(&s.state, s.version)
	return nil
}

// Do executes one transactional mutation callback under lock.
//
// The callback can register rollback hooks with Tx.BeforeChange before each
// in-place mutation. If callback returns error, rollback hooks are executed in
// reverse order and the version is not incremented.
func (s *Store[S]) Do(run func(tx *Tx[S]) error) (uint64, error) {
	if s == nil {
		return 0, ErrNilStore
	}
	if run == nil {
		return 0, ErrNilTransaction
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	tx := Tx[S]{
		state:     &s.state,
		version:   s.version,
		rollbacks: make([]func(*S), 0, 8),
	}

	if err := run(&tx); err != nil {
		tx.rollback()
		return s.version, err
	}

	s.version++
	return s.version, nil
}

// State returns mutable state pointer bound to current transaction.
func (tx *Tx[S]) State() *S {
	if tx == nil {
		return nil
	}
	return tx.state
}

// Version returns committed version observed at transaction start.
func (tx *Tx[S]) Version() uint64 {
	if tx == nil {
		return 0
	}
	return tx.version
}

// BeforeChange registers one rollback hook.
//
// Hooks run in reverse order when transaction callback returns error.
func (tx *Tx[S]) BeforeChange(undo func(state *S)) error {
	if tx == nil || tx.state == nil {
		return ErrNilTransaction
	}
	if undo == nil {
		return ErrNilRollback
	}
	tx.rollbacks = append(tx.rollbacks, undo)
	return nil
}

func (tx *Tx[S]) rollback() {
	for i := len(tx.rollbacks) - 1; i >= 0; i-- {
		tx.rollbacks[i](tx.state)
	}
	tx.rollbacks = tx.rollbacks[:0]
}
