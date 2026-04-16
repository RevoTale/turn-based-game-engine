# AGENTS.md

## Overview
`engine/` is a standalone Go module and git submodule that provides domain-agnostic turn-based engine primitives used by Sea Battle Server, runnable examples, and future RevoTale games.

## Base Policy Links (Load First)
- Router: https://github.com/RevoTale/agent-docs/blob/main/AGENTS.router.md
- Common module: https://github.com/RevoTale/agent-docs/blob/main/modules/common/AGENTS.md
- Taskfile module: https://github.com/RevoTale/agent-docs/blob/main/modules/taskfile/AGENTS.md
- Go module: https://github.com/RevoTale/agent-docs/blob/main/modules/go/AGENTS.md

## Local Details
### Packages
- `automation/`: scoped sequential automation loops and schedulers.
- `state/`: typed single-writer state store and version helpers.
- `runtime/`: topic streams and delay policy helpers for integration layers.
- `runtime/events/`: deterministic command/event execution primitives.
- `turnbased/`: deterministic turn orchestration for multiplayer turn-based games.
- `grid2d/`: 2D grid primitives with sparse typed state layers.
- `examples/`: runnable reference integrations for engine consumers.

### Contract
- Keep engine packages domain-agnostic.
- Expose stable primitives that game domains can adapt to.
- Keep behavior tests alongside the engine code.
- Keep engine testable via `testing/synctest` internally and externally by higher-level abstractions.

### Code Quality
- Follow recommended Go programming language patterns.
- Follow recommendations from "100 Go Mistakes and How to Avoid Them".
- Ensure tests are passing.
- Use idiomatic Go following standard conventions.

### Documentation and comments
- Every exported (capitalized) function, type, constant, and variable must have a doc comment immediately preceding its declaration, with no intervening blank lines.
- Use clear, simple language and complete sentences.
- Code should explain what it does. Comments should explain non-obvious design decisions, potential side effects, or complex logic.
- Outdated comments are detrimental. Update comments whenever behavior changes.
