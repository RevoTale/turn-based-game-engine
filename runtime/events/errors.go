package events

import "errors"

var (
	// ErrNilRuntime means a nil runtime pointer was used.
	ErrNilRuntime = errors.New("events runtime is nil")
	// ErrNilPatchFactory means ExecuteCommand was called without a patch factory.
	ErrNilPatchFactory = errors.New("events patch factory is nil")
	// ErrNilPatch means patch factory returned nil.
	ErrNilPatch = errors.New("events patch is nil")
	// ErrNilCommand means command definition is empty.
	ErrNilCommand = errors.New("events command is nil")
	// ErrNilEvent means emitted event definition is empty.
	ErrNilEvent = errors.New("events event is nil")
	// ErrNilHandler means command/event handler is nil.
	ErrNilHandler = errors.New("events handler is nil")
)
