// Package automation provides generic bounded actor automation loops.
//
// The package is domain-agnostic and suitable for bot/NPC/background action
// chains where each pass may trigger follow-up work. It also provides a scoped
// scheduler to run automation sessions synchronously or asynchronously with
// per-scope serialization.
//
// This package is independent from runtime/events and can be used standalone.
// It owns timing and orchestration policy only, not domain state transitions.
// When composed with runtime/events, callbacks commonly call commands through
// events.Execute.
package automation
