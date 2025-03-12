package mvt

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/pchchv/geo/geojson"
)

func convertIntID(i int) *uint64 {
	if i < 0 {
		return nil
	}

	v := uint64(i)
	return &v
}

func convertID(id interface{}) *uint64 {
	switch id := id.(type) {
	case nil:
		return nil
	case int:
		return convertIntID(id)
	case int8:
		return convertIntID(int(id))
	case int16:
		return convertIntID(int(id))
	case int32:
		return convertIntID(int(id))
	case int64:
		return convertIntID(int(id))
	case uint:
		v := uint64(id)
		return &v
	case uint8:
		v := uint64(id)
		return &v
	case uint16:
		v := uint64(id)
		return &v
	case uint32:
		v := uint64(id)
		return &v
	case uint64:
		v := uint64(id)
		return &v
	case float32:
		return convertIntID(int(id))
	case float64:
		return convertIntID(int(id))
	case string:
		if i, err := strconv.Atoi(id); err == nil {
			return convertIntID(i)
		}
	}

	return nil
}

func encodeProperties(kve *keyValueEncoder, properties geojson.Properties) ([]uint32, error) {
	tags := make([]uint32, 0, 2*len(properties))
	kve.keySortBuffer = kve.keySortBuffer[:0]
	for k := range properties {
		kve.keySortBuffer = append(kve.keySortBuffer, k)
	}
	sort.Strings(kve.keySortBuffer)

	for _, k := range kve.keySortBuffer {
		ki := kve.Key(k)
		if vi, err := kve.Value(properties[k]); err != nil {
			return nil, fmt.Errorf("property %s: %v", k, err)
		} else {
			tags = append(tags, ki, vi)
		}
	}

	return tags, nil
}
