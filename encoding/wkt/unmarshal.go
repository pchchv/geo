package wkt

import (
	"errors"
	"strconv"
	"strings"

	"github.com/pchchv/geo"
)

var ErrNotWKT = errors.New("wkt: invalid data") // returned when unmarshalling WKT and the data is not valid

func trimSpace(s string) string {
	if len(s) == 0 {
		return s
	}

	var start int
	for start = 0; start < len(s); start++ {
		if v := s[start]; v != ' ' && v != '\t' && v != '\n' {
			break
		}
	}

	var end int
	for end = len(s) - 1; end >= 0; end-- {
		if v := s[end]; v != ' ' && v != '\t' && v != '\n' {
			break
		}
	}

	if start >= end {
		return ""
	}

	return s[start : end+1]
}

// trimSpaceBrackets trim space and brackets
func trimSpaceBrackets(s string) (string, error) {
	s = trimSpace(s)
	if len(s) == 0 {
		return s, nil
	}

	if s[0] == '(' {
		s = s[1:]
	} else {
		return "", ErrNotWKT
	}

	if s[len(s)-1] == ')' {
		s = s[:len(s)-1]
	} else {
		return "", ErrNotWKT
	}

	return trimSpace(s), nil
}

// upperPrefix gets the ToUpper case of the first 20 chars.
func upperPrefix(s string) []byte {
	prefix := make([]byte, 20)
	for i := 0; i < 20 && i < len(s); i++ {
		if 'a' <= s[i] && s[i] <= 'z' {
			prefix[i] = s[i] - ('a' - 'A')
		} else {
			prefix[i] = s[i]
		}
	}

	return prefix
}

func unmarshalPoint(s string) (geo.Point, error) {
	s, err := trimSpaceBrackets(s[5:])
	if err != nil {
		return geo.Point{}, err
	}

	tp, err := parsePoint(s)
	if err != nil {
		return geo.Point{}, err
	}

	return tp, nil
}

// parsePoint pases point by (x y).
func parsePoint(s string) (p geo.Point, err error) {
	one, two, ok := strings.Cut(s, " ")
	if !ok {
		return geo.Point{}, ErrNotWKT
	}

	x, err := strconv.ParseFloat(one, 64)
	if err != nil {
		return geo.Point{}, ErrNotWKT
	}

	y, err := strconv.ParseFloat(two, 64)
	if err != nil {
		return geo.Point{}, ErrNotWKT
	}

	return geo.Point{x, y}, nil
}
