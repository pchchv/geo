package wkt

import "errors"

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
