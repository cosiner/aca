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

func New(strs []string, cleaner Cleaner) *ACA {
	a := &ACA{
		cleaner: cleaner,
		root:    &treeNode{},
	}
	for _, str := range strs {
		a.addStr(str)
	}
	a.build()
	return a
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

func (a *ACA) build() {
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
}

func (a *ACA) Process(str string, processor Processor) {
	var (
		curr     = a.root
		rawRunes = newIndexedRunesByString(str)
		rs       = a.cleaner.Clean(rawRunes.copy())
	)

	processor.Init(str, rawRunes)
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
				if !processor.Process(rs[begin:end], tmp.str) {
					return
				}
			}
		}
	}
}
