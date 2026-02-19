package tictactoe

import (
	"strings"
	"sync"

	"github.com/RevoTale/turn-based-game-engine/grid2d"
	"github.com/RevoTale/turn-based-game-engine/runtime/events"
	"github.com/RevoTale/turn-based-game-engine/turnbased"
)

// Game is a deterministic single-match runtime.
type Game struct {
	mu      sync.Mutex
	grid    *grid2d.Grid
	board   *grid2d.SparseLayer[Player]
	turns   *turnbased.Engine[Player, Move]
	runtime *events.Runtime
	log     []string
	ev      gameEvents
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
		grid:  grid,
		board: board,
		turns: turns,
		log:   make([]string, 0, 16),
	}

	builder := events.NewBuilder()
	registeredEvents, err := g.registerEvents(builder)
	if err != nil {
		return nil, err
	}
	g.ev = registeredEvents

	runtime, err := builder.Build()
	if err != nil {
		return nil, err
	}
	g.runtime = runtime

	return g, nil
}

// Play applies one move in a thread-safe way.
func (g *Game) Play(index int) error {
	if g == nil {
		return ErrNilGame
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	return events.Execute(g.runtime, g.ev.play, Move{Index: index})
}

// Snapshot returns board, winner, and match result.
func (g *Game) Snapshot() ([9]Player, Player, turnbased.MatchResult) {
	g.mu.Lock()
	defer g.mu.Unlock()

	var board [9]Player
	g.board.ForEach(func(pos grid2d.Position, value Player) bool {
		index, err := g.grid.Index(pos)
		if err != nil {
			return false
		}
		board[int(index)] = value
		return true
	})
	winner, _ := g.turns.Winner()
	return board, winner, g.turns.Result()
}

// Logs returns a copy of runtime log entries.
func (g *Game) Logs() []string {
	g.mu.Lock()
	defer g.mu.Unlock()

	out := make([]string, len(g.log))
	copy(out, g.log)
	return out
}

// BoardString returns a printable board view.
func (g *Game) BoardString() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	cells := make([]string, 0, 3)
	for row := 0; row < g.grid.Height(); row++ {
		rowVals := make([]string, 0, 3)
		for col := 0; col < g.grid.Width(); col++ {
			mark, ok, err := g.board.Get(grid2d.Position{X: col, Y: row})
			if err != nil || !ok {
				rowVals = append(rowVals, ".")
				continue
			}
			rowVals = append(rowVals, string(mark))
		}
		cells = append(cells, strings.Join(rowVals, " "))
	}
	return strings.Join(cells, "\n")
}

func (g *Game) hasWinner(player Player) bool {
	lines := [8][3]grid2d.Position{
		{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
		{{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}},
		{{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}},
		{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}},
		{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}},
		{{X: 2, Y: 0}, {X: 2, Y: 1}, {X: 2, Y: 2}},
		{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2}},
		{{X: 2, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 2}},
	}
	for _, line := range lines {
		if g.cellOwner(line[0]) == player && g.cellOwner(line[1]) == player && g.cellOwner(line[2]) == player {
			return true
		}
	}
	return false
}

func (g *Game) isBoardFull() bool {
	return g.board.Len() == g.grid.CellCount()
}

func (g *Game) cellOwner(pos grid2d.Position) Player {
	value, ok, err := g.board.Get(pos)
	if err != nil || !ok {
		return 0
	}
	return value
}
