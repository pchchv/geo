package quadtree

import "github.com/pchchv/geo"

type heapItem struct {
	point    geo.Pointer
	distance float64
}

// maxHeap is used for the knearest list.
// Needed a way to keep the furthest point from the query point in the list,
// so maxHeap is used.
// When a point is found closer than the furthest point,
// remove furthest and add the new point to the heap.
type maxHeap []heapItem
