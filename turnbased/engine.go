package turnbased

import (
	"fmt"
)

// Engine tracks turn order, current turn, and winner state.
type Engine[P comparable, A any] struct {
	order         []P
	indexByPlayer map[P]int
	turnIndex     int
	winner        P
	hasWinner     bool
}

// New creates an Engine with player order and first turn.
//
// Player order is kept as provided and used for turn rotation.
// Returns:
// - ErrEmptyPlayers when players is empty
// - ErrDuplicatePlayer when players contains duplicates
// - ErrUnknownPlayer when firstTurn is not in players
func New[P comparable, A any](players []P, firstTurn P) (*Engine[P, A], error) {
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

	return &Engine[P, A]{
		order:         order,
		indexByPlayer: indexByPlayer,
		turnIndex:     startIndex,
	}, nil
}

// CurrentPlayer returns the player who can act now.
func (e *Engine[P, A]) CurrentPlayer() P {
	return e.order[e.turnIndex]
}

// PlayerCount returns number of players in turn order.
func (e *Engine[P, A]) PlayerCount() int {
	return len(e.order)
}

// HasPlayer reports whether player is part of this game.
func (e *Engine[P, A]) HasPlayer(player P) bool {
	_, ok := e.indexByPlayer[player]
	return ok
}

// PlayerAt returns player at index in turn order.
//
// Second return value is false when index is out of range.
func (e *Engine[P, A]) PlayerAt(index int) (P, bool) {
	if index < 0 || index >= len(e.order) {
		var zero P
		return zero, false
	}
	return e.order[index], true
}

// ForEachPlayer calls visit for players in turn order.
//
// Returning false from visit stops iteration early.
func (e *Engine[P, A]) ForEachPlayer(visit func(player P) bool) {
	if visit == nil {
		return
	}
	for _, player := range e.order {
		if !visit(player) {
			return
		}
	}
}

// IsOver reports whether the game already has a winner.
func (e *Engine[P, A]) IsOver() bool {
	return e.hasWinner
}

// Winner returns winning player when game is over.
//
// Second return value is false before game over.
func (e *Engine[P, A]) Winner() (P, bool) {
	if !e.hasWinner {
		var zero P
		return zero, false
	}
	return e.winner, true
}

// Act applies one player action through resolve callback.
//
// Act checks:
// - known actor and correct turn
// - no actions after game over
// - winner (if provided by resolver) must be a known player
//
// Turn moves to next player unless outcome says to keep turn.
// If resolve returns an error, turn state is unchanged.
func (e *Engine[P, A]) Act(actor P, action A, resolve ActionResolver[P, A]) (ActionOutcome[P], error) {
	if resolve == nil {
		return ActionOutcome[P]{}, ErrNilResolver
	}
	if e.IsOver() {
		return ActionOutcome[P]{}, ErrGameOver
	}

	actorIndex, ok := e.indexByPlayer[actor]
	if !ok {
		return ActionOutcome[P]{}, fmt.Errorf("%w: %v", ErrUnknownPlayer, actor)
	}
	if actorIndex != e.turnIndex {
		return ActionOutcome[P]{}, ErrWrongTurn
	}

	outcome, err := resolve(actor, action)
	if err != nil {
		return ActionOutcome[P]{}, err
	}

	if winner, hasWinner := outcome.Winner(); hasWinner {
		if _, exists := e.indexByPlayer[winner]; !exists {
			return ActionOutcome[P]{}, fmt.Errorf("%w: %v", ErrUnknownPlayer, winner)
		}
		e.winner = winner
		e.hasWinner = true
		return outcome, nil
	}

	if !outcome.KeepsTurn() {
		e.turnIndex = (e.turnIndex + 1) % len(e.order)
	}

	return outcome, nil
}
