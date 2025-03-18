package wkbcommon

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/pchchv/geo"
)

const (
	bigEndian              byteOrder = 0
	littleEndian           byteOrder = 1
	pointType              uint32    = 1
	multiPointType         uint32    = 4
	lineStringType         uint32    = 2
	multiLineStringType    uint32    = 5
	polygonType            uint32    = 3
	multiPolygonType       uint32    = 6
	geometryCollectionType uint32    = 7
	ewkbType               uint32    = 0x20000000
)

var DefaultByteOrder binary.ByteOrder = binary.LittleEndian // order used for marshalling or encoding

// ByteOrder represents little or big endian encoding.
// binary.ByteOrder is not used because that is an
// interface that leaks to the heap all over the place.
type byteOrder int

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

// Marshal encodes the geometry with the given byte order.
func Marshal(geom geo.Geometry, srid int, byteOrder ...binary.ByteOrder) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, GeomLength(geom, srid != 0)))
	e := NewEncoder(buf)
	if len(byteOrder) > 0 {
		e.SetByteOrder(byteOrder[0])
	}

	if err := e.Encode(geom, srid); err != nil {
		return nil, err
	}

	if buf.Len() == 0 {
		return nil, nil
	}

	return buf.Bytes(), nil
}

// MustMarshal encodes geometry and panics when an error occurs.
func MustMarshal(geom geo.Geometry, srid int, byteOrder ...binary.ByteOrder) []byte {
	if d, err := Marshal(geom, srid, byteOrder...); err != nil {
		panic(err)
	} else {
		return d
	}
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

// Decoder decodes (E)WKB geometry off of the stream.
type Decoder struct {
	r io.Reader
}

// NewDecoder will create a new (E)WKB decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

// GeomLength helps to perform pre-allocation during a marshal.
func GeomLength(geom geo.Geometry, ewkb bool) (ewkbExtra int) {
	if ewkb {
		ewkbExtra = 4
	}

	switch g := geom.(type) {
	case geo.Point:
		return 21 + ewkbExtra
	case geo.MultiPoint:
		return 9 + 21*len(g) + ewkbExtra
	case geo.LineString:
		return 9 + 16*len(g) + ewkbExtra
	case geo.MultiLineString:
		var sum int
		for _, ls := range g {
			sum += 9 + 16*len(ls)
		}

		return 9 + sum + ewkbExtra
	case geo.Polygon:
		var sum int
		for _, r := range g {
			sum += 4 + 16*len(r)
		}

		return 9 + sum + ewkbExtra
	case geo.MultiPolygon:
		var sum int
		for _, c := range g {
			sum += GeomLength(c, false)
		}

		return 9 + sum + ewkbExtra
	case geo.Collection:
		var sum int
		for _, c := range g {
			sum += GeomLength(c, false)
		}

		return 9 + sum + ewkbExtra
	}

	return 0
}

func unmarshalUint32(order byteOrder, buf []byte) uint32 {
	if order == littleEndian {
		return binary.LittleEndian.Uint32(buf)
	}

	return binary.BigEndian.Uint32(buf)
}

func readUint32(r io.Reader, order byteOrder, buf []byte) (uint32, error) {
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}

	return unmarshalUint32(order, buf), nil
}

func byteOrderType(buf []byte) (byteOrder, uint32, error) {
	if len(buf) < 6 {
		return 0, 0, ErrNotWKB
	}

	var order byteOrder
	switch buf[0] {
	case 0:
		order = bigEndian
	case 1:
		order = littleEndian
	default:
		return 0, 0, ErrNotWKBHeader
	}

	// the type which is 4 bytes
	typ := unmarshalUint32(order, buf[1:])
	return order, typ, nil
}

func readByteOrderType(r io.Reader, buf []byte) (order byteOrder, typ uint32, srid int, err error) {
	// the byte order is the first byte
	if _, err := r.Read(buf[:1]); err != nil {
		return 0, 0, 0, err
	}

	if buf[0] == 0 {
		order = bigEndian
	} else if buf[0] == 1 {
		order = littleEndian
	} else {
		return 0, 0, 0, ErrNotWKB
	}

	// the type which is 4 bytes
	typ, err = readUint32(r, order, buf[:4])
	if err != nil {
		return 0, 0, 0, err
	} else if typ&ewkbType == 0 {
		return order, typ, 0, nil
	}

	if u, err := readUint32(r, order, buf[:4]); err != nil {
		return 0, 0, 0, err
	} else {
		srid = int(u)
	}

	return order, typ & 0x0ff, srid, nil
}
