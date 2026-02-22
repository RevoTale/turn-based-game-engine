# Runtime Events: Target Interfaces

This file defines the exact API shape for the refactor.

## Core Types

```go
type Runtime struct {
	// serializes command trees per runtime instance
}

type Context[S any, P any] struct {
	// carries typed state, patch, and event emission for one tree
}

type Event[S any, P any] struct {
	// typed internal step token
}

type Command[S any, P any, In any] struct {
	// typed root command token
}
```

## Registration

```go
type Emit[S any, P any] func(Event[S, P])

type CommandHandler[S any, P any, In any] func(ctx Context[S, P], input In) error
type EventHandler[S any, P any] func(ctx Context[S, P]) error

func DefineCommand[S any, P any, In any](
	handler CommandHandler[S, P, In],
) (Command[S, P, In], error)

func DefineEvent[S any, P any](
	handler EventHandler[S, P],
) (Event[S, P], error)
```

## Execution

```go
func NewRuntime() *Runtime

func ExecuteCommand[S any, P any, In any](
	runtime *Runtime,
	state S,
	command Command[S, P, In],
	input In,
	newPatch func() *P,
) (*P, error)
```

## Behavioral Contract

1. Events do not receive payload.
2. Command input is consumed only at root command handler.
3. Events and commands mutate only patch through `ctx.Patch()`.
4. State visibility is controlled by state type `S` used in `Context[S,P]` (use read-only interfaces).
5. Runtime executes one command tree at a time per runtime instance.
6. `ExecuteCommand` never commits store directly; it returns patch to caller.
