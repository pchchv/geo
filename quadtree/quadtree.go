package quadtree

import (
	"math"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/planar"
)

// FilterFunc is a function that filters the points to search for.
type FilterFunc func(p geo.Pointer) bool

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

// node represents a node of the quad tree.
// Each node stores a Value and has links to its 4 children.
type node struct {
	Value    geo.Pointer
	Children [4]*node
}

type visitor interface {
	// Bound returns the current relevant bound so that it is possible to drop irrelevant nodes from the search.
	Bound() *geo.Bound
	Visit(n *node)
	// Point returns the specific point searched for,
	// or null if none exists (i.e. boundary search).
	// This helps guide the search to the best child node first.
	Point() geo.Point
}

// visit provides a framework for walking the quad tree.
// Currently used by the `Find` and `InBound` functions.
type visit struct {
	visitor visitor
}

func newVisit(v visitor) *visit {
	return &visit{
		visitor: v,
	}
}

func (v *visit) Visit(n *node, left, right, bottom, top float64) {
	b := v.visitor.Bound()
	if left > b.Max[0] || right < b.Min[0] || bottom > b.Max[1] || top < b.Min[1] {
		return
	}

	if n.Value != nil {
		v.visitor.Visit(n)
	}

	if n.Children[0] == nil && n.Children[1] == nil && n.Children[2] == nil && n.Children[3] == nil {
		// no children check
		return
	}

	cx := (left + right) / 2.0
	cy := (bottom + top) / 2.0
	i := childIndex(cx, cy, v.visitor.Point())
	for j := i; j < i+4; j++ {
		if n.Children[j%4] == nil {
			continue
		}

		if k := j % 4; k == 0 {
			v.Visit(n.Children[0], left, cx, cy, top)
		} else if k == 1 {
			v.Visit(n.Children[1], cx, right, cy, top)
		} else if k == 2 {
			v.Visit(n.Children[2], left, cx, bottom, cy)
		} else if k == 3 {
			v.Visit(n.Children[3], cx, right, bottom, cy)
		}
	}
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

func childIndex(cx, cy float64, point geo.Point) (i int) {
	if point[1] <= cy {
		i = 2
	}

	if point[0] >= cx {
		i++
	}

	return
}
