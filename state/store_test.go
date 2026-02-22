package state_test

import (
	"errors"
	"testing"

	"github.com/RevoTale/turn-based-game-engine/state"
	"github.com/stretchr/testify/require"
)

func TestDoCommitIncrementsVersion(t *testing.T) {
	t.Parallel()

	type payload struct {
		Count int
	}
	store := state.New(payload{})

	version, err := store.Do(func(current *payload, _ uint64) error {
		current.Count++
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, uint64(1), version)

	var count int
	err = store.View(func(s *payload, _ uint64) {
		count = s.Count
	})
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestDoFailureDoesNotIncrementVersion(t *testing.T) {
	t.Parallel()

	type payload struct {
		Count int
	}
	store := state.New(payload{})
	expectedErr := errors.New("boom")

	_, err := store.Do(func(current *payload, _ uint64) error {
		current.Count++
		return expectedErr
	})
	require.ErrorIs(t, err, expectedErr)

	var (
		count   int
		version uint64
	)
	err = store.View(func(s *payload, v uint64) {
		count = s.Count
		version = v
	})
	require.NoError(t, err)
	require.Equal(t, 1, count)
	require.Equal(t, uint64(0), version)
}

func TestDoRejectsNilCallback(t *testing.T) {
	t.Parallel()

	store := state.New(struct{}{})
	_, err := store.Do(nil)
	require.ErrorIs(t, err, state.ErrNilTransaction)
}
