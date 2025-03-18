package wkbcommon

import "github.com/pchchv/geo"

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
