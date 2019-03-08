package aca

import "unicode/utf8"

func newIndexedRunesByString(s string) IndexedRunes {
	rs := make(IndexedRunes, 0, utf8.RuneCountInString(s))
	var i int
	for _, r := range s {
		rs = append(rs, IndexedRune{
			Rune:  r,
			Index: i,
		})
		i++
	}
	return rs
}

type IndexedRune struct {
	Rune  rune
	Index int
}

type IndexedRunes []IndexedRune

func (c IndexedRunes) Runes() []rune {
	rs := make([]rune, len(c))
	for i := range c {
		rs[i] = c[i].Rune
	}
	return rs
}

func (c IndexedRunes) String() string {
	return string(c.Runes())
}

func (c IndexedRunes) copy() IndexedRunes {
	nc := make(IndexedRunes, len(c))
	copy(nc, c)
	return nc
}

type Cleaner interface {
	Clean(runes IndexedRunes) IndexedRunes
}

type skipsCleaner map[rune]struct{}

func NewSkipsCleaner(skips []rune) Cleaner {
	c := make(skipsCleaner)
	for _, r := range skips {
		c[r] = struct{}{}
	}
	return c
}

func (c skipsCleaner) Clean(rs IndexedRunes) IndexedRunes {
	var end int
	for i, r := range rs {
		if _, has := c[r.Rune]; !has {
			if i != end {
				rs[end] = rs[i]
			}
			end++
		}
	}
	return rs[:end]
}

type ignoreCaseCleaner struct{}

func NewIgnoreCaseCleaner() Cleaner {
	return ignoreCaseCleaner{}
}

func (c ignoreCaseCleaner) Clean(rs IndexedRunes) IndexedRunes {
	for i, r := range rs {
		if 'A' <= r.Rune && r.Rune <= 'Z' {
			rs[i].Rune = 'a' + r.Rune - 'A'
		}
	}
	return rs
}

type groupCleaners []Cleaner

func (g groupCleaners) Clean(rs IndexedRunes) IndexedRunes {
	for _, c := range g {
		rs = c.Clean(rs)
	}
	return rs
}

func GroupCleaners(cls ...Cleaner) Cleaner {
	return groupCleaners(cls)
}
