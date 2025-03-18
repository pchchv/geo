package wkbcommon

import (
	"encoding/binary"
	"io"

	"github.com/pchchv/geo"
)

const (
	pointType              uint32 = 1
	multiPointType         uint32 = 4
	lineStringType         uint32 = 2
	multiLineStringType    uint32 = 5
	polygonType            uint32 = 3
	multiPolygonType       uint32 = 6
	geometryCollectionType uint32 = 7
	ewkbType               uint32 = 0x20000000
)

var DefaultByteOrder binary.ByteOrder = binary.LittleEndian // order used for marshalling or encoding

// Encoder encodes a geometry as (E)WKB for the
// writer specified at creation.
type Encoder struct {
	buf   []byte
	w     io.Writer
	order binary.ByteOrder
}

// NewEncoder creates a new Encoder for the given writer.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:     w,
		order: DefaultByteOrder,
	}
}

// Encode will write the geometry encoded as (E)WKB to the given writer.
func (e *Encoder) Encode(geom geo.Geometry, srid int) error {
	switch g := geom.(type) {
	case nil:
		// nil values should not write any data
		return nil
	case geo.MultiPoint:
		// empty sizes will still write an empty version of that type
		if g == nil {
			return nil
		}
	case geo.LineString:
		if g == nil {
			return nil
		}
	case geo.MultiLineString:
		if g == nil {
			return nil
		}
	case geo.Polygon:
		if g == nil {
			return nil
		}
	case geo.MultiPolygon:
		if g == nil {
			return nil
		}
	case geo.Collection:
		if g == nil {
			return nil
		}
	// deal with types that are not supported by wkb
	case geo.Ring:
		if g == nil {
			return nil
		}
		geom = geo.Polygon{g}
	case geo.Bound:
		geom = g.ToPolygon()
	}

	var b []byte
	if e.order == binary.LittleEndian {
		b = []byte{1}
	} else {
		b = []byte{0}
	}

	if _, err := e.w.Write(b); err != nil {
		return err
	}

	if e.buf == nil {
		e.buf = make([]byte, 16)
	}

	switch g := geom.(type) {
	case geo.Point:
		return e.writePoint(g, srid)
	case geo.MultiPoint:
		return e.writeMultiPoint(g, srid)
	case geo.LineString:
		return e.writeLineString(g, srid)
	case geo.MultiLineString:
		return e.writeMultiLineString(g, srid)
	case geo.Polygon:
		return e.writePolygon(g, srid)
	case geo.MultiPolygon:
		return e.writeMultiPolygon(g, srid)
	case geo.Collection:
		return e.writeCollection(g, srid)
	default:
		panic("unsupported type")
	}
}

// SetByteOrder overrides the default byte order set when the encoder was created.
func (e *Encoder) SetByteOrder(bo binary.ByteOrder) {
	e.order = bo
}

func (e *Encoder) writeTypePrefix(t uint32, l int, srid int) error {
	if srid == 0 {
		e.order.PutUint32(e.buf, t)
		e.order.PutUint32(e.buf[4:], uint32(l))
		_, err := e.w.Write(e.buf[:8])
		return err
	}

	e.order.PutUint32(e.buf, t|ewkbType)
	e.order.PutUint32(e.buf[4:], uint32(srid))
	e.order.PutUint32(e.buf[8:], uint32(l))
	_, err := e.w.Write(e.buf[:12])
	return err
}
