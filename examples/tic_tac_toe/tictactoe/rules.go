package tictactoe

import (
	"fmt"

	"github.com/RevoTale/turn-based-game-engine/grid2d"
	"github.com/RevoTale/turn-based-game-engine/turnbased"
)

func applyMove(st *gameState, move Move) (Player, turnbased.ActionOutcome[Player], error) {
	actor := st.turns.CurrentPlayer()
	outcome, err := st.turns.Act(actor, move, func(actor Player, action Move) (turnbased.ActionOutcome[Player], error) {
		pos, ok := st.grid.Position(grid2d.CellIndex(action.Index))
		if !ok {
			return turnbased.ActionOutcome[Player]{}, ErrMoveBounds
		}
		if _, occupied, getErr := st.board.Get(pos); getErr != nil {
			return turnbased.ActionOutcome[Player]{}, getErr
		} else if occupied {
			return turnbased.ActionOutcome[Player]{}, ErrCellBusy
		}

		if setErr := st.board.Set(pos, actor); setErr != nil {
			return turnbased.ActionOutcome[Player]{}, setErr
		}
		if hasWinner(st, actor) {
			return turnbased.PassTurn[Player]().WithWinner(actor), nil
		}
		if isBoardFull(st) {
			return turnbased.PassTurn[Player]().WithDraw(), nil
		}
		return turnbased.PassTurn[Player](), nil
	})
	return actor, outcome, err
}

func appendMoveLog(st *gameState, player Player, index int) {
	st.log = append(st.log, fmt.Sprintf("move: %c -> %d", player, index))
}

func appendMatchLog(st *gameState, result turnbased.MatchResult, winner Player) {
	if result == turnbased.MatchResultDraw {
		st.log = append(st.log, "match finished: draw")
		return
	}
	st.log = append(st.log, fmt.Sprintf("match finished: winner=%c", winner))
}

func hasWinner(st *gameState, player Player) bool {
	lines := [8][3]grid2d.Position{
		{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
		{{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}},
		{{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}},
		{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}},
		{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
		{{X: 2, Y: 0}, {X: 2, Y: 1}, {X: 2, Y: 2}},
		{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2}},
		{{X: 2, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 2}},
	}
	for _, line := range lines {
		if cellOwner(st, line[0]) == player && cellOwner(st, line[1]) == player && cellOwner(st, line[2]) == player {
			return true
		}
	}
	return false
}

func isBoardFull(st *gameState) bool {
	return st.board.Len() == st.grid.CellCount()
}

func cellOwner(st *gameState, pos grid2d.Position) Player {
	value, ok, err := st.board.Get(pos)
	if err != nil || !ok {
		return 0
	}
	return value
}
