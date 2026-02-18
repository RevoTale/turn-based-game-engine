package automation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_StopsWhenNoActorActs(t *testing.T) {
	actors := []string{"a", "b"}
	processed := 0

	results, err := Run(Config[string, int]{
		MaxIterations: 10,
		ForEachActor: func(visit func(actor string) bool) {
			for _, actor := range actors {
				if !visit(actor) {
					return
				}
			}
		},
		ProcessActor: func(actor string, emit func(result int)) (bool, bool, error) {
			processed++
			return false, false, nil
		},
	})

	require.NoError(t, err)
	assert.Empty(t, results)
	assert.Equal(t, 2, processed)
}

func TestRun_EmitsAndStopsOnSignal(t *testing.T) {
	actors := []int{1, 2, 3}
	emitted := make([]int, 0)

	results, err := Run(Config[int, int]{
		MaxIterations: 5,
		ForEachActor: func(visit func(actor int) bool) {
			for _, actor := range actors {
				if !visit(actor) {
					return
				}
			}
		},
		ProcessActor: func(actor int, emit func(result int)) (bool, bool, error) {
			emit(actor)
			return true, actor == 2, nil
		},
		OnEmit: func(result int) {
			emitted = append(emitted, result)
		},
	})

	require.NoError(t, err)
	assert.Equal(t, []int{1, 2}, results)
	assert.Equal(t, []int{1, 2}, emitted)
}

func TestRun_OnErrorDoesNotBreakLoop(t *testing.T) {
	actors := []string{"first", "second"}
	fail := errors.New("step failed")
	onErrorCalls := make([]string, 0)
	successCount := 0

	results, err := Run(Config[string, string]{
		MaxIterations: 2,
		ForEachActor: func(visit func(actor string) bool) {
			for _, actor := range actors {
				if !visit(actor) {
					return
				}
			}
		},
		ProcessActor: func(actor string, emit func(result string)) (bool, bool, error) {
			if actor == "first" {
				return false, false, fail
			}
			successCount++
			emit(actor)
			return true, false, nil
		},
		OnError: func(actor string, err error) {
			if err != nil {
				onErrorCalls = append(onErrorCalls, actor)
			}
		},
	})

	require.ErrorIs(t, err, ErrMaxIterationsReached)
	assert.Equal(t, []string{"second", "second"}, results)
	assert.Equal(t, 2, successCount)
	assert.Equal(t, []string{"first", "first"}, onErrorCalls)
}

func TestRun_ReturnsErrorWhenIterationBudgetExhausted(t *testing.T) {
	actors := []int{1}

	results, err := Run(Config[int, int]{
		MaxIterations: 3,
		ForEachActor: func(visit func(actor int) bool) {
			for _, actor := range actors {
				if !visit(actor) {
					return
				}
			}
		},
		ProcessActor: func(actor int, emit func(result int)) (bool, bool, error) {
			emit(actor)
			return true, false, nil
		},
	})

	require.ErrorIs(t, err, ErrMaxIterationsReached)
	assert.Equal(t, []int{1, 1, 1}, results)
}

func TestRun_ValidatesConfig(t *testing.T) {
	_, err := Run(Config[int, int]{})
	require.ErrorIs(t, err, ErrInvalidMaxIterations)

	_, err = Run(Config[int, int]{MaxIterations: 1})
	require.ErrorIs(t, err, ErrNilForEachActor)

	_, err = Run(Config[int, int]{
		MaxIterations: 1,
		ForEachActor:  func(func(int) bool) {},
	})
	require.ErrorIs(t, err, ErrNilProcessActor)
}
