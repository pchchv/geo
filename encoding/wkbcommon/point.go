package wkbcommon

import (
	"math"

	"github.com/pchchv/geo"
)

func (e *Encoder) writePoint(p geo.Point, srid int) (err error) {
	if srid != 0 {
		e.order.PutUint32(e.buf, pointType|ewkbType)
		e.order.PutUint32(e.buf[4:], uint32(srid))
		if _, err = e.w.Write(e.buf[:8]); err != nil {
			return
		}
	} else {
		e.order.PutUint32(e.buf, pointType)
		if _, err = e.w.Write(e.buf[:4]); err != nil {
			return
		}
	}

	e.order.PutUint64(e.buf, math.Float64bits(p[0]))
	e.order.PutUint64(e.buf[8:], math.Float64bits(p[1]))
	_, err = e.w.Write(e.buf)
	return
}

func (e *Encoder) writeMultiPoint(mp geo.MultiPoint, srid int) (err error) {
	if err = e.writeTypePrefix(multiPointType, len(mp), srid); err != nil {
		return
	}

	for _, p := range mp {
		if err = e.Encode(p, 0); err != nil {
			return
		}
	}

	return nil
}
