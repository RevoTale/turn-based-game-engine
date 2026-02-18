package automation

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScheduler_Run_ReturnsErrScopeBusyWhenScopeIsActive(t *testing.T) {
	scheduler := NewScheduler[int, int, int](SchedulerOptions{})

	started := make(chan struct{})
	unblock := make(chan struct{})
	startOnce := sync.Once{}

	startedAsync := scheduler.Trigger(1, SessionConfig[int, int, int]{
		MaxIterations: 1,
		ForEachActor: func(scope int, visit func(actor int) bool) {
			visit(1)
		},
		ProcessActor: func(scope int, actor int, emit func(result int)) (bool, bool, error) {
			startOnce.Do(func() { close(started) })
			<-unblock
			return false, false, nil
		},
	}, nil)
	require.True(t, startedAsync)

	<-started

	_, err := scheduler.Run(1, SessionConfig[int, int, int]{
		MaxIterations: 1,
		ForEachActor: func(scope int, visit func(actor int) bool) {
			visit(1)
		},
		ProcessActor: func(scope int, actor int, emit func(result int)) (bool, bool, error) {
			return false, false, nil
		},
	})
	require.ErrorIs(t, err, ErrScopeBusy)

	close(unblock)
	scheduler.Wait()
}

func TestScheduler_Trigger_EmitsResultsAndRunsDoneCallback(t *testing.T) {
	scheduler := NewScheduler[string, int, int](SchedulerOptions{})

	done := make(chan struct{})
	var got []int
	var gotErr error

	ok := scheduler.Trigger("room", SessionConfig[string, int, int]{
		MaxIterations: 2,
		ForEachActor: func(scope string, visit func(actor int) bool) {
			visit(1)
		},
		ProcessActor: func(scope string, actor int, emit func(result int)) (bool, bool, error) {
			emit(actor)
			return true, true, nil
		},
	}, func(results []int, err error) {
		got = results
		gotErr = err
		close(done)
	})
	require.True(t, ok)

	<-done
	scheduler.Wait()

	require.NoError(t, gotErr)
	assert.Equal(t, []int{1}, got)
}

func TestScheduler_Run_AppliesDelayHookViaSleepOption(t *testing.T) {
	var sleeps []time.Duration
	scheduler := NewScheduler[int, int, int](SchedulerOptions{
		Sleep: func(d time.Duration) {
			sleeps = append(sleeps, d)
		},
	})

	results, err := scheduler.Run(1, SessionConfig[int, int, int]{
		MaxIterations: 1,
		ForEachActor: func(scope int, visit func(actor int) bool) {
			visit(1)
		},
		DelayForActor: func(scope int, actor int) time.Duration {
			return 5 * time.Millisecond
		},
		ProcessActor: func(scope int, actor int, emit func(result int)) (bool, bool, error) {
			emit(actor)
			return true, true, nil
		},
	})

	require.NoError(t, err)
	assert.Equal(t, []int{1}, results)
	assert.Equal(t, []time.Duration{5 * time.Millisecond}, sleeps)
}

func TestScheduler_Run_PropagatesMaxIterationsError(t *testing.T) {
	scheduler := NewScheduler[int, int, int](SchedulerOptions{})

	results, err := scheduler.Run(7, SessionConfig[int, int, int]{
		MaxIterations: 3,
		ForEachActor: func(scope int, visit func(actor int) bool) {
			visit(1)
		},
		ProcessActor: func(scope int, actor int, emit func(result int)) (bool, bool, error) {
			emit(actor)
			return true, false, nil
		},
	})

	require.ErrorIs(t, err, ErrMaxIterationsReached)
	assert.Equal(t, []int{1, 1, 1}, results)
}
