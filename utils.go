package aca

import "unicode"

var _FLAG = struct{}{}

type RunSet map[rune]struct{}

func NewRuneSet(rs string) RunSet {
	if rs == "" {
		return nil
	}
	set := make(map[rune]struct{})
	for _, r := range rs {
		set[r] = _FLAG
	}
	return RunSet(set)
}

func (s RunSet) Has(r rune) bool {
	_, has := s[r]
	return has
}

func (s RunSet) Clean(rs []rune) []rune {
	var prev int
	for _, r := range rs {
		if !s.Has(r) {
			rs[prev] = r
			prev += 1
		}
	}
	return rs[:prev]
}

func ToLower(r rune) rune {
	if r <= unicode.MaxASCII {
		return unicode.ToLower(r)
	}
	return r
}
