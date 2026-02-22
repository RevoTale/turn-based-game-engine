# Turn-Based Game Engine

## Principles

- One runtime represents one match/session.
- Execution is deterministic and sequential per runtime.
- Modules are independent and composable.
- Domain rules stay outside the engine.

## Runtime Schema

```text
external input
  -> root command call
    -> command handler
      -> event chain (optional)
        -> state update
          -> result/effects
```

Notes:
- One command tree is executed at a time per runtime.
- Parallelism is done by running many runtime instances, not by parallel command trees in one runtime.

## Composition Pattern

- Feature: Shared domain logic
  Description: Reuse pure rule functions across commands/events.
  Example: `decision, writes, err := rules.ApplyMove(state, patch, move)`

- Feature: Typed context-driven handlers
  Description: Handlers use `Context[S,P]` for state read, patch mutation, and event emission.
  Example: `func(ctx events.Context[ReadState, Patch], in Shoot) error`

## Public API

### `state`

- Feature: `state.New(initial)`
  Description: Create a typed state store for one runtime.
  Example: `store := state.New(GameState{})`

- Feature: `(*Store).Do(run)`
  Description: Run one state mutation under lock.
  Example: `_, err := store.Do(func(s *GameState, v uint64) error { return nil })`

- Feature: `(*Store).View(read)`
  Description: Read state consistently under lock.
  Example: `err := store.View(func(s *GameState, v uint64) {})`

- Feature: `(*Store).Version()`
  Description: Get committed version number.
  Example: `v, err := store.Version()`

### `runtime/events`

- Feature: `events.DefineCommand(...)`
  Description: Define one root command actor.
  Example: `play, err := events.DefineCommand(handlePlay)`

- Feature: `events.DefineEvent(...)`
  Description: Define one internal event actor.
  Example: `resolve, err := events.DefineEvent(handleResolve)`

- Feature: `events.ExecuteCommand(...)`
  Description: Execute one root command tree and return patch.
  Example: `patch, err := events.ExecuteCommand(runtime, state, play, Move{Index: 4}, newPatch)`

- Feature: `(events.Context).Emit(...)`
  Description: Queue follow-up internal event from command/event handler.
  Example: `ctx.Emit(resolveShotEvent)`

### `turnbased`

- Feature: `turnbased.New(order, first)`
  Description: Create deterministic turn engine.
  Example: `tb, err := turnbased.New[Player]([]Player{"A", "B"}, "A")`

- Feature: `(*Engine).Step(actor, decision)`
  Description: Validate turn owner and apply one domain decision.
  Example: `delta, err := tb.Step(player, turnbased.NextTurn[Player]())`

- Feature: `turnbased.NextTurn/KeepTurn/Win/Draw`
  Description: Build typed turn decisions from domain rules.
  Example: `decision := turnbased.Win(winnerID)`

- Feature: `(*Engine).CurrentPlayer()`
  Description: Get player that owns the current turn.
  Example: `p := tb.CurrentPlayer()`

- Feature: `(*Engine).Result()`
  Description: Get match result state.
  Example: `result := tb.Result()`

### `grid2d`

- Feature: `grid2d.NewGrid(w, h)`
  Description: Create bounded grid geometry.
  Example: `grid, err := grid2d.NewGrid(10, 10)`

- Feature: `grid2d.NewSparseLayer[T](grid)`
  Description: Create sparse typed cell storage on top of grid.
  Example: `shots, err := grid2d.NewSparseLayer[CellStatus](grid)`

- Feature: `(*SparseLayer).Set/Get`
  Description: Write/read one cell value.
  Example: `err = shots.Set(pos, Hit)`

### `automation` (optional)

- Feature: `automation.Run(config)`
  Description: Run bounded actor loop in-process.
  Example: `results, err := automation.Run(cfg)`

- Feature: `automation.NewScheduler(opts)`
  Description: Create sequential scheduler wrapper.
  Example: `sched := automation.NewScheduler[int, Player, Action](automation.SchedulerOptions{})`

- Feature: `Runtime.Run`
  Description: Run one scoped automation session (sequential).
  Example: `results, err := sched.Run(roomID, sessionCfg)`

## More Docs

- `engine/architecture.md`
- `engine/runtime/events/schema.md`
- `engine/examples/tic_tac_toe.md`
