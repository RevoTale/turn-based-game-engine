// Package turnbased provides reusable turn orchestration for server-authoritative
// multiplayer turn-based games.
//
// Design goals:
// - deterministic turn progression
// - domain-agnostic action execution via callbacks
// - explicit game-over lifecycle with validated winner
//
// This package intentionally avoids dependencies on app transport, storage, or
// game-domain models so it can be extracted to a standalone repository.
package turnbased
