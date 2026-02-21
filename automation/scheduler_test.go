package automation

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScheduler_Run_ReturnsErrScopeBusyWhenSessionIsActive(t *testing.T) {
	scheduler := NewScheduler[int, int, int](SchedulerOptions{})

	started := make(chan struct{})
	unblock := make(chan struct{})
	done := make(chan struct{})
	runErr := make(chan error, 1)
	startOnce := sync.Once{}

	go func() {
		_, err := scheduler.Run(1, SessionConfig[int, int, int]{
			MaxIterations: 1,
			ForEachActor: func(scope int, visit func(actor int) bool) {
				visit(1)
			},
			ProcessActor: func(scope int, actor int, emit func(result int)) (bool, bool, error) {
				startOnce.Do(func() { close(started) })
				<-unblock
				return false, false, nil
			},
		})
		runErr <- err
		close(done)
	}()

	<-started

	_, err := scheduler.Run(2, SessionConfig[int, int, int]{
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
	<-done
	require.NoError(t, <-runErr)
}

func TestScheduler_Run_ReportsResults(t *testing.T) {
	scheduler := NewScheduler[string, int, int](SchedulerOptions{})

	got, gotErr := scheduler.Run("room", SessionConfig[string, int, int]{
		MaxIterations: 2,
		ForEachActor: func(scope string, visit func(actor int) bool) {
			visit(1)
		},
		ProcessActor: func(scope string, actor int, emit func(result int)) (bool, bool, error) {
			emit(actor)
			return true, true, nil
		},
	})

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
