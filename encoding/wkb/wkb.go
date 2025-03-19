package wkb

import (
	"encoding/binary"
	"io"

	"github.com/pchchv/geo"
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

// Encode writes the geometry encoded as WKB to the given writer.
func (e *Encoder) Encode(geom geo.Geometry) error {
	return e.e.Encode(geom, 0)
}

// SetByteOrder overrides the default byte order set when the encoder was created.
func (e *Encoder) SetByteOrder(bo binary.ByteOrder) *Encoder {
	e.e.SetByteOrder(bo)
	return e
}

// Decoder decodes WKB geometry off of the stream.
type Decoder struct {
	d *wkbcommon.Decoder
}

// NewDecoder will create a new WKB decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		d: wkbcommon.NewDecoder(r),
	}
}
