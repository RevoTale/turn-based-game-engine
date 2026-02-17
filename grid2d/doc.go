// Package grid2d provides reusable 2D grid primitives for turn-based engines.
//
// Design goals:
// - non-square grid support (width/height are independent)
// - sparse lazy state storage (state maps allocate only when written)
// - typed generic layers without runtime type assertions
// - multi-grid support for matches that operate on multiple boards
//
// Package grid2d is domain-agnostic and can be composed with turn orchestration
// from engine/turnbased.
package grid2d
