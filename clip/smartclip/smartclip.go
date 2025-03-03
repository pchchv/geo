package smartclip

import "github.com/pchchv/geo"

type endpoint struct {
	Point    geo.Point
	Start    bool
	Used     bool
	Side     uint8
	Index    int
	OtherEnd int
}
