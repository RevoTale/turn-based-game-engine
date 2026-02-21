// Package state provides typed single-writer state transactions.
//
// The package is intended for one match/session aggregate state where one
// command and its full event tree should run under one lock.
//
// Store.Do executes one transactional callback under lock. The callback can
// register rollback hooks through Tx.BeforeChange to undo in-place mutations
// when the transaction fails.
package state
