package aca

import (
	"container/list"
)

type treeNode struct {
	r         rune
	str       string
	runeCount int
	children  map[rune]*treeNode
	failNode  *treeNode
}

type ACA struct {
	cleaner Cleaner
	root    *treeNode
}

func New(cleaners ...Cleaner) *ACA {
	return &ACA{
		cleaner: groupCleaners(cleaners),
		root:    &treeNode{},
	}
}

func (a *ACA) addStr(str string) {
	rs := a.cleaner.Clean(newIndexedRunesByString(str))
	runeCount := len(rs)
	if runeCount == 0 {
		return
	}
	str = rs.String()

	curr := a.root
	for {
		r := rs[0]
		if curr.children == nil {
			curr.children = make(map[rune]*treeNode)
		}
		child, has := curr.children[r.Rune]
		if !has {
			child = &treeNode{
				r: r.Rune,
			}
			curr.children[r.Rune] = child
		}

		if len(rs) == 1 {
			child.str = str
			child.runeCount = runeCount
			break
		}
		rs = rs[1:]
		curr = child
	}
}

func (a *ACA) Add(strings ...string) *ACA {
	for _, str := range strings {
		a.addStr(str)
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
	Init(raw IndexedRunes)
	Process(a *ACA, runes IndexedRunes, matched string) (continu bool)
}

func (a *ACA) Process(str string, processor Processor) {
	var (
		curr     = a.root
		rawRunes = newIndexedRunesByString(str)
		rs       = a.cleaner.Clean(rawRunes.copy())
	)

	processor.Init(rawRunes)
	for i, r := range rs {
		var rooIterated bool
		for {
			child, has := curr.children[r.Rune]
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
				end := i + 1
				begin := end - tmp.runeCount
				if !processor.Process(a, rs[begin:end], tmp.str) {
					return
				}
			}
		}
	}
}

type queryMatched struct {
	raw     IndexedRunes
	matched []string
}

func (m *queryMatched) Init(raw IndexedRunes) {
	m.raw = raw
}

func (m *queryMatched) Process(a *ACA, runes IndexedRunes, matched string) bool {
	rs := make([]rune, len(runes))
	for i, r := range runes {
		rs[i] = m.raw[r.Index].Rune
	}

	m.matched = append(m.matched, string(rs))
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

func (p *queryHasContainedIn) Init(raw IndexedRunes) {}

func (p *queryHasContainedIn) Process(*ACA, IndexedRunes, string) bool {
	p.has = true
	return false
}

func (a *ACA) HasContainedIn(str string) bool {
	var p queryHasContainedIn
	a.Process(str, &p)
	return p.has
}

type replaceMatched struct {
	raw         IndexedRunes
	replaced    bool
	replacement rune
}

func (p *replaceMatched) Init(raw IndexedRunes) {
	p.raw = raw
}

func (p *replaceMatched) Process(a *ACA, runes IndexedRunes, matched string) bool {
	for _, r := range runes {
		p.raw[r.Index].Rune = p.replacement
		p.replaced = true
	}
	return true
}

func (a *ACA) Replace(str string, replacement rune) string {
	p := replaceMatched{
		replacement: replacement,
	}
	a.Process(str, &p)
	if !p.replaced {
		return str
	}
	return p.raw.String()
}
