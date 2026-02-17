# Engine Directory

`engine/` contains reusable game-engine code intended to be extracted into a standalone repository later.

## Current Modules

- `turnbased/`: generic deterministic turn orchestration for multiplayer turn-based games.
  - supports generic comparable player identifiers (`P comparable`)

## Contract

- Engine packages must stay domain-agnostic (no imports from `app`, `graph`, `battleship`, or transport layers).
- Engine APIs should expose stable primitives that game domains adapt to.
- Tests for engine behavior should live alongside engine code.

Detailed documentation is in `docs/engine/README.md`.
