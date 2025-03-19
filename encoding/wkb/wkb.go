package wkb

import (
	"encoding/binary"
	"io"

	"github.com/pchchv/geo/encoding/wkb/wkbcommon"
)

var DefaultByteOrder binary.ByteOrder = binary.LittleEndian // the order used for marshalling or encoding

// An Encoder will encode a geometry as WKB to the writer given at
// creation time.
type Encoder struct {
	e *wkbcommon.Encoder
}

// NewEncoder creates a new Encoder for the given writer.
func NewEncoder(w io.Writer) *Encoder {
	e := wkbcommon.NewEncoder(w)
	e.SetByteOrder(DefaultByteOrder)
	return &Encoder{e: e}
}
