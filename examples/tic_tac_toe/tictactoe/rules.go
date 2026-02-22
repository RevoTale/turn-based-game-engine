package tictactoe

import (
	"fmt"

	"github.com/RevoTale/turn-based-game-engine/grid2d"
	"github.com/RevoTale/turn-based-game-engine/turnbased"
)

func applyMove(st *gameState, patch *gamePatch, move Move) (Player, turnbased.Decision[Player], turnbased.Delta[Player], []boardWrite, error) {
	actor := st.turns.CurrentPlayer()
	writes := make([]boardWrite, 0, 1)

	pos, ok := st.grid.Position(grid2d.CellIndex(move.Index))
	if !ok {
		return 0, turnbased.Decision[Player]{}, turnbased.Delta[Player]{}, nil, ErrMoveBounds
	}
	if isOccupied(st, patch, writes, pos) {
		return 0, turnbased.Decision[Player]{}, turnbased.Delta[Player]{}, nil, ErrCellBusy
	}

	writes = append(writes, boardWrite{
		pos:    pos,
		player: actor,
	})

	decision := turnbased.NextTurn[Player]()
	if hasWinner(st, patch, writes, actor) {
		decision = turnbased.Win(actor)
	} else if isBoardFull(st, patch, writes) {
		decision = turnbased.Draw[Player]()
	}

	turns := *st.turns
	delta, err := turns.Step(actor, decision)
	if err != nil {
		return 0, turnbased.Decision[Player]{}, turnbased.Delta[Player]{}, nil, err
	}
	return actor, decision, delta, writes, nil
}

func applyPatch(st *gameState, patch *gamePatch) error {
	if patch == nil {
		return nil
	}

	for _, write := range patch.boardWrites {
		if err := st.board.Set(write.pos, write.player); err != nil {
			return err
		}
	}
	if patch.lastApplied {
		if _, err := st.turns.Step(patch.lastActor, patch.lastDecision); err != nil {
			return err
		}
	}
	if len(patch.log) > 0 {
		st.log = append(st.log, patch.log...)
	}
	return nil
}

func formatMoveLog(player Player, index int) string {
	return fmt.Sprintf("move: %c -> %d", player, index)
}

func formatTurnChangedLog(from Player, to Player) string {
	return fmt.Sprintf("turn changed: %c -> %c", from, to)
}

func formatMatchLog(delta turnbased.Delta[Player]) string {
	if delta.Result == turnbased.MatchResultDraw {
		return "match finished: draw"
	}
	if delta.HasWinner {
		return fmt.Sprintf("match finished: winner=%c", delta.Winner)
	}
	return "match finished"
}

func hasWinner(st *gameState, patch *gamePatch, writes []boardWrite, player Player) bool {
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
		if cellOwner(st, patch, writes, line[0]) == player &&
			cellOwner(st, patch, writes, line[1]) == player &&
			cellOwner(st, patch, writes, line[2]) == player {
			return true
		}
	}
	return false
}

func isBoardFull(st *gameState, patch *gamePatch, writes []boardWrite) bool {
	for y := 0; y < st.grid.Height(); y++ {
		for x := 0; x < st.grid.Width(); x++ {
			if cellOwner(st, patch, writes, grid2d.Position{X: x, Y: y}) == 0 {
				return false
			}
		}
	}
	return true
}

func isOccupied(st *gameState, patch *gamePatch, writes []boardWrite, pos grid2d.Position) bool {
	return cellOwner(st, patch, writes, pos) != 0
}

func cellOwner(st *gameState, patch *gamePatch, writes []boardWrite, pos grid2d.Position) Player {
	for i := len(writes) - 1; i >= 0; i-- {
		if writes[i].pos == pos {
			return writes[i].player
		}
	}
	if patch != nil {
		for i := len(patch.boardWrites) - 1; i >= 0; i-- {
			if patch.boardWrites[i].pos == pos {
				return patch.boardWrites[i].player
			}
		}
	}
	value, ok, err := st.board.Get(pos)
	if err != nil || !ok {
		return 0
	}
	return value
}
