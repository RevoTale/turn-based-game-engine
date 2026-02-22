# Tic-Tac-Toe Example

Path: `engine/examples/tic_tac_toe/`

## Features

- Feature: Single command entrypoint
  Description: `Play(index)` executes one move deterministically.
  Example: `err := game.Play(4)`

- Feature: Locked state mutation
  Description: Full command tree runs under one state lock.
  Example: `store.Do(func(s *gameState, v uint64) error { ... })`

- Feature: Internal event chaining
  Description: Command/event views emit follow-up events without payload.
  Example: `emit(resolveMoveEvent)`

- Feature: Patch-return execution + single commit
  Description: Runtime returns patch; root applies patch to store once.
  Example: `patch, err := events.ExecuteCommand(runtime, state, playCmd, move, newPatch)`

- Feature: Shared pure rules
  Description: Domain rules are reusable functions; handlers stay thin.
  Example: `delta, nextTurns, writes, err := applyMove(state, patch, move)`

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
