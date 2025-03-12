package mvt

import "github.com/pchchv/geo/encoding/mvt/vectortile"

type keyValueEncoder struct {
	Keys          []string
	keyMap        map[string]uint32
	Values        []*vectortile.Tile_Value
	valueMap      map[interface{}]uint32
	keySortBuffer []string
}

func newKeyValueEncoder() *keyValueEncoder {
	return &keyValueEncoder{
		keyMap:   make(map[string]uint32),
		valueMap: make(map[interface{}]uint32),
	}
}
