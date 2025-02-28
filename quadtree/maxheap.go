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

// Pop returns the "greatest" item in the list.
// The returned item should not be saved across push/pop operations.
func (h *maxHeap) Pop() {
	lastItem := (*h)[len(*h)-1]
	(*h) = (*h)[:len(*h)-1]
	mh := (*h)
	if len(mh) == 0 {
		return
	}

	var i int
	// move the last item to the top and reset the heap
	mh[0].point = lastItem.point
	mh[0].distance = lastItem.distance
	for {
		right := (i + 1) << 1
		left := right - 1
		childIndex := i
		child := mh[childIndex]
		// swap with biggest child
		if left < len(mh) && child.distance < mh[left].distance {
			childIndex = left
			child = mh[left]
		}

		if right < len(mh) && child.distance < mh[right].distance {
			childIndex = right
			child = mh[right]
		}

		// non bigger, so quit
		if childIndex == i {
			break
		}

		// swap the nodes
		mh[i].point = child.point
		mh[i].distance = child.distance
		mh[childIndex].point = lastItem.point
		mh[childIndex].distance = lastItem.distance
		i = childIndex
	}
}
