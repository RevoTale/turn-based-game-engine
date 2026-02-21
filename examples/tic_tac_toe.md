# Tic-Tac-Toe Example

Path: `engine/examples/tic_tac_toe/`

## Features

- Feature: Single command entrypoint
  Description: `Play(index)` executes one move deterministically.
  Example: `err := game.Play(4)`

- Feature: Locked state mutation
  Description: Full command tree runs under one state lock.
  Example: `store.Do(func(tx *state.Tx[gameState]) error { ... })`

- Feature: Internal event chaining
  Description: Move command emits follow-up events.
  Example: `ctx.Emit(events.Next(moveApplied, payload))`

- Feature: Shared pure rules
  Description: Domain rules are reusable functions; handlers stay thin.
  Example: `actor, outcome, err := applyMove(statePatch, move)`

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
