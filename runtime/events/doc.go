// Package events provides deterministic in-memory command/event orchestration.
//
// One execution starts with a root command through ExecuteCommand.
// Command handlers consume input and may emit internal events. Events do not
// receive payload.
//
// Command/event handlers work through Context[S,P], which carries typed state
// and mutable patch plus event emitter.
//
// Runtime execution is serialized: one command tree runs at a time per runtime
// instance. ExecuteCommand returns patch to caller; store commit stays outside
// this package.
package events
