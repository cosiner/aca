package aca

import (
	"container/list"
	"unicode/utf8"
)

type ACA struct {
	str      string
	r        rune
	children map[rune]*ACA
	failNode *ACA
}

func (a *ACA) addRunes(str string, rs []rune) {
	if len(rs) == 0 {
		return
	}

	curr := a
	for {
		r := rs[0]
		if curr.children == nil {
			curr.children = make(map[rune]*ACA)
		}
		child, has := curr.children[r]
		if !has {
			child = &ACA{r: r}
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

func (a *ACA) Build() *ACA {
	a.failNode = a

	queue := list.New()
	queue.PushBack(a)
	for queue.Len() > 0 {
		front := queue.Front()
		queue.Remove(front)
		topNode := front.Value.(*ACA)

		for r, child := range topNode.children {
			node := topNode
			for node = node.failNode; node != nil; node = node.failNode {
				failNode, has := node.children[r]
				if has && failNode != child {
					child.failNode = failNode
				}
				if node == a || child.failNode != nil {
					break
				}
			}
			if child.failNode == nil {
				child.failNode = a
			}
			queue.PushBack(child)
		}
	}
	return a
}

type Processor interface {
	Prepare(r rune) (res rune, valid bool)
	Process(rs []rune, index int, matched string) (continu bool)
}

func (a *ACA) Process(str string, processor Processor) {
	var (
		curr = a
		rs   = []rune(str)
	)
	for i, ru := range rs {
		r, valid := processor.Prepare(ru)
		if !valid {
			continue
		}

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

			if curr == a {
				if rooIterated {
					break
				}
				rooIterated = true
			}
		}
		for tmp := curr; tmp != a; tmp = tmp.failNode {
			if tmp.str != "" {
				if !processor.Process(rs, i, tmp.str) {
					return
				}
			}
		}
	}
}

type queryMatched struct {
	matched []string
}

func (m *queryMatched) Prepare(r rune) (rune, bool) {
	return r, true
}

func (m *queryMatched) Process(_ []rune, index int, matched string) bool {
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

func (p *queryHasContainedIn) Prepare(r rune) (rune, bool) {
	return r, true
}

func (p *queryHasContainedIn) Process([]rune, int, string) bool {
	p.has = true
	return false
}

func (a *ACA) HasContainedIn(str string) bool {
	var p queryHasContainedIn
	a.Process(str, &p)
	return p.has
}

type ReplaceOptions struct {
	Skips         RunSet
	Replacement   rune
	ReplaceSkip   bool
	CaseSensitive bool
}

type replaceMatched struct {
	rs []rune
	*ReplaceOptions
}

func (p *replaceMatched) Prepare(r rune) (rune, bool) {
	if p.Skips.Has(r) {
		return r, false
	}

	if p.CaseSensitive {
		return r, true
	}
	return ToLower(r), true
}

func (p *replaceMatched) Process(runes []rune, index int, matched string) bool {
	if p.rs == nil {
		p.rs = runes
	}
	for n, size, j := 0, utf8.RuneCountInString(matched), index; n < size && j >= 0; j-- {
		_, valid := p.Prepare(p.rs[j])
		if valid {
			n += 1
		}
		if valid || p.ReplaceSkip {
			p.rs[j] = p.Replacement
		}
	}
	return true
}

func (a *ACA) Replace(str string, options *ReplaceOptions) string {
	p := replaceMatched{
		ReplaceOptions: options,
	}
	a.Process(str, &p)
	if p.rs == nil {
		return str
	}
	return string(p.rs)
}
