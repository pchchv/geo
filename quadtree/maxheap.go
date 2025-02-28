package quadtree

import "github.com/pchchv/geo"

type heapItem struct {
	point    geo.Pointer
	distance float64
}
