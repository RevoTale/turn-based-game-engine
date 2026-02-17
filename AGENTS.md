## Packages

- `turnbased/`: deterministic turn orchestration for multiplayer turn-based games.
- `grid2d/`: 2D grid primitives with sparse typed state layers.
- `runtime/`: scheduler and conditional event stream abstractions.

## Contract

- Keep engine packages domain-agnostic.
- Expose stable primitives that game domains can adapt to.
- Keep behavior tests alongside the engine code.
- Keep engine to be testable via `testing/synctest` internally and externally (by the higher level abtractions)

## Code Quality
- Follow the recommended Go programming language pattern
- Follow the recommendations of the "100 Go Mistakes and how to aavoid them" book
- Ensure tests are passing. 