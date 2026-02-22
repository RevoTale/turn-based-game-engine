package events

// Emit queues one internal follow-up event in the active command tree.
type Emit[S any, P any] func(Event[S, P])

// Context carries typed state, patch, and event emitter for one command tree.
type Context[S any, P any] struct {
	state S
	patch *P
	emit  Emit[S, P]
}

// State returns the state view bound to this execution tree.
func (c Context[S, P]) State() S {
	return c.state
}

// Patch returns the mutable patch bound to this execution tree.
func (c Context[S, P]) Patch() *P {
	return c.patch
}

// Emit queues one follow-up event.
func (c Context[S, P]) Emit(event Event[S, P]) {
	c.emit(event)
}

// CommandHandler processes one root command input.
type CommandHandler[S any, P any, In any] func(ctx Context[S, P], input In) error

// EventHandler processes one internal event step.
type EventHandler[S any, P any] func(ctx Context[S, P]) error

// Event is one internal event actor.
type Event[S any, P any] struct {
	run EventHandler[S, P]
}

// Command is one root command actor.
type Command[S any, P any, In any] struct {
	run CommandHandler[S, P, In]
}

// DefineEvent creates one typed internal event actor.
func DefineEvent[S any, P any](handler EventHandler[S, P]) (Event[S, P], error) {
	if handler == nil {
		return Event[S, P]{}, ErrNilHandler
	}

	return Event[S, P]{
		run: handler,
	}, nil
}

// DefineCommand creates one typed root command actor.
func DefineCommand[S any, P any, In any](handler CommandHandler[S, P, In]) (Command[S, P, In], error) {
	if handler == nil {
		return Command[S, P, In]{}, ErrNilHandler
	}

	return Command[S, P, In]{
		run: handler,
	}, nil
}
