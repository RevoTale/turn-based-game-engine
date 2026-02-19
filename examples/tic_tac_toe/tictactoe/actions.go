package tictactoe

import "errors"

// Player is a tic-tac-toe mark.
type Player rune

const (
	// PlayerX is the first player mark.
	PlayerX Player = 'X'
	// PlayerO is the second player mark.
	PlayerO Player = 'O'
)

var (
	// ErrGameFinished indicates move submission after match completion.
	ErrGameFinished = errors.New("match is finished")
	// ErrMoveBounds indicates move index outside board bounds.
	ErrMoveBounds = errors.New("move index out of bounds")
	// ErrCellBusy indicates move to an occupied cell.
	ErrCellBusy = errors.New("cell already occupied")
	// ErrNilGame indicates method call on nil game pointer.
	ErrNilGame = errors.New("game is nil")
)

// Move is one user move command.
type Move struct {
	Index int
}
