package wkbcommon

import (
	"math"

	"github.com/pchchv/geo"
)

func (e *Encoder) writePolygon(p geo.Polygon, srid int) (err error) {
	if err = e.writeTypePrefix(polygonType, len(p), srid); err != nil {
		return
	}

	for _, r := range p {
		e.order.PutUint32(e.buf, uint32(len(r)))
		if _, err = e.w.Write(e.buf[:4]); err != nil {
			return
		}

		for _, p := range r {
			e.order.PutUint64(e.buf, math.Float64bits(p[0]))
			e.order.PutUint64(e.buf[8:], math.Float64bits(p[1]))
			if _, err = e.w.Write(e.buf); err != nil {
				return
			}
		}
	}

	return nil
}

func (e *Encoder) writeMultiPolygon(mp geo.MultiPolygon, srid int) (err error) {
	if err = e.writeTypePrefix(multiPolygonType, len(mp), srid); err != nil {
		return
	}

	for _, p := range mp {
		if err = e.Encode(p, 0); err != nil {
			return
		}
	}

	return nil
}

func readPolygon(r io.Reader, order byteOrder, buf []byte) (geo.Polygon, error) {
	num, err := readUint32(r, order, buf[:4])
	if err != nil {
		return nil, err
	}

	alloc := num
	if alloc > MaxMultiAlloc {
		// invalid data can come in here and allocate tons of memory.
		alloc = MaxMultiAlloc
	}
	result := make(geo.Polygon, 0, alloc)

	for i := 0; i < int(num); i++ {
		ls, err := readLineString(r, order, buf)
		if err != nil {
			return nil, err
		}

		result = append(result, geo.Ring(ls))
	}

	return result, nil
}

func readMultiPolygon(r io.Reader, order byteOrder, buf []byte) (geo.MultiPolygon, error) {
	num, err := readUint32(r, order, buf[:4])
	if err != nil {
		return nil, err
	}

	alloc := num
	if alloc > MaxMultiAlloc {
		// invalid data can come in here and allocate tons of memory.
		alloc = MaxMultiAlloc
	}
	result := make(geo.MultiPolygon, 0, alloc)

	for i := 0; i < int(num); i++ {
		pOrder, typ, _, err := readByteOrderType(r, buf)
		if err != nil {
			return nil, err
		}

		if typ != polygonType {
			return nil, errors.New("expect multipolygons to contains polygons, did not find a polygon")
		}

		p, err := readPolygon(r, pOrder, buf)
		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, nil
}
