package turnbased

import (
	"fmt"
)

// Engine tracks turn order, current turn, and terminal match state.
type Engine[P comparable] struct {
	order         []P
	indexByPlayer map[P]int
	turnIndex     int
	result        MatchResult
	winner        P
}

// New creates an Engine with player order and first turn.
//
// Player order is kept as provided and used for turn rotation.
// Returns:
// - ErrEmptyPlayers when players is empty
// - ErrDuplicatePlayer when players contains duplicates
// - ErrUnknownPlayer when firstTurn is not in players
func New[P comparable](players []P, firstTurn P) (*Engine[P], error) {
	if len(players) == 0 {
		return nil, ErrEmptyPlayers
	}

	indexByPlayer := make(map[P]int, len(players))
	order := make([]P, len(players))
	copy(order, players)

	for i, player := range order {
		if _, exists := indexByPlayer[player]; exists {
			return nil, fmt.Errorf("%w: %v", ErrDuplicatePlayer, player)
		}
		indexByPlayer[player] = i
	}

	startIndex, ok := indexByPlayer[firstTurn]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrUnknownPlayer, firstTurn)
	}

	return &Engine[P]{
		order:         order,
		indexByPlayer: indexByPlayer,
		turnIndex:     startIndex,
	}, nil
}

// CurrentPlayer returns the player who can act now.
func (e *Engine[P]) CurrentPlayer() P {
	return e.order[e.turnIndex]
}

// PlayerCount returns number of players in turn order.
func (e *Engine[P]) PlayerCount() int {
	return len(e.order)
}

// HasPlayer reports whether player is part of this game.
func (e *Engine[P]) HasPlayer(player P) bool {
	_, ok := e.indexByPlayer[player]
	return ok
}

// PlayerAt returns player at index in turn order.
//
// Second return value is false when index is out of range.
func (e *Engine[P]) PlayerAt(index int) (P, bool) {
	if index < 0 || index >= len(e.order) {
		var zero P
		return zero, false
	}
	return e.order[index], true
}

// ForEachPlayer calls visit for players in turn order.
//
// Returning false from visit stops iteration early.
func (e *Engine[P]) ForEachPlayer(visit func(player P) bool) {
	if visit == nil {
		return
	}
	for _, player := range e.order {
		if !visit(player) {
			return
		}
	}
}

// IsOver reports whether the game already reached any terminal state.
func (e *Engine[P]) IsOver() bool {
	return e.result != MatchResultOngoing
}

// Winner returns winning player when game is over.
//
// Second return value is false before game over.
func (e *Engine[P]) Winner() (P, bool) {
	if e.result != MatchResultWinner {
		var zero P
		return zero, false
	}
	return e.winner, true
}

// IsDraw reports whether the game ended in draw.
func (e *Engine[P]) IsDraw() bool {
	return e.result == MatchResultDraw
}

// Result returns current terminal result of the match.
func (e *Engine[P]) Result() MatchResult {
	return e.result
}

// Step applies one domain decision for the actor currently taking a turn.
//
// The method validates actor ownership, validates decision semantics, mutates
// turn state, and returns one Delta describing what changed.
func (e *Engine[P]) Step(actor P, decision Decision[P]) (Delta[P], error) {
	before := e.CurrentPlayer()
	if e.IsOver() {
		return Delta[P]{}, ErrGameOver
	}

	actorIndex, ok := e.indexByPlayer[actor]
	if !ok {
		return Delta[P]{}, fmt.Errorf("%w: %v", ErrUnknownPlayer, actor)
	}
	if actorIndex != e.turnIndex {
		return Delta[P]{}, ErrWrongTurn
	}

	switch decision.Result() {
	case MatchResultOngoing, MatchResultWinner, MatchResultDraw:
	default:
		return Delta[P]{}, ErrInvalidDecision
	}

	if winner, hasWinner := decision.Winner(); hasWinner {
		if decision.Result() != MatchResultWinner {
			return Delta[P]{}, ErrInvalidDecision
		}
		if _, exists := e.indexByPlayer[winner]; !exists {
			return Delta[P]{}, fmt.Errorf("%w: %v", ErrUnknownPlayer, winner)
		}
		e.winner = winner
		e.result = MatchResultWinner
		return e.buildDelta(actor, before), nil
	}
	if decision.Result() == MatchResultWinner {
		return Delta[P]{}, ErrInvalidDecision
	}
	if decision.Draw() {
		e.result = MatchResultDraw
		return e.buildDelta(actor, before), nil
	}

	if !decision.KeepsTurn() {
		e.turnIndex = (e.turnIndex + 1) % len(e.order)
	}

	return e.buildDelta(actor, before), nil
}

func (e *Engine[P]) buildDelta(actor P, previousPlayer P) Delta[P] {
	current := e.CurrentPlayer()
	delta := Delta[P]{
		Actor:          actor,
		PreviousPlayer: previousPlayer,
		CurrentPlayer:  current,
		TurnChanged:    current != previousPlayer,
		Result:         e.result,
		IsTerminal:     e.result != MatchResultOngoing,
	}

	if winner, ok := e.Winner(); ok {
		delta.Winner = winner
		delta.HasWinner = true
	}

	return delta
}
