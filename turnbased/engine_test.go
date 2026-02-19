package turnbased

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_ValidatesPlayers(t *testing.T) {
	t.Parallel()

	_, err := New[int, int](nil, 1)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrEmptyPlayers)

	_, err = New[int, int]([]int{1, 1}, 1)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrDuplicatePlayer)

	_, err = New[int, int]([]int{1, 2}, 3)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrUnknownPlayer)
}

func TestAct_EnforcesTurnAndRotation(t *testing.T) {
	t.Parallel()

	e, err := New[int, int]([]int{1, 2, 3}, 1)
	require.NoError(t, err)

	_, err = e.Act(2, 0, func(actor int, action int) (ActionOutcome[int], error) {
		return PassTurn[int](), nil
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrWrongTurn)

	_, err = e.Act(1, 0, func(actor int, action int) (ActionOutcome[int], error) {
		return PassTurn[int](), nil
	})
	require.NoError(t, err)
	assert.Equal(t, 2, e.CurrentPlayer())

	_, err = e.Act(2, 0, func(actor int, action int) (ActionOutcome[int], error) {
		return PassTurn[int](), nil
	})
	require.NoError(t, err)
	assert.Equal(t, 3, e.CurrentPlayer())

	_, err = e.Act(3, 0, func(actor int, action int) (ActionOutcome[int], error) {
		return PassTurn[int](), nil
	})
	require.NoError(t, err)
	assert.Equal(t, 1, e.CurrentPlayer())
}

func TestAct_KeepTurn(t *testing.T) {
	t.Parallel()

	e, err := New[string, int]([]string{"p1", "p2"}, "p1")
	require.NoError(t, err)

	_, err = e.Act("p1", 42, func(actor string, action int) (ActionOutcome[string], error) {
		return KeepTurn[string](), nil
	})
	require.NoError(t, err)
	assert.Equal(t, "p1", e.CurrentPlayer())
}

func TestAct_WinnerFinalizesGame(t *testing.T) {
	t.Parallel()

	e, err := New[string, int]([]string{"alpha", "beta"}, "beta")
	require.NoError(t, err)

	winner := "beta"
	_, err = e.Act("beta", 0, func(actor string, action int) (ActionOutcome[string], error) {
		return PassTurn[string]().WithWinner(winner), nil
	})
	require.NoError(t, err)
	assert.True(t, e.IsOver())
	assert.Equal(t, MatchResultWinner, e.Result())

	gotWinner, ok := e.Winner()
	require.True(t, ok)
	assert.Equal(t, winner, gotWinner)

	_, err = e.Act("beta", 0, func(actor string, action int) (ActionOutcome[string], error) {
		return PassTurn[string](), nil
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrGameOver)
}

func TestAct_DrawFinalizesGame(t *testing.T) {
	t.Parallel()

	e, err := New[string, int]([]string{"alpha", "beta"}, "alpha")
	require.NoError(t, err)

	_, err = e.Act("alpha", 0, func(actor string, action int) (ActionOutcome[string], error) {
		return PassTurn[string]().WithDraw(), nil
	})
	require.NoError(t, err)
	assert.True(t, e.IsOver())
	assert.Equal(t, MatchResultDraw, e.Result())
	assert.True(t, e.IsDraw())

	_, hasWinner := e.Winner()
	assert.False(t, hasWinner)

	_, err = e.Act("beta", 0, func(actor string, action int) (ActionOutcome[string], error) {
		return PassTurn[string](), nil
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrGameOver)
}

func TestAct_PropagatesResolverErrorWithoutChangingTurn(t *testing.T) {
	t.Parallel()

	e, err := New[string, int]([]string{"first", "second"}, "first")
	require.NoError(t, err)

	resolveErr := errors.New("boom")
	_, err = e.Act("first", 0, func(actor string, action int) (ActionOutcome[string], error) {
		return PassTurn[string](), resolveErr
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, resolveErr)
	assert.Equal(t, "first", e.CurrentPlayer())
}

func TestAct_RejectsInvalidTerminalOutcome(t *testing.T) {
	t.Parallel()

	e, err := New[string, int]([]string{"first", "second"}, "first")
	require.NoError(t, err)

	_, err = e.Act("first", 0, func(actor string, action int) (ActionOutcome[string], error) {
		return ActionOutcome[string]{
			result: MatchResult(99),
		}, nil
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidOutcome)
	assert.Equal(t, "first", e.CurrentPlayer())
	assert.False(t, e.IsOver())
}

func TestActionOutcomeEncapsulation(t *testing.T) {
	t.Parallel()

	outcome := PassTurn[string]()
	assert.False(t, outcome.KeepsTurn())
	_, hasWinner := outcome.Winner()
	assert.False(t, hasWinner)
	assert.False(t, outcome.Draw())
	assert.Equal(t, MatchResultOngoing, outcome.Result())

	outcome = KeepTurn[string]().WithWinner("p1")
	assert.True(t, outcome.KeepsTurn())
	winner, hasWinner := outcome.Winner()
	require.True(t, hasWinner)
	assert.Equal(t, "p1", winner)
	assert.False(t, outcome.Draw())
	assert.Equal(t, MatchResultWinner, outcome.Result())

	outcome = PassTurn[string]().WithDraw()
	assert.False(t, outcome.KeepsTurn())
	assert.True(t, outcome.Draw())
	_, hasWinner = outcome.Winner()
	assert.False(t, hasWinner)
	assert.Equal(t, MatchResultDraw, outcome.Result())
}

func TestPlayerAccessors_NoCopyAPI(t *testing.T) {
	t.Parallel()

	e, err := New[string, int]([]string{"p1", "p2", "p3"}, "p1")
	require.NoError(t, err)

	assert.Equal(t, 3, e.PlayerCount())
	assert.True(t, e.HasPlayer("p2"))
	assert.False(t, e.HasPlayer("p4"))

	p, ok := e.PlayerAt(1)
	require.True(t, ok)
	assert.Equal(t, "p2", p)

	_, ok = e.PlayerAt(-1)
	assert.False(t, ok)
	_, ok = e.PlayerAt(3)
	assert.False(t, ok)
}

func TestForEachPlayer_OrderAndEarlyStop(t *testing.T) {
	t.Parallel()

	e, err := New[int, int]([]int{10, 20, 30}, 10)
	require.NoError(t, err)

	visited := make([]int, 0, 3)
	e.ForEachPlayer(func(player int) bool {
		visited = append(visited, player)
		return true
	})
	assert.Equal(t, []int{10, 20, 30}, visited)

	visited = visited[:0]
	e.ForEachPlayer(func(player int) bool {
		visited = append(visited, player)
		return len(visited) < 2
	})
	assert.Equal(t, []int{10, 20}, visited)
}
