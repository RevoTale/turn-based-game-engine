// Package multi contains optional multi-grid and layer-registry utilities
// built on top of core grid2d primitives.
//
// Use this package when you need to manage many grids and/or many named layers
// across those grids. Registry manages grids and per-grid LayerSpace
// containers. LayerSpace manages layer creation and access.
//
// Single-board games can usually use grid2d.Grid and grid2d.SparseLayer
// directly.
package multi
