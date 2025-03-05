package simplify

import (
	"encoding/json"
	"os"

	"github.com/pchchv/geo"
)

func benchmarkData() (ls geo.LineString) {
	f, err := os.Open("testdata/lisbon2portugal.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var points []float64
	if err = json.NewDecoder(f).Decode(&points); err != nil {
		panic(err)
	}

	for i := 0; i < len(points); i += 2 {
		ls = append(ls, geo.Point{points[i], points[i+1]})
	}

	return
}
