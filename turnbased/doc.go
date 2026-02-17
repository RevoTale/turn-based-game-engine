// Package turnbased provides deterministic turn coordination for multiplayer
// games.
//
// The package validates turn ownership, executes domain actions through a
// resolver callback, advances turn order, and finalizes winner state.
//
// Domain rules stay outside the package. The package only manages generic turn
// lifecycle semantics.
package turnbased
