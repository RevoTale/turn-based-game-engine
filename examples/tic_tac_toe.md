# Tic-Tac-Toe Runnable Example

A runnable standalone module is available at:

- `engine/examples/tic_tac_toe/`

It uses:

- `turnbased.Engine` for turn order and terminal result ownership (`winner` and `draw`)
- `runtime/events` with command-driven execution (`events.Execute`) and two-phase emissions (`ctx.Emit(events.Next(...))`)
- `grid2d.Grid` + `grid2d.SparseLayer` for board layout and per-cell state

This example does not use `automation`; it demonstrates the events runtime in
standalone mode.

## Run

```bash
cd engine/examples/tic_tac_toe
go run .
```

## Test

```bash
cd engine/examples/tic_tac_toe
go test ./...
```

## Files

- `engine/examples/tic_tac_toe/go.mod`: standalone module config (`replace` points to local engine)
- `engine/examples/tic_tac_toe/main.go`: showcase runner
- `engine/examples/tic_tac_toe/game_test.go`: behavior tests (win, draw, invalid moves)
- `engine/examples/tic_tac_toe/tictactoe/actions.go`: domain actions and error contract
- `engine/examples/tic_tac_toe/tictactoe/events.go`: event payloads and event registration
- `engine/examples/tic_tac_toe/tictactoe/game.go`: game runtime and board mechanics

## Design Notes

- `Play(index)` is the single command entrypoint.
- `Play` is self-locked (`sync.Mutex`) for deterministic single-writer updates.
- The root event validates/applies a move and returns follow-up events.
- Follow-up events are dispatched in deterministic order by the engine dispatcher.
