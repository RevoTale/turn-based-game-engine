package events_test

import (
	"errors"
	"testing"
	"time"

	"github.com/RevoTale/turn-based-game-engine/runtime/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefineRejectsNilHandler(t *testing.T) {
	t.Parallel()

	_, err := events.DefineEvent[int, int](nil)
	require.ErrorIs(t, err, events.ErrNilHandler)

	_, err = events.DefineCommand[int, int, int](nil)
	require.ErrorIs(t, err, events.ErrNilHandler)
}

func TestExecuteCommandDrainsEmittedEventTree(t *testing.T) {
	t.Parallel()

	type state struct{}
	type patch struct {
		coord string
		hit   bool
		trace []string
	}
	type input struct {
		Coordinate string
		Hit        bool
	}

	var hitEv events.Event[*state, patch]
	var missEv events.Event[*state, patch]

	resolveEv, err := events.DefineEvent[*state, patch](func(ctx events.Context[*state, patch]) error {
		p := ctx.Patch()
		if p.hit {
			ctx.Emit(hitEv)
			return nil
		}
		ctx.Emit(missEv)
		return nil
	})
	require.NoError(t, err)

	hitEv, err = events.DefineEvent[*state, patch](func(ctx events.Context[*state, patch]) error {
		p := ctx.Patch()
		p.trace = append(p.trace, "hit:"+p.coord)
		return nil
	})
	require.NoError(t, err)

	missEv, err = events.DefineEvent[*state, patch](func(ctx events.Context[*state, patch]) error {
		p := ctx.Patch()
		p.trace = append(p.trace, "miss:"+p.coord)
		return nil
	})
	require.NoError(t, err)

	fireCmd, err := events.DefineCommand[*state, patch, input](func(ctx events.Context[*state, patch], in input) error {
		p := ctx.Patch()
		p.coord = in.Coordinate
		p.hit = in.Hit
		p.trace = append(p.trace, "fire:"+in.Coordinate)
		ctx.Emit(resolveEv)
		return nil
	})
	require.NoError(t, err)

	runtime := events.NewRuntime()
	s := &state{}

	firstPatch, err := events.ExecuteCommand(runtime, s, fireCmd, input{Coordinate: "B4", Hit: true}, func() *patch {
		return &patch{trace: make([]string, 0, 4)}
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"fire:B4", "hit:B4"}, firstPatch.trace)

	secondPatch, err := events.ExecuteCommand(runtime, s, fireCmd, input{Coordinate: "A1", Hit: false}, func() *patch {
		return &patch{trace: make([]string, 0, 4)}
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"fire:A1", "miss:A1"}, secondPatch.trace)
}

func TestRuntimeRejectsNilInputs(t *testing.T) {
	t.Parallel()

	type patch struct{}
	cmd, err := events.DefineCommand[int, patch, int](func(_ events.Context[int, patch], _ int) error { return nil })
	require.NoError(t, err)

	_, err = events.ExecuteCommand[int, patch, int](nil, 1, cmd, 1, func() *patch { return &patch{} })
	require.ErrorIs(t, err, events.ErrNilRuntime)

	runtime := events.NewRuntime()
	_, err = events.ExecuteCommand(runtime, 1, cmd, 1, nil)
	require.ErrorIs(t, err, events.ErrNilPatchFactory)

	_, err = events.ExecuteCommand(runtime, 1, cmd, 1, func() *patch { return nil })
	require.ErrorIs(t, err, events.ErrNilPatch)

	var nilCmd events.Command[int, patch, int]
	_, err = events.ExecuteCommand(runtime, 1, nilCmd, 1, func() *patch { return &patch{} })
	require.ErrorIs(t, err, events.ErrNilCommand)
}

func TestRuntimeRejectsNilEmittedEvent(t *testing.T) {
	t.Parallel()

	type patch struct{}
	cmd, err := events.DefineCommand[int, patch, int](func(ctx events.Context[int, patch], _ int) error {
		var unknown events.Event[int, patch]
		ctx.Emit(unknown)
		return nil
	})
	require.NoError(t, err)

	runtime := events.NewRuntime()
	_, err = events.ExecuteCommand(runtime, 1, cmd, 0, func() *patch { return &patch{} })
	require.ErrorIs(t, err, events.ErrNilEvent)
}

func TestRuntimeReturnsEventAndCommandErrors(t *testing.T) {
	t.Parallel()

	type patch struct{}
	expectedCmdErr := errors.New("cmd failed")
	expectedEvErr := errors.New("event failed")

	cmdErrCmd, err := events.DefineCommand[int, patch, int](func(_ events.Context[int, patch], _ int) error {
		return expectedCmdErr
	})
	require.NoError(t, err)

	runtime := events.NewRuntime()
	_, err = events.ExecuteCommand(runtime, 1, cmdErrCmd, 1, func() *patch { return &patch{} })
	require.ErrorIs(t, err, expectedCmdErr)

	failEv, err := events.DefineEvent[int, patch](func(_ events.Context[int, patch]) error {
		return expectedEvErr
	})
	require.NoError(t, err)

	cmdEmitFail, err := events.DefineCommand[int, patch, int](func(ctx events.Context[int, patch], _ int) error {
		ctx.Emit(failEv)
		return nil
	})
	require.NoError(t, err)

	_, err = events.ExecuteCommand(runtime, 1, cmdEmitFail, 1, func() *patch { return &patch{} })
	require.ErrorIs(t, err, expectedEvErr)
}

func TestRuntimeSerializesConcurrentCommandTrees(t *testing.T) {
	t.Parallel()

	type state struct {
		firstStarted chan struct{}
		releaseFirst chan struct{}
		secondEnter  chan struct{}
	}
	type patch struct{}

	s := &state{
		firstStarted: make(chan struct{}),
		releaseFirst: make(chan struct{}),
		secondEnter:  make(chan struct{}),
	}

	cmd, err := events.DefineCommand[*state, patch, int](func(ctx events.Context[*state, patch], input int) error {
		st := ctx.State()
		if input == 1 {
			close(st.firstStarted)
			<-st.releaseFirst
			return nil
		}
		if input == 2 {
			close(st.secondEnter)
		}
		return nil
	})
	require.NoError(t, err)

	runtime := events.NewRuntime()

	errCh := make(chan error, 2)
	go func() {
		_, err := events.ExecuteCommand(runtime, s, cmd, 1, func() *patch { return &patch{} })
		errCh <- err
	}()
	<-s.firstStarted

	go func() {
		_, err := events.ExecuteCommand(runtime, s, cmd, 2, func() *patch { return &patch{} })
		errCh <- err
	}()

	select {
	case <-s.secondEnter:
		t.Fatal("second command entered before first completed")
	case <-time.After(20 * time.Millisecond):
	}

	close(s.releaseFirst)
	require.NoError(t, <-errCh)
	require.NoError(t, <-errCh)

	select {
	case <-s.secondEnter:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("second command did not execute after first completion")
	}
}
