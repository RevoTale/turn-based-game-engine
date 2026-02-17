# Engine Directory

`engine/` contains reusable game-engine code intended to be extracted into a standalone repository later.

## Current Modules

- `turnbased/`: generic deterministic turn orchestration for multiplayer turn-based games.
  - supports generic comparable player identifiers (`P comparable`)
- `grid2d/`: generic non-square 2D grid primitives with sparse typed state layers.
  - supports multiple independent grids with different dimensions
  - supports typed layer spaces (for example `bool`, `string`, numeric states)
- `runtime/`: scheduler and conditional event stream abstractions for app/runtime orchestration.
- `errcode/`: unified mapping from engine errors to stable machine-readable codes.

## Contract

- Engine packages must stay domain-agnostic (no imports from `app`, `graph`, `battleship`, or transport layers).
- Engine APIs should expose stable primitives that game domains adapt to.
- Tests for engine behavior should live alongside engine code.

Detailed documentation is in `docs/engine/README.md`.
