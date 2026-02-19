# Trigger Rights Schema

Scope: trigger rights inside `runtime/events`.

## Flow

```text
External/API code ───────────────┐
automation callback (optional) ──┴──> Command (events.Execute)
                                      |
                                      v
                           Command handler + hooks
                                      |
                                      v
                          Event handler + hooks chain
                              (ctx.Emit + events.Next)
```

## Allowed

- External/API code -> `Command` via `events.Execute(...)`
- Automation callback (if composed) -> `Command` via `events.Execute(...)`
- Command handler/hook -> `Event` via `ctx.Emit(events.Next(...))`
- Event handler/hook -> `Event` via `ctx.Emit(events.Next(...))`

## Not Allowed

- External/API code -> `Event` directly
- Command/event handlers -> `Command` directly

## Responsibility

- `Command`: root intent entrypoint.
- `Event`: internal transition/reaction in one execution tree.
- `automation`: optional timing/orchestration layer that may call commands.
