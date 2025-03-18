package wkbcommon

import (
	"encoding/binary"
	"io"
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
