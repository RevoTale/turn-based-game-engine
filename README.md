# Turn-Based Game Engine

Engine modules are independent and composable.

## Modules

- `turnbased`: deterministic turn ownership and match result state.
- `grid2d`: grid indexing and sparse layer state.
- `runtime/events`: deterministic command and internal event execution tree.
- `automation`: scheduling and delayed actor execution.

## Independence

- `runtime/events` does not depend on `automation`.
- `automation` does not depend on `runtime/events`.
- Both can be used alone.

## Composition (Optional)

Common pattern:

1. `automation` decides when an actor should act.
2. Automation callback triggers a root command via `events.Execute(...)`.
3. `runtime/events` runs the command and all emitted internal events deterministically.

## Responsibilities At A Glance

- `Command` in `runtime/events`: root intent entrypoint.
- `Event` in `runtime/events`: internal reaction node inside one execution tree.
- `automation`: timing/retry/delay policy for when to invoke root commands.

Detailed docs:

- `engine/turnbased/doc.go`
- `engine/grid2d/doc.go`
- `engine/runtime/doc.go`
- `engine/runtime/events/schema.md`
- `engine/examples/tic_tac_toe.md`
