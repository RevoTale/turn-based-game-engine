// Package grid2d provides 2D board utilities for games.
//
// Use it when you need:
// - a grid with bounds checks
// - per-cell data storage that only keeps written cells
// - multiple grids managed by game or room id
//
// The package is generic, so each game chooses its own id and cell value types.
package grid2d
