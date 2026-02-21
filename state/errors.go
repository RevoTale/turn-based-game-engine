package state

import "errors"

var (
	// ErrNilStore means a nil store pointer was used.
	ErrNilStore = errors.New("state store is nil")
	// ErrNilTransaction means a nil transactional callback was provided.
	ErrNilTransaction = errors.New("state transaction callback is nil")
	// ErrNilRollback means a nil rollback hook was provided.
	ErrNilRollback = errors.New("state rollback hook is nil")
)
