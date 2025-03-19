package wkb

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/pchchv/geo"
	"github.com/pchchv/geo/encoding/wkb/wkbcommon"
)

var (
	ErrNotWKB              = errors.New("wkb: invalid data")              // returned when unmarshalling WKB and the data is not valid
	ErrIncorrectGeometry   = errors.New("wkb: incorrect geometry")        // returned when unmarshalling WKB data into the wrong type (e.g. linestring data into a point)
	ErrUnsupportedDataType = errors.New("wkb: scan value must be []byte") // the error returned when scanning a non-byte slice
	ErrUnsupportedGeometry = errors.New("wkb: unsupported geometry")      // returned when geometry type is not supported by this package
	commonErrorMap         = map[error]error{
		wkbcommon.ErrUnsupportedDataType: ErrUnsupportedDataType,
		wkbcommon.ErrNotWKB:              ErrNotWKB,
		wkbcommon.ErrNotWKBHeader:        ErrNotWKB,
		wkbcommon.ErrIncorrectGeometry:   ErrIncorrectGeometry,
		wkbcommon.ErrUnsupportedGeometry: ErrUnsupportedGeometry,
	}
	DefaultByteOrder binary.ByteOrder = binary.LittleEndian // the order used for marshalling or encoding
)

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

// Decode decodes the next geometry off of the stream.
func (d *Decoder) Decode() (geo.Geometry, error) {
	if g, _, err := d.d.Decode(); err != nil {
		return nil, mapCommonError(err)
	} else {
		return g, nil
	}
}

func mapCommonError(err error) error {
	if e, ok := commonErrorMap[err]; ok {
		return e
	}
	return err
}
