package events_test

import (
	"errors"
	"testing"
	"time"

	"github.com/RevoTale/turn-based-game-engine/runtime/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilderRegisterAndSeal(t *testing.T) {
	t.Parallel()

	builder := events.NewBuilder()

	first, err := events.RegisterEvent(builder, func(_ *events.Context, _ int) error {
		return nil
	}, events.Hooks[int]{})
	require.NoError(t, err)

	second, err := events.RegisterCommand(builder, func(_ *events.Context, _ int) error {
		return nil
	}, events.Hooks[int]{})
	require.NoError(t, err)
	require.NotEqual(t, first, events.Event[int]{})
	require.NotEqual(t, second, events.Command[int]{})

	runtime, err := builder.Build()
	require.NoError(t, err)
	require.NotNil(t, runtime)

	_, err = events.RegisterEvent(builder, func(_ *events.Context, _ int) error {
		return nil
	}, events.Hooks[int]{})
	require.ErrorIs(t, err, events.ErrBuilderSealed)
}

func TestDispatcherDispatchDrainsEmittedEventTree(t *testing.T) {
	t.Parallel()

	type fireCommand struct {
		Coordinate string
		Hit        bool
	}
	type hitEvent struct {
		Coordinate string
	}
	type missEvent struct {
		Coordinate string
	}

	builder := events.NewBuilder()
	trace := make([]string, 0, 2)

	hit, err := events.RegisterEvent(builder, func(_ *events.Context, payload hitEvent) error {
		trace = append(trace, "hit:"+payload.Coordinate)
		return nil
	}, events.Hooks[hitEvent]{})
	require.NoError(t, err)

	miss, err := events.RegisterEvent(builder, func(_ *events.Context, payload missEvent) error {
		trace = append(trace, "miss:"+payload.Coordinate)
		return nil
	}, events.Hooks[missEvent]{})
	require.NoError(t, err)

	fire, err := events.RegisterCommand(builder, func(_ *events.Context, payload fireCommand) error {
		trace = append(trace, "fire:"+payload.Coordinate)
		return nil
	}, events.Hooks[fireCommand]{
		After: func(ctx *events.Context, payload fireCommand) error {
			if payload.Hit {
				ctx.Emit(events.Next(hit, hitEvent{Coordinate: payload.Coordinate}))
				return nil
			}
			ctx.Emit(events.Next(miss, missEvent{Coordinate: payload.Coordinate}))
			return nil
		},
	})
	require.NoError(t, err)

	runtime, err := builder.Build()
	require.NoError(t, err)

	err = events.ExecuteCommand(runtime, fire, fireCommand{Coordinate: "B4", Hit: true})
	require.NoError(t, err)
	assert.Equal(t, []string{"fire:B4", "hit:B4"}, trace)

	trace = trace[:0]
	err = events.ExecuteCommand(runtime, fire, fireCommand{Coordinate: "A1", Hit: false})
	require.NoError(t, err)
	assert.Equal(t, []string{"fire:A1", "miss:A1"}, trace)
}

func TestDispatcherOnFailRunsWhenHookFails(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("before failed")
	builder := events.NewBuilder()

	onFailCalled := false
	command, err := events.RegisterCommand(builder, func(_ *events.Context, _ int) error {
		t.Fatal("handler must not execute when before hook fails")
		return nil
	}, events.Hooks[int]{
		Before: func(_ *events.Context, _ int) error {
			return expectedErr
		},
		OnFail: func(_ *events.Context, _ int, cause error) {
			onFailCalled = true
			assert.ErrorIs(t, cause, expectedErr)
		},
	})
	require.NoError(t, err)

	runtime, err := builder.Build()
	require.NoError(t, err)

	err = events.ExecuteCommand(runtime, command, 10)
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.True(t, onFailCalled)
}

func TestRuntimeRejectsUnknownRootCommand(t *testing.T) {
	t.Parallel()

	builder := events.NewBuilder()
	runtime, err := builder.Build()
	require.NoError(t, err)

	var unknown events.Command[int]
	err = events.ExecuteCommand(runtime, unknown, 7)
	require.ErrorIs(t, err, events.ErrUnknownCommand)
}

func TestRuntimeRejectsUnknownReturnedEvent(t *testing.T) {
	t.Parallel()

	builder := events.NewBuilder()

	root, err := events.RegisterCommand(builder, func(ctx *events.Context, _ int) error {
		var unknown events.Event[int]
		ctx.Emit(events.Next(unknown, 1))
		return nil
	}, events.Hooks[int]{})
	require.NoError(t, err)

	runtime, err := builder.Build()
	require.NoError(t, err)

	err = events.ExecuteCommand(runtime, root, 0)
	require.ErrorIs(t, err, events.ErrUnknownEvent)
}

func TestExecuteCommandRunsCommand(t *testing.T) {
	t.Parallel()

	builder := events.NewBuilder()
	ran := false

	command, err := events.RegisterCommand(builder, func(_ *events.Context, payload int) error {
		ran = payload == 42
		return nil
	}, events.Hooks[int]{})
	require.NoError(t, err)

	runtime, err := builder.Build()
	require.NoError(t, err)

	require.NoError(t, events.ExecuteCommand(runtime, command, 42))
	require.True(t, ran)
}

func TestRuntimeSerializesConcurrentCommandTrees(t *testing.T) {
	t.Parallel()

	builder := events.NewBuilder()
	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})
	secondEntered := make(chan struct{})

	command, err := events.RegisterCommand(builder, func(_ *events.Context, payload int) error {
		if payload == 1 {
			close(firstStarted)
			<-releaseFirst
			return nil
		}
		if payload == 2 {
			close(secondEntered)
		}
		return nil
	}, events.Hooks[int]{})
	require.NoError(t, err)

	runtime, err := builder.Build()
	require.NoError(t, err)

	errCh := make(chan error, 2)
	go func() {
		errCh <- events.ExecuteCommand(runtime, command, 1)
	}()
	<-firstStarted

	go func() {
		errCh <- events.ExecuteCommand(runtime, command, 2)
	}()

	select {
	case <-secondEntered:
		t.Fatal("second command entered before first completed")
	case <-time.After(20 * time.Millisecond):
	}

	close(releaseFirst)

	require.NoError(t, <-errCh)
	require.NoError(t, <-errCh)

	select {
	case <-secondEntered:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("second command did not execute after first completion")
	}
}
