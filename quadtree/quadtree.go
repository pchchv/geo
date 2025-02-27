package quadtree

import "github.com/pchchv/geo"

// FilterFunc is a function that filters the points to search for.
type FilterFunc func(p geo.Pointer) bool

// node represents a node of the quad tree.
// Each node stores a Value and has links to its 4 children.
type node struct {
	Value    geo.Pointer
	Children [4]*node
}

// Quadtree implements a two-dimensional recursive
// spatial subdivision of geo.Pointers.
// This implementation uses rectangular partitions.
type Quadtree struct {
	bound geo.Bound
	root  *node
}

// New creates a new quadtree for the given bound.
// Added points must be within this bound.
func New(bound geo.Bound) *Quadtree {
	return &Quadtree{bound: bound}
}
