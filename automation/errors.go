package automation

import "errors"

var (
	// ErrInvalidMaxIterations means config.MaxIterations is not positive.
	ErrInvalidMaxIterations = errors.New("max iterations must be positive")
	// ErrNilForEachActor means config.ForEachActor is missing.
	ErrNilForEachActor = errors.New("for-each actor callback is required")
	// ErrNilProcessActor means config.ProcessActor is missing.
	ErrNilProcessActor = errors.New("process actor callback is required")
	// ErrMaxIterationsReached means loop budget was exhausted without stabilizing.
	ErrMaxIterationsReached = errors.New("max iterations reached")
	// ErrScopeBusy means the requested scope is already being processed.
	ErrScopeBusy = errors.New("automation scope is already being processed")
)
