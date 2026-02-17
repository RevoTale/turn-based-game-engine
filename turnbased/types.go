package turnbased

// ActionOutcome describes the turn consequences of a successfully resolved action.
type ActionOutcome[P comparable] struct {
	keepTurn  bool
	hasWinner bool
	winner    P
}

// ActionResolver runs domain-specific action logic and returns turn consequences.
type ActionResolver[P comparable, A any] func(actor P, action A) (ActionOutcome[P], error)

// PassTurn indicates that turn should move to the next player if game is not over.
func PassTurn[P comparable]() ActionOutcome[P] {
	return ActionOutcome[P]{}
}

// KeepTurn indicates that current player keeps the turn if game is not over.
func KeepTurn[P comparable]() ActionOutcome[P] {
	return ActionOutcome[P]{keepTurn: true}
}

// WithWinner marks game-over and attaches winner identity.
func (o ActionOutcome[P]) WithWinner(winner P) ActionOutcome[P] {
	o.hasWinner = true
	o.winner = winner
	return o
}

func (o ActionOutcome[P]) KeepsTurn() bool {
	return o.keepTurn
}

func (o ActionOutcome[P]) Winner() (P, bool) {
	if !o.hasWinner {
		var zero P
		return zero, false
	}
	return o.winner, true
}
