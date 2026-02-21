# Trigger Rights Schema

Scope: trigger rights for `runtime/events`.

## Flow

```text
external/api or automation callback
  -> root command (`events.ExecuteCommand`)
    -> command handler
      -> emitted events (`ctx.Emit(events.Next(...))`)
```

## Feature Matrix

- Feature: External trigger
  Description: External code starts one root command tree.
  Example: `events.ExecuteCommand(runtime, command, payload)`

- Feature: Internal trigger
  Description: Command/event handlers emit only internal events.
  Example: `ctx.Emit(events.Next(eventToken, payload))`

- Feature: Execution mode
  Description: One command tree executes at a time per runtime.
  Example: two concurrent `events.ExecuteCommand(...)` calls are serialized per runtime.

## Not Allowed

- External code calling event handlers directly.
- Command/event handlers executing root commands.
