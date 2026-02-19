package turnbased

import "errors"

var (
	ErrEmptyPlayers    = errors.New("turn engine requires at least one player")
	ErrDuplicatePlayer = errors.New("turn engine has duplicate player")
	ErrUnknownPlayer   = errors.New("unknown player")
	ErrWrongTurn       = errors.New("wrong player turn")
	ErrGameOver        = errors.New("game is over")
	ErrNilResolver     = errors.New("action resolver is required")
	ErrInvalidOutcome  = errors.New("invalid action outcome")
)
