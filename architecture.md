# Engine Architecture

## Goal

Keep runtime behavior deterministic, sequential, and domain-agnostic.

## Schema

```text
External/API
  -> Root command call
    -> Command
      -> Domain rules produce turn Decision (optional)
        -> turnbased.Step applies Decision (optional)
      -> Event chain (optional)
        -> State mutation
          -> Result
```

Rules:
- One command tree at a time per runtime.
- No parallel command-tree execution inside one runtime.
- Commands consume input and may emit events.
- Events do not receive payload.
- Runtime returns patch; root orchestration commits once.

## Public Features

- Feature: Root command execution
  Description: Execute root command tree and receive produced patch.
  Example: `patch, err := events.ExecuteCommand(runtime, state, command, input, newPatch)`

- Feature: Event emission
  Description: Chain internal reactions in deterministic FIFO order.
  Example: `emit(resolveShotEvent)`

- Feature: Turn decision step
  Description: Apply one validated turn decision and get typed delta.
  Example: `delta, err := turns.Step(actor, turnbased.NextTurn[Player]())`

- Feature: State lock scope
  Description: Keep one lock over the full command tree.
  Example: `_, err := store.Do(func(s *GameState, v uint64) error { ... })`
