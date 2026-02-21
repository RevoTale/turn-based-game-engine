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
  Example: `actor, outcome, err := rules.ApplyMove(patch, move)`

- Feature: Local event payloads
  Description: Keep payloads command-local; avoid one global shared payload type.
  Example: `events.Command[playCommand]`, `events.Event[moveApplied]`

## Public API

### `state`

- Feature: `state.New(initial)`
  Description: Create a typed state store for one runtime.
  Example: `store := state.New(GameState{})`

- Feature: `(*Store).Do(run)`
  Description: Run one state mutation under lock.
  Example: `_, err := store.Do(func(tx *state.Tx[GameState]) error { return nil })`

- Feature: `(*Store).View(read)`
  Description: Read state consistently under lock.
  Example: `err := store.View(func(s *GameState, v uint64) {})`

- Feature: `(*Store).Version()`
  Description: Get committed version number.
  Example: `v, err := store.Version()`

### `runtime/events`

- Feature: `events.NewBuilder()`
  Description: Start runtime command/event registration.
  Example: `builder := events.NewBuilder()`

- Feature: `events.RegisterCommand(...)`
  Description: Register one root command handler.
  Example: `play, err := events.RegisterCommand(builder, handlePlay, events.Hooks[Move]{})`

- Feature: `events.RegisterEvent(...)`
  Description: Register one internal event handler.
  Example: `moved, err := events.RegisterEvent(builder, handleMoved, events.Hooks[Moved]{})`

- Feature: `(*Builder).Build()`
  Description: Build immutable runtime.
  Example: `runtime, err := builder.Build()`

- Feature: `events.ExecuteCommand(...)`
  Description: Execute one root command tree.
  Example: `err := events.ExecuteCommand(runtime, play, Move{Index: 4})`

- Feature: `ctx.Emit(events.Next(...))`
  Description: Emit follow-up internal event from command/event handler.
  Example: `ctx.Emit(events.Next(moved, Moved{Index: 4}))`

### `turnbased`

- Feature: `turnbased.New(order, first)`
  Description: Create deterministic turn engine.
  Example: `tb, err := turnbased.New[Player, Move]([]Player{"A", "B"}, "A")`

- Feature: `(*Engine).Act(actor, action, resolve)`
  Description: Validate turn owner and apply one domain action.
  Example: `outcome, err := tb.Act(player, move, resolve)`

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
