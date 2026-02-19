package turnbased

// MatchResult is the terminal status of a match.
type MatchResult uint8

const (
	// MatchResultOngoing means the match is not finished.
	MatchResultOngoing MatchResult = iota
	// MatchResultWinner means the match finished with a winner.
	MatchResultWinner
	// MatchResultDraw means the match finished as a draw.
	MatchResultDraw
)

// ActionOutcome describes what happens to turn state after an action.
type ActionOutcome[P comparable] struct {
	keepTurn bool
	result   MatchResult
	winner   P
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
	o.result = MatchResultWinner
	o.winner = winner
	return o
}

// WithDraw marks action result as game over with draw state.
func (o ActionOutcome[P]) WithDraw() ActionOutcome[P] {
	o.result = MatchResultDraw
	var zero P
	o.winner = zero
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
	if o.result != MatchResultWinner {
		var zero P
		return zero, false
	}
	return o.winner, true
}

// Draw reports whether this outcome ends the match as a draw.
func (o ActionOutcome[P]) Draw() bool {
	return o.result == MatchResultDraw
}

// Result returns terminal result of this outcome.
func (o ActionOutcome[P]) Result() MatchResult {
	return o.result
}
