// Package grid2d provides generic 2D grid primitives for game state storage.
//
// It includes bounded grid indexing, neighbor traversal, connected-component
// traversal, sparse per-cell layers, and straight-line helpers.
// Optional lock-aware layer wrappers are available for concurrent access.
//
// Optional multi-grid and layer-registry utilities are provided by the
// sibling package grid2d/multi.
//
// The package is domain-agnostic and designed to compose with higher-level game
// rules.
package grid2d
