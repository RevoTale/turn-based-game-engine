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

// Decision describes what the current action decided about turn flow.
//
// A decision is produced by domain rules and then applied via Engine.Step.
type Decision[P comparable] struct {
	keepTurn  bool
	result    MatchResult
	winner    P
	hasWinner bool
}

// Delta describes the turn engine state changes after one Step call.
type Delta[P comparable] struct {
	Actor          P
	PreviousPlayer P
	CurrentPlayer  P
	TurnChanged    bool
	Result         MatchResult
	Winner         P
	HasWinner      bool
	IsTerminal     bool
}

// NextTurn tells the engine to move turn to the next player.
func NextTurn[P comparable]() Decision[P] {
	return Decision[P]{}
}

// KeepTurn tells the engine that current player acts again.
func KeepTurn[P comparable]() Decision[P] {
	return Decision[P]{keepTurn: true}
}

// Win marks match as finished with a winner.
func Win[P comparable](winner P) Decision[P] {
	return Decision[P]{
		result:    MatchResultWinner,
		winner:    winner,
		hasWinner: true,
	}
}

// Draw marks match as finished with draw state.
func Draw[P comparable]() Decision[P] {
	return Decision[P]{
		result: MatchResultDraw,
	}
}

// KeepsTurn reports whether current actor keeps the turn.
func (d Decision[P]) KeepsTurn() bool {
	return d.keepTurn
}

// Winner returns winner when this decision ends the game with a winner.
//
// Second return value is false when decision has no winner.
func (d Decision[P]) Winner() (P, bool) {
	if d.result != MatchResultWinner || !d.hasWinner {
		var zero P
		return zero, false
	}
	return d.winner, true
}

// Draw reports whether this decision ends the match as a draw.
func (d Decision[P]) Draw() bool {
	return d.result == MatchResultDraw
}

// Result returns terminal result of this decision.
func (d Decision[P]) Result() MatchResult {
	return d.result
}
