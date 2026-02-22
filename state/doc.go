// Package state provides typed single-writer state locking helpers.
//
// The package is intended for one match/session aggregate state where one
// command and its full event tree should run under one lock.
//
// Store.Do executes one mutation callback under lock and increments store
// version on success.
package state
