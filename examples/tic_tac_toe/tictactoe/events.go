package tictactoe

import (
	"errors"

	"github.com/RevoTale/turn-based-game-engine/runtime/events"
	"github.com/RevoTale/turn-based-game-engine/turnbased"
)

var errNilStatePatch = errors.New("state patch is nil")

type playCommand struct {
	State *gameState
	Move  Move
}

type moveApplied struct {
	State  *gameState
	Player Player
	Index  int
}

type matchFinished struct {
	State  *gameState
	Result turnbased.MatchResult
	Winner Player
}

type gameEvents struct {
	play          events.Command[playCommand]
	moveApplied   events.Event[moveApplied]
	matchFinished events.Event[matchFinished]
}

func (g *Game) registerEvents(builder *events.Builder) (gameEvents, error) {
	registered := gameEvents{}
	var err error

	registered.moveApplied, err = events.RegisterEvent(builder, handleMoveApplied, events.Hooks[moveApplied]{})
	if err != nil {
		return gameEvents{}, err
	}
	registered.matchFinished, err = events.RegisterEvent(builder, handleMatchFinished, events.Hooks[matchFinished]{})
	if err != nil {
		return gameEvents{}, err
	}
	registered.play, err = events.RegisterCommand(builder, handlePlay(registered), events.Hooks[playCommand]{})
	if err != nil {
		return gameEvents{}, err
	}

	return registered, nil
}

func handleMoveApplied(_ *events.Context, payload moveApplied) error {
	if payload.State == nil {
		return errNilStatePatch
	}
	appendMoveLog(payload.State, payload.Player, payload.Index)
	return nil
}

func handleMatchFinished(_ *events.Context, payload matchFinished) error {
	if payload.State == nil {
		return errNilStatePatch
	}
	appendMatchLog(payload.State, payload.Result, payload.Winner)
	return nil
}

func handlePlay(ev gameEvents) events.Handler[playCommand] {
	return func(ctx *events.Context, payload playCommand) error {
		if payload.State == nil {
			return errNilStatePatch
		}
		st := payload.State

		if st.turns.IsOver() {
			return ErrGameFinished
		}

		actor, outcome, err := applyMove(st, payload.Move)
		if err != nil {
			return err
		}

		ctx.Emit(events.Next(ev.moveApplied, moveApplied{
			State:  st,
			Player: actor,
			Index:  payload.Move.Index,
		}))

		if winner, ok := outcome.Winner(); ok {
			ctx.Emit(events.Next(ev.matchFinished, matchFinished{
				State:  st,
				Result: turnbased.MatchResultWinner,
				Winner: winner,
			}))
			return nil
		}
		if outcome.Draw() {
			ctx.Emit(events.Next(ev.matchFinished, matchFinished{
				State:  st,
				Result: turnbased.MatchResultDraw,
			}))
		}
		return nil
	}
}
