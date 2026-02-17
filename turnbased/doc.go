// Package turnbased tracks player turns for multiplayer turn-based games.
//
// Use it when you want one place that:
// - checks whose turn it is
// - rotates turn order
// - marks the game as finished when a winner appears
//
// Game-specific rules stay in callbacks, so this package can be reused by
// different games.
package turnbased
