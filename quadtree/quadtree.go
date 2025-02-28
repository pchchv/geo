package quadtree

import (
	"errors"
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

// Find returns the closest Value/Pointer in the quadtree.
// This function is thread safe.
// Multiple goroutines can read from a pre-created tree.
func (q *Quadtree) Find(p geo.Point) geo.Pointer {
	return q.Matching(p, nil)
}

// Matching returns the closest Value/Pointer in the
// quadtree for which the given filter function returns true.
// This function is thread safe.
// Multiple goroutines can read from a pre-created tree.
func (q *Quadtree) Matching(p geo.Point, f FilterFunc) geo.Pointer {
	if q.root == nil {
		return nil
	}

	b := q.bound
	v := &findVisitor{
		point:          p,
		filter:         f,
		closestBound:   &b,
		minDistSquared: math.MaxFloat64,
	}

	newVisit(v).Visit(q.root, q.bound.Min[0], q.bound.Max[0], q.bound.Min[1], q.bound.Max[1])

	if v.closest == nil {
		return nil
	}

	return v.closest.Value
}

// InBoundMatching returns a slice with all the pointers in the quadtree that are
// within the given bound and matching the give filter function.
// An optional buffer parameter is provided to
// allow for the reuse of result slice memory.
// This function is thread safe.
// Multiple goroutines can read from a pre-created tree.
func (q *Quadtree) InBoundMatching(buf []geo.Pointer, b geo.Bound, f FilterFunc) (p []geo.Pointer) {
	if q.root == nil {
		return nil
	}

	if buf != nil {
		p = buf[:0]
	}

	v := &inBoundVisitor{
		bound:    &b,
		pointers: p,
		filter:   f,
	}

	newVisit(v).Visit(q.root, q.bound.Min[0], q.bound.Max[0], q.bound.Min[1], q.bound.Max[1])

	return v.pointers
}

// Bound returns the bounds used for the quad tree.
func (q *Quadtree) Bound() geo.Bound {
	return q.bound
}

// InBound returns a slice with all the pointers in the
// quadtree that are within the given bound.
// An optional buffer parameter is provided to allow
// for the reuse of result slice memory.
// This function is thread safe.
// Multiple goroutines can read from a pre-created tree.
func (q *Quadtree) InBound(buf []geo.Pointer, b geo.Bound) []geo.Pointer {
	return q.InBoundMatching(buf, b, nil)
}

// Add puts an object into the quad tree,
// must be within the quadtree bounds.
// This function is not thread-safe,
// ie. multiple goroutines cannot insert into a single quadtree.
func (q *Quadtree) Add(p geo.Pointer) error {
	if p == nil {
		return nil
	}

	point := p.Point()
	if !q.bound.Contains(point) {
		// point is outside the bounds used to create the tree
		return errors.New("quadtree: point outside of bounds")
	}

	if q.root == nil {
		q.root = &node{
			Value: p,
		}
		return nil
	} else if q.root.Value == nil {
		q.root.Value = p
		return nil
	}

	q.add(q.root, p, p.Point(), q.bound.Min[0], q.bound.Max[0], q.bound.Min[1], q.bound.Max[1])

	return nil
}

// Remove will remove the pointer from the quadtree.
// By default it'll match using the points,
// but a FilterFunc can be provided for a more specific test
// if there are elements with the same point value in the tree.
// i. e.:
//
//	func(pointer geo.Pointer) {
//		return pointer.(*MyType).ID == lookingFor.ID
//	}
func (q *Quadtree) Remove(p geo.Pointer, eq FilterFunc) bool {
	if eq == nil {
		point := p.Point()
		eq = func(pointer geo.Pointer) bool {
			return point.Equal(pointer.Point())
		}
	}

	b := q.bound
	v := &findVisitor{
		point:          p.Point(),
		filter:         eq,
		closestBound:   &b,
		minDistSquared: math.MaxFloat64,
	}

	newVisit(v).Visit(q.root, q.bound.Min[0], q.bound.Max[0], q.bound.Min[1], q.bound.Max[1])
	if v.closest == nil {
		return false
	}

	v.closest.Value = nil
	// if v.closest is NOT a leaf node, values will be shuffled up into this node.
	// if v.closest IS a leaf node, the call is a no-op but we can't delete
	// the now empty node because we don't know the parent here.
	// Future adds will reuse this node if applicable.
	// Removing v.closest parent will cause this node to be removed,
	// but the parent will be a leaf with a nil value.
	removeNode(v.closest)
	return true
}

// KNearestMatching returns k closest Value/Pointer in the
// quadtree for which the given filter function returns true.
// This function is thread safe.
// Multiple goroutines can read from a pre-created tree.
// An optional buffer parameter allows the result slice memory to be reused.
// The points are returned in sorted order, nearest first.
// KNearestMatching allows defining a maximum distance in order to reduce search iterations.
func (q *Quadtree) KNearestMatching(buf []geo.Pointer, p geo.Point, k int, f FilterFunc, maxDistance ...float64) []geo.Pointer {
	if q.root == nil {
		return nil
	}

	b := q.bound
	v := &nearestVisitor{
		point:          p,
		filter:         f,
		k:              k,
		maxHeap:        make(maxHeap, 0, k+1),
		closestBound:   &b,
		maxDistSquared: math.MaxFloat64,
	}
	if len(maxDistance) > 0 {
		v.maxDistSquared = maxDistance[0] * maxDistance[0]
	}

	newVisit(v).Visit(q.root, q.bound.Min[0], q.bound.Max[0], q.bound.Min[1], q.bound.Max[1])

	// repack result
	if cap(buf) < len(v.maxHeap) {
		buf = make([]geo.Pointer, len(v.maxHeap))
	} else {
		buf = buf[:len(v.maxHeap)]
	}

	for i := len(v.maxHeap) - 1; i >= 0; i-- {
		buf[i] = v.maxHeap[0].point
		v.maxHeap.Pop()
	}

	return buf
}

// KNearest returns k closest Value/Pointer in the quadtree.
// KNearest is thread safe.
// Multiple goroutines can read from a pre-created tree.
// An optional buffer parameter is provided to allow for the reuse of result slice memory.
// The points are returned in a sorted order, nearest first.
// KNearest allows defining a maximum distance in order to reduce search iterations.
func (q *Quadtree) KNearest(buf []geo.Pointer, p geo.Point, k int, maxDistance ...float64) []geo.Pointer {
	return q.KNearestMatching(buf, p, k, nil, maxDistance...)
}

// add is the recursive search to find a place to add the point.
func (q *Quadtree) add(n *node, p geo.Pointer, point geo.Point, left, right, bottom, top float64) {
	var i int
	// figure which child of this internal node the point is in
	if cy := (bottom + top) / 2.0; point[1] <= cy {
		top = cy
		i = 2
	} else {
		bottom = cy
	}

	if cx := (left + right) / 2.0; point[0] >= cx {
		left = cx
		i++
	} else {
		right = cx
	}

	if n.Children[i] == nil {
		n.Children[i] = &node{Value: p}
		return
	} else if n.Children[i].Value == nil {
		n.Children[i].Value = p
		return
	}

	// proceed down to the child to see if it's a
	// leaf yet and can add the pointer there
	q.add(n.Children[i], p, point, left, right, bottom, top)
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

type inBoundVisitor struct {
	bound    *geo.Bound
	pointers []geo.Pointer
	filter   FilterFunc
}

func (v *inBoundVisitor) Visit(n *node) {
	if v.filter != nil && !v.filter(n.Value) {
		return
	}

	p := n.Value.Point()
	if v.bound.Min[0] > p[0] || v.bound.Max[0] < p[0] ||
		v.bound.Min[1] > p[1] || v.bound.Max[1] < p[1] {
		return

	}

	v.pointers = append(v.pointers, n.Value)
}

func (v *inBoundVisitor) Bound() *geo.Bound {
	return v.bound
}

func (v *inBoundVisitor) Point() (p geo.Point) {
	return
}

type nearestVisitor struct {
	point          geo.Point
	filter         FilterFunc
	k              int
	maxHeap        maxHeap
	closestBound   *geo.Bound
	maxDistSquared float64
}

func (v *nearestVisitor) Visit(n *node) {
	// skip this pointer if have a filter and it doesn't match
	if v.filter != nil && !v.filter(n.Value) {
		return
	}

	point := n.Value.Point()
	if d := planar.DistanceSquared(point, v.point); d < v.maxDistSquared {
		v.maxHeap.Push(n.Value, d)
		if len(v.maxHeap) > v.k {
			v.maxHeap.Pop()
			// actually this is a hack
			// know how heap works and obtain
			// top element without function call
			top := v.maxHeap[0]
			v.maxDistSquared = top.distance
			// have filled queue,
			// so start to restrict searching range
			d = math.Sqrt(top.distance)
			v.closestBound.Min[0] = v.point[0] - d
			v.closestBound.Max[0] = v.point[0] + d
			v.closestBound.Min[1] = v.point[1] - d
			v.closestBound.Max[1] = v.point[1] + d
		}
	}
}

func (v *nearestVisitor) Bound() *geo.Bound {
	return v.closestBound
}

func (v *nearestVisitor) Point() geo.Point {
	return v.point
}

type pointsQueueItem struct {
	point    geo.Pointer
	distance float64 // distance to point and priority inside the queue
	index    int     // point index in queue
}

type pointsQueue []pointsQueueItem

func newPointsQueue(capacity int) pointsQueue {
	// makes capacity+1 because is needed additional place for the greatest element
	return make([]pointsQueueItem, 0, capacity+1)
}

func (pq *pointsQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(pointsQueueItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *pointsQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq pointsQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq pointsQueue) Len() int { return len(pq) }

func (pq pointsQueue) Less(i, j int) bool {
	// want pop longest distances so Less was inverted
	return pq[i].distance > pq[j].distance
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

// removeNode is the recursive fixing up of the tree when we remove a node.
// It will pull up a child value into it's place.
// It will try to remove leaf nodes that are now empty,
// since their values got pulled up.
func removeNode(n *node) bool {
	var i int
	if n.Children[0] != nil {
		i = 0
	} else if n.Children[1] != nil {
		i = 1
	} else if n.Children[2] != nil {
		i = 2
	} else if n.Children[3] != nil {
		i = 3
	} else {
		// all children are nil, can remove
		// n.value ==  nil because it pulled up or removed by the caller
		return true
	}

	n.Value = n.Children[i].Value
	n.Children[i].Value = nil
	removeThisChild := removeNode(n.Children[i])
	if removeThisChild {
		n.Children[i] = nil
	}

	return false
}
