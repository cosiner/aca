package aca

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

	stack := []*ACA{a}
	for l := len(stack); l > 0; l = len(stack) {
		topNode := stack[l-1]
		stack = stack[:l-1]

		for r, child := range topNode.children {
			node := topNode
			for node = node.failNode; ; node = node.failNode {
				failNode, has := node.children[r]
				if has && failNode != child {
					child.failNode = failNode
				}
				if node == a || child.failNode != nil {
					if child.failNode == nil {
						child.failNode = a
					}
					break
				}
			}
			stack = append(stack, child)
		}
	}
	return a
}

type Processor interface {
	Skip(r rune) bool
	Process(wholeStr string, index int, matched string) (continu bool)
}

func (a *ACA) Process(str string, processor Processor) {
	var curr = a
	for i, r := range str {
		if processor.Skip(r) {
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
				if !processor.Process(str, i, tmp.str) {
					return
				}
			}
		}
	}
}

type queryMatched struct {
	matched []string
}

func (m *queryMatched) Skip(rune) bool {
	return false
}

func (m *queryMatched) Process(wholeStr string, index int, matched string) bool {
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

func (p *queryHasContainedIn) Skip(rune) bool {
	return false
}

func (p *queryHasContainedIn) Process(string, int, string) bool {
	p.has = true
	return false
}

func (a *ACA) HasContainedIn(str string) bool {
	var p queryHasContainedIn
	a.Process(str, &p)
	return p.has
}

type replaceMatched struct {
	rs      []rune
	skips   map[rune]struct{}
	replace rune
}

func (p *replaceMatched) Skip(r rune) bool {
	_, has := p.skips[r]
	return has
}

func (p *replaceMatched) Process(wholeStr string, index int, matched string) bool {
	if p.rs == nil {
		p.rs = []rune(wholeStr)
	}

	for n, size, j := 0, len(matched), index; n < size && j >= 0; j-- {
		if !p.Skip(p.rs[j]) {
			p.rs[j] = p.replace
			n += 1
		}
	}
	return true
}

var _FLAG = struct{}{}

func NewRuneSet(rs string) map[rune]struct{} {
	if len(rs) == 0 {
		return nil
	}
	set := make(map[rune]struct{})
	for _, r := range rs {
		set[r] = _FLAG
	}
	return set
}

func (a *ACA) Replace(str string, replace rune, skips map[rune]struct{}) string {
	p := replaceMatched{
		replace: replace,
		skips:   skips,
	}
	a.Process(str, &p)
	if p.rs == nil {
		return str
	}
	return string(p.rs)
}
