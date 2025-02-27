package quadtree

import (
	"math"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/planar"
)

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

type visitor interface {
	// Bound returns the current relevant bound so we can prune irrelevant nodes
	// from the search. Using a pointer was benchmarked to be 5% faster than
	// having to copy the bound on return. go1.9
	Bound() *geo.Bound
	Visit(n *node)
	// Point should return the specific point being search for, or null if there
	// isn't one (ie. searching by bound). This helps guide the search to the
	// best child node first.
	Point() geo.Point
}

type findVisitor struct {
	point          geo.Point
	filter         FilterFunc
	closest        *node
	closestBound   *geo.Bound
	minDistSquared float64
}

func (v *findVisitor) Visit(n *node) {
	// skip this pointer if we have a filter and it doesn't match
	if v.filter != nil && !v.filter(n.Value) {
		return
	}

	point := n.Value.Point()
	if d := planar.DistanceSquared(point, v.point); d < v.minDistSquared {
		v.minDistSquared = d
		v.closest = n
		d = math.Sqrt(d)
		v.closestBound.Min[0] = v.point[0] - d
		v.closestBound.Max[0] = v.point[0] + d
		v.closestBound.Min[1] = v.point[1] - d
		v.closestBound.Max[1] = v.point[1] + d
	}
}

func (v *findVisitor) Bound() *geo.Bound {
	return v.closestBound
}

func (v *findVisitor) Point() geo.Point {
	return v.point
}
