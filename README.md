# Turn-Based Game Engine

Reusable, domain-agnostic primitives for turn-based multiplayer game backends.

## Packages

- `turnbased/`: deterministic turn orchestration for multiplayer turn-based games.
- `grid2d/`: non-square 2D grid primitives with sparse typed state layers.
- `runtime/`: scheduler and conditional event stream abstractions.

## Contract

- Keep engine packages domain-agnostic.
- Expose stable primitives that game domains can adapt to.
- Keep behavior tests alongside the engine code.
