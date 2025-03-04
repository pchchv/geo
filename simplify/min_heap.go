package simplify

type visItem struct {
	area       float64  // triangle area
	pointIndex int      // index of point in original path
	next       *visItem // to keep a virtual linked list to help restore the triangle areas when points are removed
	previous   *visItem
	index      int // internal index in heap, for removal and update
}

// minHeap creates a priority queue or min heap.
type minHeap []*visItem

func (h *minHeap) Push(item *visItem) {
	item.index = len(*h)
	*h = append(*h, item)
	h.up(item.index)
}

func (h *minHeap) Pop() *visItem {
	removed := (*h)[0]
	lastItem := (*h)[len(*h)-1]
	(*h) = (*h)[:len(*h)-1]
	if len(*h) > 0 {
		lastItem.index = 0
		(*h)[0] = lastItem
		h.down(0)
	}

	return removed
}

func (h minHeap) Update(item *visItem, area float64) {
	if item.area > area {
		// area got smaller
		item.area = area
		h.up(item.index)
	} else {
		// area got larger
		item.area = area
		h.down(item.index)
	}
}

func (h minHeap) up(i int) {
	object := h[i]
	for i > 0 {
		up := ((i + 1) >> 1) - 1
		parent := h[up]
		if parent.area <= object.area {
			break
		}

		// swap nodes
		parent.index = i
		h[i] = parent
		object.index = up
		h[up] = object
		i = up
	}
}

func (h minHeap) down(i int) {
	object := h[i]
	for {
		right := (i + 1) << 1
		left := right - 1
		down := i
		child := h[down]

		// swap with smallest child
		if left < len(h) && h[left].area < child.area {
			down = left
			child = h[down]
		}

		if right < len(h) && h[right].area < child.area {
			down = right
			child = h[down]
		}

		// non smaller, so quit
		if down == i {
			break
		}

		// swap the nodes
		child.index = i
		h[child.index] = child
		object.index = down
		h[down] = object
		i = down
	}
}
