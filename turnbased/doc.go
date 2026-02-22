// Package turnbased provides deterministic turn coordination for multiplayer
// games.
//
// The package validates turn ownership, applies domain decisions through Step,
// advances turn order, and finalizes match result state (winner or draw).
//
// Domain rules stay outside the package. The package only manages generic turn
// lifecycle semantics.
package turnbased
