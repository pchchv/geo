package mvt

import "strconv"

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
