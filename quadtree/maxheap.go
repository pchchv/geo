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

func (h *maxHeap) Push(point geo.Pointer, distance float64) {
	prevLen := len(*h)
	*h = (*h)[:prevLen+1]
	(*h)[prevLen].point = point
	(*h)[prevLen].distance = distance
	i := len(*h) - 1
	for i > 0 {
		up := ((i + 1) >> 1) - 1
		parent := (*h)[up]
		if distance < parent.distance {
			// parent is further so we're done fixing up the heap.
			break
		}

		// swap nodes
		(*h)[i].point = parent.point
		(*h)[i].distance = parent.distance
		(*h)[up].point = point
		(*h)[up].distance = distance
		i = up
	}
}
