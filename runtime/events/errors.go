package events

import "errors"

var (
	// ErrNilBuilder means a nil builder was used.
	ErrNilBuilder = errors.New("events builder is nil")
	// ErrBuilderSealed means registration was attempted after Build.
	ErrBuilderSealed = errors.New("events builder is sealed")
	// ErrNilHandler means registration was called without a handler.
	ErrNilHandler = errors.New("event handler is required")
	// ErrNilRuntime means execution was attempted on a nil runtime.
	ErrNilRuntime = errors.New("events runtime is nil")
	// ErrNilDispatcher is kept as compatibility alias for ErrNilRuntime.
	ErrNilDispatcher = ErrNilRuntime
	// ErrUnknownCommand means a referenced command id is not registered.
	ErrUnknownCommand = errors.New("unknown command")
	// ErrUnknownEvent means a referenced event id is not registered.
	ErrUnknownEvent = errors.New("unknown event")
	// ErrPayloadTypeMismatch means runtime payload does not match event type.
	ErrPayloadTypeMismatch = errors.New("event payload type mismatch")
)
