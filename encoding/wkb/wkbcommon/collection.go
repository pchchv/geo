package wkbcommon

import (
	"io"

	"github.com/pchchv/geo"
)

func (e *Encoder) writeCollection(c geo.Collection, srid int) (err error) {
	if err = e.writeTypePrefix(geometryCollectionType, len(c), srid); err != nil {
		return
	}

	for _, geom := range c {
		if err = e.Encode(geom, 0); err != nil {
			return
		}
	}

	return nil
}

func readCollection(r io.Reader, order byteOrder, buf []byte) (geo.Collection, error) {
	num, err := readUint32(r, order, buf[:4])
	if err != nil {
		return nil, err
	}

	alloc := num
	if alloc > MaxMultiAlloc {
		// invalid data can come in here and allocate tons of memory.
		alloc = MaxMultiAlloc
	}
	result := make(geo.Collection, 0, alloc)

	d := NewDecoder(r)
	for i := 0; i < int(num); i++ {
		geom, _, err := d.Decode()
		if err != nil {
			return nil, err
		}

		result = append(result, geom)
	}

	return result, nil
}
