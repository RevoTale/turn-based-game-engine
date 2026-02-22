package tictactoe

import (
	"errors"

	"github.com/RevoTale/turn-based-game-engine/grid2d"
	"github.com/RevoTale/turn-based-game-engine/runtime/events"
	"github.com/RevoTale/turn-based-game-engine/turnbased"
)

var errMoveNotPrepared = errors.New("move is not prepared")

type boardWrite struct {
	pos    grid2d.Position
	player Player
}

type gamePatch struct {
	movePrepared bool
	move         Move

	lastActor    Player
	lastDecision turnbased.Decision[Player]
	lastStep     turnbased.Delta[Player]
	lastApplied  bool

	boardWrites []boardWrite
	log         []string
}

type gameEvents struct {
	play          events.Command[*gameState, gamePatch, Move]
	resolveMove   events.Event[*gameState, gamePatch]
	moveApplied   events.Event[*gameState, gamePatch]
	turnChanged   events.Event[*gameState, gamePatch]
	matchFinished events.Event[*gameState, gamePatch]
}

type playContext interface {
	State() *gameState
	Patch() *gamePatch
	EmitResolveMove()
}

type resolveContext interface {
	State() *gameState
	Patch() *gamePatch
	EmitMoveApplied()
	EmitTurnChanged()
	EmitMatchFinished()
}

type playCtx struct {
	events.Context[*gameState, gamePatch]
	ev gameEvents
}

func (c playCtx) EmitResolveMove() {
	c.Emit(c.ev.resolveMove)
}

type resolveCtx struct {
	events.Context[*gameState, gamePatch]
	ev gameEvents
}

func (c resolveCtx) EmitMoveApplied() {
	c.Emit(c.ev.moveApplied)
}

func (c resolveCtx) EmitTurnChanged() {
	c.Emit(c.ev.turnChanged)
}

func (c resolveCtx) EmitMatchFinished() {
	c.Emit(c.ev.matchFinished)
}

func (g *Game) registerEvents() (gameEvents, error) {
	registered := gameEvents{}
	var err error

	registered.moveApplied, err = events.DefineEvent(handleMoveApplied)
	if err != nil {
		return gameEvents{}, err
	}
	registered.turnChanged, err = events.DefineEvent(handleTurnChanged)
	if err != nil {
		return gameEvents{}, err
	}
	registered.matchFinished, err = events.DefineEvent(handleMatchFinished)
	if err != nil {
		return gameEvents{}, err
	}
	registered.resolveMove, err = events.DefineEvent(bindResolveHandler(registered))
	if err != nil {
		return gameEvents{}, err
	}
	registered.play, err = events.DefineCommand(bindPlayHandler(registered))
	if err != nil {
		return gameEvents{}, err
	}
	return registered, nil
}

func bindPlayHandler(ev gameEvents) events.CommandHandler[*gameState, gamePatch, Move] {
	return func(ctx events.Context[*gameState, gamePatch], input Move) error {
		return handlePlay(playCtx{
			Context: ctx,
			ev:      ev,
		}, input)
	}
}

func bindResolveHandler(ev gameEvents) events.EventHandler[*gameState, gamePatch] {
	return func(ctx events.Context[*gameState, gamePatch]) error {
		return handleResolveMove(resolveCtx{
			Context: ctx,
			ev:      ev,
		})
	}
}

func handlePlay(ctx playContext, input Move) error {
	state := ctx.State()
	if state.turns.IsOver() {
		return ErrGameFinished
	}

	patch := ctx.Patch()
	patch.movePrepared = true
	patch.move = input
	ctx.EmitResolveMove()
	return nil
}

func handleResolveMove(ctx resolveContext) error {
	patch := ctx.Patch()
	if !patch.movePrepared {
		return errMoveNotPrepared
	}

	actor, stepDecision, stepDelta, writes, err := applyMove(ctx.State(), patch, patch.move)
	if err != nil {
		return err
	}

	patch.lastActor = actor
	patch.lastDecision = stepDecision
	patch.lastStep = stepDelta
	patch.lastApplied = true
	patch.boardWrites = append(patch.boardWrites, writes...)

	ctx.EmitMoveApplied()
	if stepDelta.TurnChanged {
		ctx.EmitTurnChanged()
	}
	if stepDelta.IsTerminal {
		ctx.EmitMatchFinished()
	}
	return nil
}

func handleMoveApplied(ctx events.Context[*gameState, gamePatch]) error {
	patch := ctx.Patch()
	if !patch.lastApplied {
		return nil
	}
	patch.log = append(patch.log, formatMoveLog(patch.lastStep.Actor, patch.move.Index))
	return nil
}

func handleTurnChanged(ctx events.Context[*gameState, gamePatch]) error {
	patch := ctx.Patch()
	if !patch.lastApplied || !patch.lastStep.TurnChanged {
		return nil
	}
	patch.log = append(patch.log, formatTurnChangedLog(patch.lastStep.PreviousPlayer, patch.lastStep.CurrentPlayer))
	return nil
}

func handleMatchFinished(ctx events.Context[*gameState, gamePatch]) error {
	patch := ctx.Patch()
	if !patch.lastApplied {
		return nil
	}
	patch.log = append(patch.log, formatMatchLog(patch.lastStep))
	return nil
}
