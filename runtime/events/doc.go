// Package events provides deterministic in-memory command/event dispatch.
//
// This package is independent from automation and can be used standalone.
//
// One execution starts with a root Command through ExecuteCommand.
// Command and Event handlers can emit internal follow-up events through
// Context.Emit.
//
// Commands/events are registered once at boot time and receive internal
// incrementing ids. The built runtime is immutable.
//
// Dispatch uses a two-phase model: handlers and hooks enqueue emissions, and
// the runtime appends those emissions only after the current event stage
// finishes. Runtime execution is serialized: one command tree runs at a time
// per runtime instance.
package events
