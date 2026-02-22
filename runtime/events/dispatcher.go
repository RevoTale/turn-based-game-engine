package events

import (
	"fmt"
	"sync"
)

// Runtime serializes command-tree execution for one match/session instance.
type Runtime struct {
	mu sync.Mutex
}

// NewRuntime creates a runtime that executes one command tree at a time.
func NewRuntime() *Runtime {
	return &Runtime{}
}

// ExecuteCommand executes one root command and drains all emitted events in
// deterministic FIFO order.
//
// ExecuteCommand never commits state. It returns the produced patch to caller.
func ExecuteCommand[S any, P any, In any](runtime *Runtime, state S, command Command[S, P, In], input In, newPatch func() *P) (*P, error) {
	if runtime == nil {
		return nil, ErrNilRuntime
	}
	if newPatch == nil {
		return nil, ErrNilPatchFactory
	}
	if command.run == nil {
		return nil, ErrNilCommand
	}

	runtime.mu.Lock()
	defer runtime.mu.Unlock()

	patch := newPatch()
	if patch == nil {
		return nil, ErrNilPatch
	}

	queue := make([]Event[S, P], 0, 8)
	var emitErr error
	emit := func(event Event[S, P]) {
		if emitErr != nil {
			return
		}
		if event.run == nil {
			emitErr = ErrNilEvent
			return
		}
		queue = append(queue, event)
	}

	ctx := Context[S, P]{
		state: state,
		patch: patch,
		emit:  emit,
	}

	if err := command.run(ctx, input); err != nil {
		return nil, fmt.Errorf("command failed: %w", err)
	}
	if emitErr != nil {
		return nil, emitErr
	}

	for i := 0; i < len(queue); i++ {
		if err := queue[i].run(ctx); err != nil {
			return nil, fmt.Errorf("event failed: %w", err)
		}
		if emitErr != nil {
			return nil, emitErr
		}
	}

	return patch, nil
}
