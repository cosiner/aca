package aca

import (
	"container/list"
	"unicode/utf8"
)

type treeNode struct {
	r        rune
	str      string
	children map[rune]*treeNode
	failNode *treeNode
}

type ACA struct {
	skips         RunSet
	caseSensitive bool
	root          *treeNode
}

func New(skips RunSet, caseSensitive bool) *ACA {
	return &ACA{
		skips:         skips,
		caseSensitive: caseSensitive,
		root:          &treeNode{},
	}
}

func (a *ACA) Skips() RunSet {
	return a.skips
}

func (a *ACA) CaseSensitive() bool {
	return a.caseSensitive
}

func (a *ACA) UpdateSkips(set RunSet) {
	a.skips = set
}

func (a *ACA) addRunes(str string, rs []rune) {
	rs = a.skips.Clean(rs)
	if len(rs) == 0 {
		return
	}
	rs = a.PrepareRunes(rs)
	str = string(rs)

	curr := a.root
	for {
		r := rs[0]
		if curr.children == nil {
			curr.children = make(map[rune]*treeNode)
		}
		child, has := curr.children[r]
		if !has {
			child = &treeNode{
				r: r,
			}
			curr.children[r] = child
		}

		if len(rs) == 1 {
			child.str = str
			break
		}
		rs = rs[1:]
		curr = child
	}
}

func (a *ACA) Add(strings ...string) *ACA {
	for _, str := range strings {
		a.addRunes(str, []rune(str))
	}

	return a
}

func (a *ACA) iterate(fn func(parent, child *treeNode)) {
	queue := list.New()
	queue.PushBack(a.root)
	for queue.Len() > 0 {
		front := queue.Front()
		queue.Remove(front)

		parent := front.Value.(*treeNode)
		for _, child := range parent.children {
			fn(parent, child)
			queue.PushBack(child)
		}
	}
}

func (a *ACA) Build() *ACA {
	a.root.failNode = a.root

	a.iterate(func(parent, child *treeNode) {
		node := parent
		for node = node.failNode; ; node = node.failNode {
			failNode, has := node.children[child.r]
			if has && failNode != child {
				child.failNode = failNode
			}
			if node == a.root || child.failNode != nil {
				break
			}
		}
		if child.failNode == nil {
			child.failNode = a.root
		}
	})

	return a
}

type Processor interface {
	Process(a *ACA, rs []rune, index int, matched string) (continu bool)
}

func (a *ACA) PrepareRune(r rune) rune {
	if a.caseSensitive {
		return r
	}
	return ToLower(r)
}

func (a *ACA) PrepareRunes(rs []rune) []rune {
	for i, r := range rs {
		rs[i] = r
	}
	return rs
}

func (a *ACA) Process(str string, processor Processor) {
	var (
		curr = a.root
		rs   = []rune(str)
	)
	for i, r := range rs {
		if a.skips.Has(r) {
			continue
		}

		r = a.PrepareRune(r)
		var rooIterated bool
		for {
			child, has := curr.children[r]
			if has {
				curr = child
			} else {
				curr = curr.failNode
			}
			if has {
				break
			}

			if curr == a.root {
				if rooIterated {
					break
				}
				rooIterated = true
			}
		}
		for tmp := curr; tmp != a.root; tmp = tmp.failNode {
			if tmp.str != "" {
				if !processor.Process(a, rs, i, tmp.str) {
					return
				}
			}
		}
	}
}

type queryMatched struct {
	matched []string
}

func (m *queryMatched) Process(_ *ACA, _ []rune, index int, matched string) bool {
	m.matched = append(m.matched, matched)
	return true
}

func (a *ACA) Match(str string) []string {
	var p queryMatched
	a.Process(str, &p)
	return p.matched
}

type queryHasContainedIn struct {
	has bool
}

func (p *queryHasContainedIn) Process(*ACA, []rune, int, string) bool {
	p.has = true
	return false
}

func (a *ACA) HasContainedIn(str string) bool {
	var p queryHasContainedIn
	a.Process(str, &p)
	return p.has
}

type replaceMatched struct {
	rs          []rune
	replacement rune
	replaceSkip bool
}

func (p *replaceMatched) Process(a *ACA, runes []rune, index int, matched string) bool {
	if p.rs == nil {
		p.rs = runes
	}
	for n, size, j := 0, utf8.RuneCountInString(matched), index; n < size && j >= 0; j-- {
		skipped := a.Skips().Has(p.rs[j])
		if !skipped {
			n += 1
		}
		if !skipped || p.replaceSkip {
			p.rs[j] = p.replacement
		}
	}
	return true
}

func (a *ACA) Replace(str string, replacement rune, replaceSkip bool) string {
	p := replaceMatched{
		replacement: replacement,
		replaceSkip: replaceSkip,
	}
	a.Process(str, &p)
	if p.rs == nil {
		return str
	}
	return string(p.rs)
}
