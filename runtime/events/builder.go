package events

import "maps"

import "fmt"

// EventID is the internal identifier assigned at registration time.
type EventID uint32

// Event is a typed token for one registered event.
type Event[T any] struct {
	id EventID
}

// Command is a typed token for one registered root command.
type Command[T any] struct {
	id EventID
}

// Handler executes one event body and can queue follow-up emissions via
// Context.Emit.
type Handler[T any] func(ctx *Context, payload T) error

// Hook executes a pre or post stage around event handling.
type Hook[T any] func(ctx *Context, payload T) error

// FailHook executes when a before, handler, or after stage fails.
type FailHook[T any] func(ctx *Context, payload T, cause error)

// Hooks groups optional lifecycle hooks for one event.
type Hooks[T any] struct {
	Before Hook[T]
	After  Hook[T]
	OnFail FailHook[T]
}

// Builder registers events and seals configuration with Build.
type Builder struct {
	nextID EventID
	defs   map[EventID]definition
	sealed bool
}

// NewBuilder creates an empty event builder.
func NewBuilder() *Builder {
	return &Builder{
		nextID: 1,
		defs:   make(map[EventID]definition),
	}
}

func (b *Builder) registerID() (EventID, error) {
	if b == nil {
		return 0, ErrNilBuilder
	}
	if b.sealed {
		return 0, ErrBuilderSealed
	}

	id := b.nextID
	b.nextID++
	return id, nil
}

// Build seals the builder and returns an immutable runtime.
func (b *Builder) Build() (*Runtime, error) {
	if b == nil {
		return nil, ErrNilBuilder
	}

	defs := make(map[EventID]definition, len(b.defs))
	maps.Copy(defs, b.defs)

	b.sealed = true
	return &Runtime{
		dispatcher: dispatcher{defs: defs},
	}, nil
}

// RegisterEvent adds one typed event definition and returns its typed token.
func RegisterEvent[T any](b *Builder, handler Handler[T], hooks Hooks[T]) (Event[T], error) {
	if handler == nil {
		return Event[T]{}, ErrNilHandler
	}

	id, err := b.registerID()
	if err != nil {
		return Event[T]{}, err
	}
	b.defs[id] = buildDefinition(id, handler, hooks)
	return Event[T]{id: id}, nil
}

// RegisterCommand adds one typed root command definition and returns its token.
func RegisterCommand[T any](b *Builder, handler Handler[T], hooks Hooks[T]) (Command[T], error) {
	if handler == nil {
		return Command[T]{}, ErrNilHandler
	}

	id, err := b.registerID()
	if err != nil {
		return Command[T]{}, err
	}
	b.defs[id] = buildDefinition(id, handler, hooks)
	return Command[T]{id: id}, nil
}

type definition struct {
	run func(ctx *Context, payload any) error
}

func buildDefinition[T any](id EventID, handler Handler[T], hooks Hooks[T]) definition {
	return definition{
		run: func(ctx *Context, payload any) error {
			typedPayload, ok := payload.(T)
			if !ok {
				return fmt.Errorf("%w: id=%d", ErrPayloadTypeMismatch, id)
			}

			stageMark := ctx.beginEmissionScope()

			if hooks.Before != nil {
				err := hooks.Before(ctx, typedPayload)
				if err != nil {
					ctx.rollbackEmissionScope(stageMark)
					if hooks.OnFail != nil {
						hooks.OnFail(ctx, typedPayload, err)
					}
					return err
				}
				if emitErr := ctx.consumeEmitError(); emitErr != nil {
					ctx.rollbackEmissionScope(stageMark)
					if hooks.OnFail != nil {
						hooks.OnFail(ctx, typedPayload, emitErr)
					}
					return emitErr
				}
			}

			err := handler(ctx, typedPayload)
			if err != nil {
				ctx.rollbackEmissionScope(stageMark)
				if hooks.OnFail != nil {
					hooks.OnFail(ctx, typedPayload, err)
				}
				return err
			}
			if emitErr := ctx.consumeEmitError(); emitErr != nil {
				ctx.rollbackEmissionScope(stageMark)
				if hooks.OnFail != nil {
					hooks.OnFail(ctx, typedPayload, emitErr)
				}
				return emitErr
			}

			if hooks.After != nil {
				err = hooks.After(ctx, typedPayload)
				if err != nil {
					ctx.rollbackEmissionScope(stageMark)
					if hooks.OnFail != nil {
						hooks.OnFail(ctx, typedPayload, err)
					}
					return err
				}
				if emitErr := ctx.consumeEmitError(); emitErr != nil {
					ctx.rollbackEmissionScope(stageMark)
					if hooks.OnFail != nil {
						hooks.OnFail(ctx, typedPayload, emitErr)
					}
					return emitErr
				}
			}

			return nil
		},
	}
}
