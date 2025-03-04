package simplify

import (
	"math"

	"github.com/pchchv/geo"
)

var _ geo.Simplifier = &VisvalingamSimplifier{}

// VisvalingamSimplifier is a reducer that performs the vivalingham algorithm.
type VisvalingamSimplifier struct {
	Threshold float64
	ToKeep    int // If 0, the default will be 2 for line, 3 for non-closed rings and 4 for closed rings.
	// The intent is to maintain valid geometry after simplification,
	// however it is still possible for the simplification to create self-intersections.
}

// Visvalingam creates a new VisvalingamSimplifier.
// If minPointsToKeep is 0 the algorithm will keep at least 2 points for lines,
// 3 for non-closed rings and 4 for closed rings.
// However it is still possible for the simplification to create self-intersections.
func Visvalingam(threshold float64, minPointsToKeep int) *VisvalingamSimplifier {
	return &VisvalingamSimplifier{
		Threshold: threshold,
		ToKeep:    minPointsToKeep,
	}
}

// VisvalingamThreshold runs the Visvalingam-Whyatt algorithm
// removing triangles whose area is below the threshold.
// Will keep at least 2 points for lines, 3 for non-closed rings and 4 for closed rings.
// The intent is to maintain valid geometry after simplification,
// however it is still possible for the simplification to create self-intersections.
func VisvalingamThreshold(threshold float64) *VisvalingamSimplifier {
	return Visvalingam(threshold, 0)
}

// VisvalingamKeep runs the Visvalingam-Wyatt algorithm,
// removing triangles of minimum area until the number of points `minPointsToKeep` is reached.
// If minPointsToKeep is 0, the algorithm will keep at least 2 points for lines,
// 3 for unclosed rings, and 4 for closed rings.
// However, it is still possible to create self-intersections when simplifying.
func VisvalingamKeep(minPointsToKeep int) *VisvalingamSimplifier {
	return Visvalingam(math.MaxFloat64, minPointsToKeep)
}

// Ring will simplify the ring using this simplifier.
func (s *VisvalingamSimplifier) Ring(r geo.Ring) geo.Ring {
	return ring(s, r)
}

// LineString will simplify the linestring using this simplifier.
func (s *VisvalingamSimplifier) LineString(ls geo.LineString) geo.LineString {
	return lineString(s, ls)
}

// MultiLineString will simplify the multi-linestring using this simplifier.
func (s *VisvalingamSimplifier) MultiLineString(mls geo.MultiLineString) geo.MultiLineString {
	return multiLineString(s, mls)
}

// Polygon will simplify the polygon using this simplifier.
func (s *VisvalingamSimplifier) Polygon(p geo.Polygon) geo.Polygon {
	return polygon(s, p)
}

// MultiPolygon will simplify the multi-polygon using this simplifier.
func (s *VisvalingamSimplifier) MultiPolygon(mp geo.MultiPolygon) geo.MultiPolygon {
	return multiPolygon(s, mp)
}

// Collection will simplify the collection using this simplifier.
func (s *VisvalingamSimplifier) Collection(c geo.Collection) geo.Collection {
	return collection(s, c)
}

// Simplify will run the simplification for any geometry type.
func (s *VisvalingamSimplifier) Simplify(g geo.Geometry) geo.Geometry {
	return simplify(s, g)
}

func (s *VisvalingamSimplifier) simplify(ls geo.LineString, area, wim bool) (geo.LineString, []int) {
	if len(ls) <= 1 {
		return ls, nil
	}

	toKeep := s.ToKeep
	if toKeep == 0 {
		if area {
			if ls[0] == ls[len(ls)-1] {
				toKeep = 4
			} else {
				toKeep = 3
			}
		} else {
			toKeep = 2
		}
	}

	var indexMap []int
	if len(ls) <= toKeep {
		if wim {
			// create identify map
			indexMap = make([]int, len(ls))
			for i := range ls {
				indexMap[i] = i
			}
		}
		return ls, indexMap
	}

	var removed int
	// edge cases checked, get on with it
	threshold := s.Threshold * 2                  // triangle area is doubled to save the multiply :)
	heap := minHeap(make([]*visItem, 0, len(ls))) // build the initial minheap linked list
	linkedListStart := &visItem{
		area:       math.Inf(1),
		pointIndex: 0,
	}
	heap.Push(linkedListStart)

	// internal path items
	items := make([]visItem, len(ls))
	previous := linkedListStart
	for i := 1; i < len(ls)-1; i++ {
		item := &items[i]
		item.area = doubleTriangleArea(ls, i-1, i, i+1)
		item.pointIndex = i
		item.previous = previous
		heap.Push(item)
		previous.next = item
		previous = item
	}

	// final item
	endItem := &items[len(ls)-1]
	endItem.area = math.Inf(1)
	endItem.pointIndex = len(ls) - 1
	endItem.previous = previous
	previous.next = endItem
	heap.Push(endItem)

	// run through the reduction process
	for len(heap) > 0 {
		current := heap.Pop()
		if current.area > threshold || len(ls)-removed <= toKeep {
			break
		}

		next := current.next
		previous := current.previous
		// remove current element from linked list
		previous.next = current.next
		next.previous = current.previous
		removed++

		// figure out the new areas
		if previous.previous != nil {
			area := doubleTriangleArea(ls,
				previous.previous.pointIndex,
				previous.pointIndex,
				next.pointIndex,
			)

			area = math.Max(area, current.area)
			heap.Update(previous, area)
		}

		if next.next != nil {
			area := doubleTriangleArea(ls,
				previous.pointIndex,
				next.pointIndex,
				next.next.pointIndex,
			)

			area = math.Max(area, current.area)
			heap.Update(next, area)
		}
	}

	var count int
	item := linkedListStart
	for item != nil {
		ls[count] = ls[item.pointIndex]
		count++

		if wim {
			indexMap = append(indexMap, item.pointIndex)
		}
		item = item.next
	}

	return ls[:count], indexMap
}

func doubleTriangleArea(ls geo.LineString, i1, i2, i3 int) float64 {
	a, b, c := ls[i1], ls[i2], ls[i3]
	return math.Abs((b[0]-a[0])*(c[1]-a[1]) - (b[1]-a[1])*(c[0]-a[0]))
}
