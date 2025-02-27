package length

import (
	"math"

	"github.com/pchchv/geo"
)

func Distance(a, b geo.Point) float64 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	return math.Sqrt(dx*dx + dy*dy)
}
