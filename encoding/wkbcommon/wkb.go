package wkbcommon

import (
	"encoding/binary"
	"io"
)


const (
	pointType           uint32 = 1
	multiPointType      uint32 = 4
	lineStringType      uint32 = 2
	multiLineStringType uint32 = 5
	polygonType         uint32 = 3
	multiPolygonType    uint32 = 6
	ewkbType            uint32 = 0x20000000
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
