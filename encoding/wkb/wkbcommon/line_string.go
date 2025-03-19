package wkbcommon

import (
	"errors"
	"io"
	"math"

	"github.com/pchchv/geo"
)

func (e *Encoder) writeLineString(ls geo.LineString, srid int) (err error) {
	if err = e.writeTypePrefix(lineStringType, len(ls), srid); err != nil {
		return
	}

	for _, p := range ls {
		e.order.PutUint64(e.buf, math.Float64bits(p[0]))
		e.order.PutUint64(e.buf[8:], math.Float64bits(p[1]))
		if _, err = e.w.Write(e.buf); err != nil {
			return
		}
	}

	return nil
}

func (e *Encoder) writeMultiLineString(mls geo.MultiLineString, srid int) (err error) {
	if err = e.writeTypePrefix(multiLineStringType, len(mls), srid); err != nil {
		return
	}

	for _, ls := range mls {
		if err = e.Encode(ls, 0); err != nil {
			return
		}
	}

	return nil
}

func readLineString(r io.Reader, order byteOrder, buf []byte) (geo.LineString, error) {
	num, err := readUint32(r, order, buf[:4])
	if err != nil {
		return nil, err
	}

	alloc := num
	if alloc > MaxPointsAlloc {
		// invalid data can come in here and allocate tons of memory.
		alloc = MaxPointsAlloc
	}
	result := make(geo.LineString, 0, alloc)

	for i := 0; i < int(num); i++ {
		p, err := readPoint(r, order, buf)
		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, nil
}

func readMultiLineString(r io.Reader, order byteOrder, buf []byte) (geo.MultiLineString, error) {
	num, err := readUint32(r, order, buf[:4])
	if err != nil {
		return nil, err
	}

	alloc := num
	if alloc > MaxMultiAlloc {
		// invalid data can come in here and allocate tons of memory.
		alloc = MaxMultiAlloc
	}
	result := make(geo.MultiLineString, 0, alloc)

	for i := 0; i < int(num); i++ {
		lOrder, typ, _, err := readByteOrderType(r, buf)
		if err != nil {
			return nil, err
		}

		if typ != lineStringType {
			return nil, errors.New("expect multilines to contains lines, did not find a line")
		}

		ls, err := readLineString(r, lOrder, buf)
		if err != nil {
			return nil, err
		}

		result = append(result, ls)
	}

	return result, nil
}

func unmarshalLineString(order byteOrder, data []byte) (geo.LineString, error) {
	ps, err := unmarshalPoints(order, data)
	if err != nil {
		return nil, err
	}

	return geo.LineString(ps), nil
}

func unmarshalMultiLineString(order byteOrder, data []byte) (geo.MultiLineString, error) {
	if len(data) < 4 {
		return nil, ErrNotWKB
	}

	num := unmarshalUint32(order, data)
	data = data[4:]
	alloc := num
	if alloc > MaxMultiAlloc {
		// invalid data can come in here and allocate tons of memory.
		alloc = MaxMultiAlloc
	}

	result := make(geo.MultiLineString, 0, alloc)
	for i := 0; i < int(num); i++ {
		ls, _, err := ScanLineString(data)
		if err != nil {
			return nil, err
		}

		data = data[16*len(ls)+9:]
		result = append(result, ls)
	}

	return result, nil
}
