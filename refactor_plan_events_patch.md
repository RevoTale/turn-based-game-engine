# Engine Refactor Plan

## Requirements

1. Events do not accept payload.
2. Commands accept input and can emit events.
3. Events mutate only command-local patch (delta) through narrow typed interfaces.
4. Full state can be read in events only through adapter-provided read-only methods.
5. Command execution returns patch.
6. Root orchestration commits patch once.
7. Runtime is deterministic and sequential: one command tree at a time.
8. Rollback mechanics are removed from engine core.
9. Automation remains an optional external command trigger layer.
10. `turnbased` and `grid2d` stay independent.

## Phases

1. Refactor `runtime/events` to remove hooks and payload emissions.
2. Move to patch-return command execution + single root commit model.
3. Simplify `state` to lock + apply APIs (no rollback hooks).
4. Simplify `grid2d/multi` to one registry surface.
5. Update docs/examples/tests to match new APIs.

## Done

1. `runtime/events` has no payload-bearing event path.
2. No rollback APIs remain in `state`.
3. Tic-tac-toe uses patch-return execution and single commit.
4. `grid2d/multi` no longer exposes overlapping `LayerSpace` + `Registry`.
5. `go test ./...` in `engine` is green.
