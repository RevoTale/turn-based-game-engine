# Engine Architecture

## Goal

Keep runtime behavior deterministic, sequential, and domain-agnostic.

## Schema

```text
External/API
  -> Root command call
    -> Command
      -> Event chain (optional)
        -> State mutation
          -> Result
```

Rules:
- One command tree at a time per runtime.
- No parallel command-tree execution inside one runtime.
- Handlers emit events; handlers do not execute commands.

## Public Features

- Feature: Root command execution
  Description: Execute root command tree from external code.
  Example: `err := events.ExecuteCommand(runtime, command, payload)`

- Feature: Event emission
  Description: Chain internal reactions in deterministic order.
  Example: `ctx.Emit(events.Next(shipHit, payload))`

- Feature: State lock scope
  Description: Keep one lock over the full command tree.
  Example: `_, err := store.Do(func(tx *state.Tx[GameState]) error { ... })`
