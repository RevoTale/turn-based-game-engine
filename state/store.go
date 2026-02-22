package state

import "sync"

// Store holds one authoritative typed state and committed version counter.
type Store[S any] struct {
	mu      sync.Mutex
	state   S
	version uint64
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

// Do executes one mutation callback under lock and increments version when the
// callback succeeds.
func (s *Store[S]) Do(run func(state *S, version uint64) error) (uint64, error) {
	if s == nil {
		return 0, ErrNilStore
	}
	if run == nil {
		return 0, ErrNilTransaction
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := run(&s.state, s.version); err != nil {
		return s.version, err
	}
	s.version++
	return s.version, nil
}
