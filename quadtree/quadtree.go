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
