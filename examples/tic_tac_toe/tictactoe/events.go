package tictactoe

import (
	"fmt"

	"github.com/RevoTale/turn-based-game-engine/grid2d"
	"github.com/RevoTale/turn-based-game-engine/runtime/events"
	"github.com/RevoTale/turn-based-game-engine/turnbased"
)

type moveApplied struct {
	Player Player
	Index  int
}

type matchFinished struct {
	Result turnbased.MatchResult
	Winner Player
}

type gameEvents struct {
	play          events.Command[Move]
	moveApplied   events.Event[moveApplied]
	matchFinished events.Event[matchFinished]
}

func (g *Game) registerEvents(builder *events.Builder) (gameEvents, error) {
	registered := gameEvents{}
	var err error

	registered.moveApplied, err = events.RegisterEvent(builder, g.handleMoveApplied, events.Hooks[moveApplied]{})
	if err != nil {
		return gameEvents{}, err
	}
	registered.matchFinished, err = events.RegisterEvent(builder, g.handleMatchFinished, events.Hooks[matchFinished]{})
	if err != nil {
		return gameEvents{}, err
	}
	registered.play, err = events.RegisterCommand(builder, g.handlePlay(registered), events.Hooks[Move]{})
	if err != nil {
		return gameEvents{}, err
	}

	return registered, nil
}

func (g *Game) handleMoveApplied(_ *events.Context, payload moveApplied) error {
	g.log = append(g.log, fmt.Sprintf("move: %c -> %d", payload.Player, payload.Index))
	return nil
}

func (g *Game) handleMatchFinished(_ *events.Context, payload matchFinished) error {
	if payload.Result == turnbased.MatchResultDraw {
		g.log = append(g.log, "match finished: draw")
		return nil
	}
	g.log = append(g.log, fmt.Sprintf("match finished: winner=%c", payload.Winner))
	return nil
}

func (g *Game) handlePlay(ev gameEvents) events.Handler[Move] {
	return func(ctx *events.Context, move Move) error {
		if g.turns.IsOver() {
			return ErrGameFinished
		}

		actor := g.turns.CurrentPlayer()
		outcome, err := g.turns.Act(actor, move, func(actor Player, action Move) (turnbased.ActionOutcome[Player], error) {
			pos, ok := g.grid.Position(grid2d.CellIndex(action.Index))
			if !ok {
				return turnbased.ActionOutcome[Player]{}, ErrMoveBounds
			}
			if _, occupied, getErr := g.board.Get(pos); getErr != nil {
				return turnbased.ActionOutcome[Player]{}, getErr
			} else if occupied {
				return turnbased.ActionOutcome[Player]{}, ErrCellBusy
			}

			if setErr := g.board.Set(pos, actor); setErr != nil {
				return turnbased.ActionOutcome[Player]{}, setErr
			}
			if g.hasWinner(actor) {
				return turnbased.PassTurn[Player]().WithWinner(actor), nil
			}
			if g.isBoardFull() {
				return turnbased.PassTurn[Player]().WithDraw(), nil
			}
			return turnbased.PassTurn[Player](), nil
		})
		if err != nil {
			return err
		}

		ctx.Emit(events.Next(ev.moveApplied, moveApplied{Player: actor, Index: move.Index}))

		if winner, ok := outcome.Winner(); ok {
			ctx.Emit(events.Next(ev.matchFinished, matchFinished{
				Result: turnbased.MatchResultWinner,
				Winner: winner,
			}))
			return nil
		}
		if outcome.Draw() {
			ctx.Emit(events.Next(ev.matchFinished, matchFinished{
				Result: turnbased.MatchResultDraw,
			}))
		}
		return nil
	}
}
