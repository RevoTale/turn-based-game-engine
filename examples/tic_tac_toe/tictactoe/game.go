package tictactoe

import (
	"strings"

	"github.com/RevoTale/turn-based-game-engine/grid2d"
	"github.com/RevoTale/turn-based-game-engine/runtime/events"
	"github.com/RevoTale/turn-based-game-engine/state"
	"github.com/RevoTale/turn-based-game-engine/turnbased"
)

// Game is a deterministic single-match runtime.
type Game struct {
	store   *state.Store[gameState]
	runtime *events.Runtime
	ev      gameEvents
}

type gameState struct {
	grid  *grid2d.Grid
	board *grid2d.SparseLayer[Player]
	turns *turnbased.Engine[Player, Move]
	log   []string
}

// NewGame creates one game runtime.
func NewGame() (*Game, error) {
	grid, err := grid2d.NewGrid(3, 3)
	if err != nil {
		return nil, err
	}
	board, err := grid2d.NewSparseLayer[Player](grid)
	if err != nil {
		return nil, err
	}

	turns, err := turnbased.New[Player, Move]([]Player{PlayerX, PlayerO}, PlayerX)
	if err != nil {
		return nil, err
	}

	g := &Game{
		store: state.New(gameState{
			grid:  grid,
			board: board,
			turns: turns,
			log:   make([]string, 0, 16),
		}),
		runtime: events.NewRuntime(),
	}

	registeredEvents, err := g.registerEvents()
	if err != nil {
		return nil, err
	}
	g.ev = registeredEvents

	return g, nil
}

// Play applies one move in a thread-safe way.
func (g *Game) Play(index int) error {
	if g == nil {
		return ErrNilGame
	}

	_, err := g.store.Do(func(st *gameState, _ uint64) error {
		patch, executeErr := events.ExecuteCommand(g.runtime, st, g.ev.play, Move{Index: index}, func() *gamePatch {
			return &gamePatch{
				log: make([]string, 0, 4),
			}
		})
		if executeErr != nil {
			return executeErr
		}
		return applyPatch(st, patch)
	})
	return err
}

// Snapshot returns board, winner, and match result.
func (g *Game) Snapshot() ([9]Player, Player, turnbased.MatchResult) {
	var board [9]Player
	var winner Player
	result := turnbased.MatchResultOngoing
	if g == nil {
		return board, winner, result
	}

	_ = g.store.View(func(st *gameState, _ uint64) {
		st.board.ForEach(func(pos grid2d.Position, value Player) bool {
			index, err := st.grid.Index(pos)
			if err != nil {
				return false
			}
			board[int(index)] = value
			return true
		})
		winner, _ = st.turns.Winner()
		result = st.turns.Result()
	})
	return board, winner, result
}

// Logs returns a copy of runtime log entries.
func (g *Game) Logs() []string {
	if g == nil {
		return nil
	}

	var out []string
	_ = g.store.View(func(st *gameState, _ uint64) {
		out = make([]string, len(st.log))
		copy(out, st.log)
	})
	return out
}

// BoardString returns a printable board view.
func (g *Game) BoardString() string {
	if g == nil {
		return ""
	}

	cells := make([]string, 0, 3)
	_ = g.store.View(func(st *gameState, _ uint64) {
		for row := 0; row < st.grid.Height(); row++ {
			rowVals := make([]string, 0, 3)
			for col := 0; col < st.grid.Width(); col++ {
				mark, ok, err := st.board.Get(grid2d.Position{X: col, Y: row})
				if err != nil || !ok {
					rowVals = append(rowVals, ".")
					continue
				}
				rowVals = append(rowVals, string(mark))
			}
			cells = append(cells, strings.Join(rowVals, " "))
		}
	})
	return strings.Join(cells, "\n")
}
