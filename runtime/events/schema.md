# Trigger Rights Schema

Scope: trigger rights for `runtime/events`.

## Flow

```text
external/api or automation callback
  -> root command (`events.ExecuteCommand`)
    -> command handler (consumes input)
      -> emitted events (`emit(event)`)
        -> event handlers (no payload)
          -> patch output
```

## Feature Matrix

- Feature: External trigger
  Description: External code starts one root command tree.
  Example: `patch, err := events.ExecuteCommand(runtime, state, command, input, newPatch)`

- Feature: Internal trigger
  Description: Command/event handlers emit only internal events.
  Example: `emit(resolveEvent)`

- Feature: Execution mode
  Description: One command tree executes at a time per runtime.
  Example: two concurrent `events.ExecuteCommand(...)` calls are serialized per runtime.

## Not Allowed

- External code calling event handlers directly.
- Command/event handlers executing root commands.
