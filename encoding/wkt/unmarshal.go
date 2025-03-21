package wkt

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/pchchv/geo"
)

var (
	ErrNotWKT            = errors.New("wkt: invalid data")       // returned when unmarshalling WKT and the data is not valid
	ErrIncorrectGeometry = errors.New("wkt: incorrect geometry") // returned when unmarshalling WKT data into the wrong type
)

// UnmarshalPoint returns the point represented by the wkt string.
// Return ErrIncorrectGeometry if the wkt is not a point.
func UnmarshalPoint(s string) (geo.Point, error) {
	s = trimSpace(s)
	prefix := upperPrefix(s)
	if !bytes.HasPrefix(prefix, []byte("POINT")) {
		return geo.Point{}, ErrIncorrectGeometry
	}

	return unmarshalPoint(s)
}

// UnmarshalMultiPoint returns the multi-point represented by the wkt string.
// Return ErrIncorrectGeometry if the wkt is not a multi-point.
func UnmarshalMultiPoint(s string) (geo.MultiPoint, error) {
	s = trimSpace(s)
	prefix := upperPrefix(s)
	if !bytes.HasPrefix(prefix, []byte("MULTIPOINT")) {
		return nil, ErrIncorrectGeometry
	}

	return unmarshalMultiPoint(s)
}

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

func unmarshalMultiPoint(s string) (geo.MultiPoint, error) {
	if strings.EqualFold(s, "MULTIPOINT EMPTY") {
		return geo.MultiPoint{}, nil
	}

	s, err := trimSpaceBrackets(s[10:])
	if err != nil {
		return nil, err
	}

	count := strings.Count(s, ",")
	mp := make(geo.MultiPoint, 0, count+1)
	if err = splitOnComma(s, func(p string) error {
		p, err := trimSpaceBrackets(p)
		if err != nil {
			return err
		}

		tp, err := parsePoint(p)
		if err != nil {
			return err
		}

		mp = append(mp, tp)
		return nil
	}); err != nil {
		return nil, err
	}

	return mp, nil
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

// splitOnComma is optimized to split on the regex [\s|\t|\n]*,[\s|\t|\n]*
// i.e. comma with possible spaces on each side. e.g. '  ,  '
// We use a yield function because it
// was faster/used less memory than allocating an array of the results.
// In WKT points are separtated by commas,
// coordinates in points are separted by spaces e.g. 1 2,3 4,5 6,7 81 2,5 4
// is needed to split this and find each point.
func splitOnComma(s string, yield func(s string) error) error {
	// at is right after the previous space-comma-space match.
	// once a space-comma-space match is found,
	// go from 'at' to the start of the match,
	// that's the split that needs to be returned.
	// start of a space-comma-space section
	var at, start int
	// a space starts a section,
	// is needed to see a comma for it to be a valid section
	var sawSpace, sawComma bool
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			if !sawSpace {
				sawSpace = true
				start = i
			}

			sawComma = true
			continue
		}

		if v := s[i]; v == ' ' || v == '\t' || v == '\n' {
			if !sawSpace {
				sawSpace = true
				start = i
			}
			continue
		}

		if sawComma {
			if err := yield(s[at:start]); err != nil {
				return err
			}
			at = i
		}

		sawSpace, sawComma = false, false
	}

	return yield(s[at:])
}

// splitByRegexpYield splits the input by the regexp.
// The first callback can be used to initialize an array with the size of the result,
// the second is the callback with the matches.
// An yield function is uised because it was faster/used less memory than allocating an array of the results.
func splitByRegexpYield(s string, re *regexp.Regexp, set func(int), yield func(string) error) error {
	indexes := re.FindAllStringSubmatchIndex(s, -1)
	set(len(indexes) + 1)
	var start int
	for _, element := range indexes {
		if err := yield(s[start:element[2]]); err != nil {
			return err
		}

		start = element[3]
	}

	return yield(s[start:])
}

// splitGeometryCollection split GEOMETRYCOLLECTION to more geometry.
func splitGeometryCollection(s string) (r []string) {
	stack := make([]rune, 0)
	r = make([]string, 0)
	l := len(s)
	for i, v := range s {
		if !strings.Contains(string(stack), "(") {
			stack = append(stack, v)
			continue
		}

		if ('A' <= v && v < 'Z') || ('a' <= v && v < 'z') {
			t := string(stack)
			r = append(r, t[:len(t)-1])
			stack = make([]rune, 0)
			stack = append(stack, v)
			continue
		}

		if i == l-1 {
			r = append(r, string(stack))
			continue
		}

		stack = append(stack, v)
	}
	return
}
