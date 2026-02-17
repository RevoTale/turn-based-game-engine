package grid2d

import "errors"

type ErrorCode string

const (
	CodeInvalidGridSize     ErrorCode = "INVALID_GRID_SIZE"
	CodeGridTooLarge        ErrorCode = "GRID_TOO_LARGE"
	CodeOutOfBounds         ErrorCode = "OUT_OF_BOUNDS"
	CodeEmptyPositions      ErrorCode = "EMPTY_POSITIONS"
	CodeNotStraightLine     ErrorCode = "NOT_STRAIGHT_LINE"
	CodeInvalidLineAxis     ErrorCode = "INVALID_LINE_AXIS"
	CodeInvalidNeighborhood ErrorCode = "INVALID_NEIGHBORHOOD"
	CodeNilGrid             ErrorCode = "NIL_GRID"
	CodeNilGridSet          ErrorCode = "NIL_GRID_SET"
	CodeNilLayerRegistry    ErrorCode = "NIL_LAYER_REGISTRY"
	CodeGridExists          ErrorCode = "GRID_EXISTS"
	CodeGridNotFound        ErrorCode = "GRID_NOT_FOUND"
	CodeLayerExists         ErrorCode = "LAYER_EXISTS"
	CodeLayerNotFound       ErrorCode = "LAYER_NOT_FOUND"
)

var (
	ErrInvalidGridSize     = errors.New("grid width and height must be positive")
	ErrGridTooLarge        = errors.New("grid cell count overflows int")
	ErrOutOfBounds         = errors.New("cell position is out of bounds")
	ErrEmptyPositions      = errors.New("positions are required")
	ErrNotStraightLine     = errors.New("positions must form a straight horizontal or vertical line")
	ErrInvalidLineAxis     = errors.New("line axis is invalid")
	ErrInvalidNeighborhood = errors.New("neighborhood is invalid")
	ErrNilGrid             = errors.New("grid is required")
	ErrNilGridSet          = errors.New("grid set is required")
	ErrNilLayerRegistry    = errors.New("layer registry is required")
	ErrGridExists          = errors.New("grid id already exists")
	ErrGridNotFound        = errors.New("grid not found")
	ErrLayerExists         = errors.New("layer key already exists")
	ErrLayerNotFound       = errors.New("layer key not found")
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
	case errors.Is(err, ErrEmptyPositions):
		return CodeEmptyPositions
	case errors.Is(err, ErrNotStraightLine):
		return CodeNotStraightLine
	case errors.Is(err, ErrInvalidLineAxis):
		return CodeInvalidLineAxis
	case errors.Is(err, ErrInvalidNeighborhood):
		return CodeInvalidNeighborhood
	case errors.Is(err, ErrNilGrid):
		return CodeNilGrid
	case errors.Is(err, ErrNilGridSet):
		return CodeNilGridSet
	case errors.Is(err, ErrNilLayerRegistry):
		return CodeNilLayerRegistry
	case errors.Is(err, ErrGridExists):
		return CodeGridExists
	case errors.Is(err, ErrGridNotFound):
		return CodeGridNotFound
	case errors.Is(err, ErrLayerExists):
		return CodeLayerExists
	case errors.Is(err, ErrLayerNotFound):
		return CodeLayerNotFound
	default:
		return ""
	}
}
