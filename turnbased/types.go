package turnbased

// ActionOutcome describes what happens to turn state after an action.
type ActionOutcome[P comparable] struct {
	keepTurn  bool
	hasWinner bool
	winner    P
}

// ActionResolver runs game-specific action logic.
type ActionResolver[P comparable, A any] func(actor P, action A) (ActionOutcome[P], error)

// PassTurn tells the engine to move turn to the next player.
func PassTurn[P comparable]() ActionOutcome[P] {
	return ActionOutcome[P]{}
}

// KeepTurn tells the engine that current player acts again.
func KeepTurn[P comparable]() ActionOutcome[P] {
	return ActionOutcome[P]{keepTurn: true}
}

// WithWinner marks action result as game over with winner.
func (o ActionOutcome[P]) WithWinner(winner P) ActionOutcome[P] {
	o.hasWinner = true
	o.winner = winner
	return o
}

// KeepsTurn reports whether current actor keeps the turn.
func (o ActionOutcome[P]) KeepsTurn() bool {
	return o.keepTurn
}

// Winner returns winner when this outcome ends the game.
//
// Second return value is false when outcome has no winner.
func (o ActionOutcome[P]) Winner() (P, bool) {
	if !o.hasWinner {
		var zero P
		return zero, false
	}
	return o.winner, true
}
