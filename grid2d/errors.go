package grid2d

import "errors"

type ErrorCode string

const (
	CodeInvalidGridSize ErrorCode = "INVALID_GRID_SIZE"
	CodeGridTooLarge    ErrorCode = "GRID_TOO_LARGE"
	CodeOutOfBounds     ErrorCode = "OUT_OF_BOUNDS"
	CodeNilGrid         ErrorCode = "NIL_GRID"
	CodeNilGridSet      ErrorCode = "NIL_GRID_SET"
	CodeGridExists      ErrorCode = "GRID_EXISTS"
	CodeGridNotFound    ErrorCode = "GRID_NOT_FOUND"
	CodeLayerExists     ErrorCode = "LAYER_EXISTS"
)

var (
	ErrInvalidGridSize = errors.New("grid width and height must be positive")
	ErrGridTooLarge    = errors.New("grid cell count overflows int")
	ErrOutOfBounds     = errors.New("cell position is out of bounds")
	ErrNilGrid         = errors.New("grid is required")
	ErrNilGridSet      = errors.New("grid set is required")
	ErrGridExists      = errors.New("grid id already exists")
	ErrGridNotFound    = errors.New("grid not found")
	ErrLayerExists     = errors.New("layer key already exists")
)

// CodeOf returns a stable machine-readable code for sentinel engine errors.
func CodeOf(err error) ErrorCode {
	switch {
	case errors.Is(err, ErrInvalidGridSize):
		return CodeInvalidGridSize
	case errors.Is(err, ErrGridTooLarge):
		return CodeGridTooLarge
	case errors.Is(err, ErrOutOfBounds):
		return CodeOutOfBounds
	case errors.Is(err, ErrNilGrid):
		return CodeNilGrid
	case errors.Is(err, ErrNilGridSet):
		return CodeNilGridSet
	case errors.Is(err, ErrGridExists):
		return CodeGridExists
	case errors.Is(err, ErrGridNotFound):
		return CodeGridNotFound
	case errors.Is(err, ErrLayerExists):
		return CodeLayerExists
	default:
		return ""
	}
}
