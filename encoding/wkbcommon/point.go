package wkbcommon

import (
	"encoding/binary"
	"errors"
	"io"
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

func readPoint(r io.Reader, order byteOrder, buf []byte) (p geo.Point, err error) {
	for i := 0; i < 2; i++ {
		if _, err = io.ReadFull(r, buf); err != nil {
			return geo.Point{}, err
		} else if order == littleEndian {
			p[i] = math.Float64frombits(binary.LittleEndian.Uint64(buf))
		} else {
			p[i] = math.Float64frombits(binary.BigEndian.Uint64(buf))
		}
	}

	return p, nil
}

func readMultiPoint(r io.Reader, order byteOrder, buf []byte) (geo.MultiPoint, error) {
	num, err := readUint32(r, order, buf[:4])
	if err != nil {
		return nil, err
	}

	alloc := num
	if alloc > MaxPointsAlloc {
		// invalid data can come in here and allocate tons of memory.
		alloc = MaxPointsAlloc
	}
	result := make(geo.MultiPoint, 0, alloc)

	for i := 0; i < int(num); i++ {
		pOrder, typ, _, err := readByteOrderType(r, buf)
		if err != nil {
			return nil, err
		}

		if typ != pointType {
			return nil, errors.New("expect multipoint to contains points, did not find a point")
		}

		p, err := readPoint(r, pOrder, buf)
		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, nil
}
