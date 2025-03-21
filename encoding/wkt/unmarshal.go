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
	ErrNotWKT              = errors.New("wkt: invalid data")         // returned when unmarshalling WKT and the data is not valid
	ErrIncorrectGeometry   = errors.New("wkt: incorrect geometry")   // returned when unmarshalling WKT data into the wrong type
	ErrUnsupportedGeometry = errors.New("wkt: unsupported geometry") // returned when geometry type is not supported by this library
	singleParen            = regexp.MustCompile(`\)([\s|\t]*,[\s|\t]*)\(`)
	doubleParen            = regexp.MustCompile(`\)[\s|\t]*\)([\s|\t]*,[\s|\t]*)\([\s|\t]*\(`)
)

// Unmarshal returns a geometry by parsing the WKT string.
func Unmarshal(s string) (geo.Geometry, error) {
	s = trimSpace(s)
	prefix := upperPrefix(s)
	if bytes.HasPrefix(prefix, []byte("POINT")) {
		return unmarshalPoint(s)
	} else if bytes.HasPrefix(prefix, []byte("LINESTRING")) {
		return unmarshalLineString(s)
	} else if bytes.HasPrefix(prefix, []byte("POLYGON")) {
		return unmarshalPolygon(s)
	} else if bytes.HasPrefix(prefix, []byte("MULTIPOINT")) {
		return unmarshalMultiPoint(s)
	} else if bytes.HasPrefix(prefix, []byte("MULTILINESTRING")) {
		return unmarshalMultiLineString(s)
	} else if bytes.HasPrefix(prefix, []byte("MULTIPOLYGON")) {
		return unmarshalMultiPolygon(s)
	} else if bytes.HasPrefix(prefix, []byte("GEOMETRYCOLLECTION")) {
		return unmarshalCollection(s)
	} else {
		return nil, ErrUnsupportedGeometry
	}
}

// UnmarshalPoint returns the point represented by the wkt string.
// Returns ErrIncorrectGeometry if the wkt is not a point.
func UnmarshalPoint(s string) (geo.Point, error) {
	s = trimSpace(s)
	prefix := upperPrefix(s)
	if !bytes.HasPrefix(prefix, []byte("POINT")) {
		return geo.Point{}, ErrIncorrectGeometry
	}

	return unmarshalPoint(s)
}

// UnmarshalMultiPoint returns the multi-point represented by the wkt string.
// Returns ErrIncorrectGeometry if the wkt is not a multi-point.
func UnmarshalMultiPoint(s string) (geo.MultiPoint, error) {
	s = trimSpace(s)
	prefix := upperPrefix(s)
	if !bytes.HasPrefix(prefix, []byte("MULTIPOINT")) {
		return nil, ErrIncorrectGeometry
	}

	return unmarshalMultiPoint(s)
}

// UnmarshalLineString returns the linestring represented by the wkt string.
// Returns ErrIncorrectGeometry if the wkt is not a linestring.
func UnmarshalLineString(s string) (geo.LineString, error) {
	s = trimSpace(s)
	prefix := upperPrefix(s)
	if !bytes.HasPrefix(prefix, []byte("LINESTRING")) {
		return nil, ErrIncorrectGeometry
	}

	return unmarshalLineString(s)
}

// UnmarshalMultiLineString returns the multi-linestring represented by the wkt string.
// Returns ErrIncorrectGeometry if the wkt is not a multi-linestring.
func UnmarshalMultiLineString(s string) (geo.MultiLineString, error) {
	s = trimSpace(s)
	prefix := upperPrefix(s)
	if !bytes.HasPrefix(prefix, []byte("MULTILINESTRING")) {
		return nil, ErrIncorrectGeometry
	}

	return unmarshalMultiLineString(s)
}

// UnmarshalPolygon returns the polygon represented by the wkt string.
// Returns ErrIncorrectGeometry if the wkt is not a polygon.
func UnmarshalPolygon(s string) (geo.Polygon, error) {
	s = trimSpace(s)
	prefix := upperPrefix(s)
	if !bytes.HasPrefix(prefix, []byte("POLYGON")) {
		return nil, ErrIncorrectGeometry
	}

	return unmarshalPolygon(s)
}

// UnmarshalMultiPolygon returns the multi-polygon represented by the wkt string.
// Returns ErrIncorrectGeometry if the wkt is not a multi-polygon.
func UnmarshalMultiPolygon(s string) (geo.MultiPolygon, error) {
	s = trimSpace(s)
	prefix := upperPrefix(s)
	if !bytes.HasPrefix(prefix, []byte("MULTIPOLYGON")) {
		return nil, ErrIncorrectGeometry
	}

	return unmarshalMultiPolygon(s)
}

// UnmarshalCollection returns the geometry collection represented by the wkt string.
// Returns ErrIncorrectGeometry if the wkt is not a geometry collection.
func UnmarshalCollection(s string) (geo.Collection, error) {
	s = trimSpace(s)
	prefix := upperPrefix(s)
	if !bytes.HasPrefix(prefix, []byte("GEOMETRYCOLLECTION")) {
		return nil, ErrIncorrectGeometry
	}

	return unmarshalCollection(s)
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

func unmarshalLineString(s string) (geo.LineString, error) {
	if strings.EqualFold(s, "LINESTRING EMPTY") {
		return geo.LineString{}, nil
	}

	s, err := trimSpaceBrackets(s[10:])
	if err != nil {
		return nil, err
	}

	count := strings.Count(s, ",")
	ls := make(geo.LineString, 0, count+1)
	if err = splitOnComma(s, func(p string) error {
		tp, err := parsePoint(p)
		if err != nil {
			return err
		}

		ls = append(ls, tp)
		return nil
	}); err != nil {
		return nil, err
	}

	return ls, nil
}

func unmarshalMultiLineString(s string) (geo.MultiLineString, error) {
	if strings.EqualFold(s, "MULTILINESTRING EMPTY") {
		return geo.MultiLineString{}, nil
	}

	s, err := trimSpaceBrackets(s[15:])
	if err != nil {
		return nil, err
	}

	var tmls geo.MultiLineString
	if err = splitByRegexpYield(s, singleParen, func(i int) {
		tmls = make(geo.MultiLineString, 0, i)
	},
		func(ls string) (err error) {
			if ls, err = trimSpaceBrackets(ls); err != nil {
				return
			}

			count := strings.Count(ls, ",")
			tls := make(geo.LineString, 0, count+1)
			if err = splitOnComma(ls, func(p string) error {
				tp, err := parsePoint(p)
				if err != nil {
					return err
				}

				tls = append(tls, tp)
				return nil
			}); err != nil {
				return
			}

			tmls = append(tmls, tls)
			return nil
		},
	); err != nil {
		return nil, err
	}

	return tmls, nil
}

func unmarshalPolygon(s string) (poly geo.Polygon, err error) {
	if strings.EqualFold(s, "POLYGON EMPTY") {
		return geo.Polygon{}, nil
	}

	if s, err = trimSpaceBrackets(s[7:]); err != nil {
		return nil, err
	}

	if err = splitByRegexpYield(s, singleParen, func(i int) {
		poly = make(geo.Polygon, 0, i)
	},
		func(r string) (err error) {
			if r, err = trimSpaceBrackets(r); err != nil {
				return err
			}

			count := strings.Count(r, ",")
			ring := make(geo.Ring, 0, count+1)
			if err = splitOnComma(r, func(p string) error {
				tp, err := parsePoint(p)
				if err != nil {
					return err
				}

				ring = append(ring, tp)
				return nil
			}); err != nil {
				return
			}

			poly = append(poly, ring)
			return nil
		},
	); err != nil {
		return nil, err
	}

	return poly, nil
}

func unmarshalMultiPolygon(s string) (mpoly geo.MultiPolygon, err error) {
	if strings.EqualFold(s, "MULTIPOLYGON EMPTY") {
		return geo.MultiPolygon{}, nil
	}

	if s, err = trimSpaceBrackets(s[12:]); err != nil {
		return nil, err
	}

	if err = splitByRegexpYield(s, doubleParen, func(i int) {
		mpoly = make(geo.MultiPolygon, 0, i)
	},
		func(poly string) (err error) {
			if poly, err = trimSpaceBrackets(poly); err != nil {
				return
			}

			var tpoly geo.Polygon
			if err = splitByRegexpYield(poly, singleParen, func(i int) {
				tpoly = make(geo.Polygon, 0, i)
			},
				func(r string) (err error) {
					if r, err = trimSpaceBrackets(r); err != nil {
						return
					}

					count := strings.Count(r, ",")
					tr := make(geo.Ring, 0, count+1)
					if err = splitOnComma(r, func(s string) (err error) {
						if tp, err := parsePoint(s); err != nil {
							return err
						} else {
							tr = append(tr, tp)
							return nil
						}
					}); err != nil {
						return
					}

					tpoly = append(tpoly, tr)
					return nil
				},
			); err != nil {
				return
			}

			mpoly = append(mpoly, tpoly)
			return nil
		},
	); err != nil {
		return nil, err
	}

	return mpoly, nil
}

func unmarshalCollection(s string) (geo.Collection, error) {
	if strings.EqualFold(s, "GEOMETRYCOLLECTION EMPTY") {
		return geo.Collection{}, nil
	}

	if len(s) == 18 { // just GEOMETRYCOLLECTION
		return nil, ErrNotWKT
	}

	geometries := splitGeometryCollection(s[18:])
	if len(geometries) == 0 {
		return geo.Collection{}, nil
	}

	c := make(geo.Collection, 0, len(geometries))
	for _, g := range geometries {
		if len(g) == 0 {
			continue
		}

		tg, err := Unmarshal(g)
		if err != nil {
			return nil, err
		}

		c = append(c, tg)
	}

	return c, nil
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
