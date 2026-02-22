package turnbased

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_ValidatesPlayers(t *testing.T) {
	t.Parallel()

	_, err := New[int](nil, 1)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrEmptyPlayers)

	_, err = New[int]([]int{1, 1}, 1)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrDuplicatePlayer)

	_, err = New[int]([]int{1, 2}, 3)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrUnknownPlayer)
}

func TestStep_EnforcesTurnAndRotation(t *testing.T) {
	t.Parallel()

	e, err := New[int]([]int{1, 2, 3}, 1)
	require.NoError(t, err)

	_, err = e.Step(2, NextTurn[int]())
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrWrongTurn)

	delta, err := e.Step(1, NextTurn[int]())
	require.NoError(t, err)
	assert.Equal(t, 1, delta.Actor)
	assert.Equal(t, 1, delta.PreviousPlayer)
	assert.Equal(t, 2, delta.CurrentPlayer)
	assert.True(t, delta.TurnChanged)
	assert.False(t, delta.IsTerminal)
	assert.Equal(t, MatchResultOngoing, delta.Result)
	assert.Equal(t, 2, e.CurrentPlayer())

	_, err = e.Step(2, NextTurn[int]())
	require.NoError(t, err)
	assert.Equal(t, 3, e.CurrentPlayer())

	_, err = e.Step(3, NextTurn[int]())
	require.NoError(t, err)
	assert.Equal(t, 1, e.CurrentPlayer())
}

func TestStep_KeepTurn(t *testing.T) {
	t.Parallel()

	e, err := New[string]([]string{"p1", "p2"}, "p1")
	require.NoError(t, err)

	delta, err := e.Step("p1", KeepTurn[string]())
	require.NoError(t, err)
	assert.False(t, delta.TurnChanged)
	assert.Equal(t, "p1", delta.CurrentPlayer)
	assert.Equal(t, "p1", e.CurrentPlayer())
}

func TestStep_WinnerFinalizesGame(t *testing.T) {
	t.Parallel()

	e, err := New[string]([]string{"alpha", "beta"}, "beta")
	require.NoError(t, err)

	delta, err := e.Step("beta", Win("beta"))
	require.NoError(t, err)
	assert.True(t, delta.IsTerminal)
	assert.Equal(t, MatchResultWinner, delta.Result)
	require.True(t, delta.HasWinner)
	assert.Equal(t, "beta", delta.Winner)
	assert.False(t, delta.TurnChanged)

	assert.True(t, e.IsOver())
	assert.Equal(t, MatchResultWinner, e.Result())
	winner, ok := e.Winner()
	require.True(t, ok)
	assert.Equal(t, "beta", winner)

	_, err = e.Step("beta", NextTurn[string]())
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrGameOver)
}

func TestStep_DrawFinalizesGame(t *testing.T) {
	t.Parallel()

	e, err := New[string]([]string{"alpha", "beta"}, "alpha")
	require.NoError(t, err)

	delta, err := e.Step("alpha", Draw[string]())
	require.NoError(t, err)
	assert.True(t, delta.IsTerminal)
	assert.Equal(t, MatchResultDraw, delta.Result)
	assert.False(t, delta.HasWinner)
	assert.True(t, e.IsDraw())

	_, hasWinner := e.Winner()
	assert.False(t, hasWinner)

	_, err = e.Step("beta", NextTurn[string]())
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrGameOver)
}

func TestStep_RejectsInvalidDecision(t *testing.T) {
	t.Parallel()

	e, err := New[string]([]string{"first", "second"}, "first")
	require.NoError(t, err)

	_, err = e.Step("first", Decision[string]{
		result: MatchResult(99),
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidDecision)
	assert.Equal(t, "first", e.CurrentPlayer())
	assert.False(t, e.IsOver())

	_, err = e.Step("first", Decision[string]{
		result: MatchResultWinner,
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidDecision)
	assert.Equal(t, "first", e.CurrentPlayer())
	assert.False(t, e.IsOver())
}

func TestStep_RejectsUnknownWinner(t *testing.T) {
	t.Parallel()

	e, err := New[string]([]string{"first", "second"}, "first")
	require.NoError(t, err)

	_, err = e.Step("first", Win("unknown"))
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrUnknownPlayer)
	assert.Equal(t, MatchResultOngoing, e.Result())
}

func TestStep_UnknownActor(t *testing.T) {
	t.Parallel()

	e, err := New[int]([]int{1, 2}, 1)
	require.NoError(t, err)

	_, err = e.Step(99, NextTurn[int]())
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrUnknownPlayer)
}

func TestDecisionHelpers(t *testing.T) {
	t.Parallel()

	next := NextTurn[string]()
	assert.False(t, next.KeepsTurn())
	assert.Equal(t, MatchResultOngoing, next.Result())
	_, ok := next.Winner()
	assert.False(t, ok)
	assert.False(t, next.Draw())

	keep := KeepTurn[string]()
	assert.True(t, keep.KeepsTurn())
	assert.Equal(t, MatchResultOngoing, keep.Result())

	win := Win("p1")
	assert.False(t, win.KeepsTurn())
	assert.Equal(t, MatchResultWinner, win.Result())
	winner, ok := win.Winner()
	require.True(t, ok)
	assert.Equal(t, "p1", winner)
	assert.False(t, win.Draw())

	draw := Draw[string]()
	assert.False(t, draw.KeepsTurn())
	assert.True(t, draw.Draw())
	assert.Equal(t, MatchResultDraw, draw.Result())
	_, ok = draw.Winner()
	assert.False(t, ok)
}

func TestPlayerAccessors_NoCopyAPI(t *testing.T) {
	t.Parallel()

	e, err := New[string]([]string{"p1", "p2", "p3"}, "p1")
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

	e, err := New[int]([]int{10, 20, 30}, 10)
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
