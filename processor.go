package aca

type Processor interface {
	Init(str string, raw IndexedRunes)
	Process(runes IndexedRunes, matched string) (continu bool)
}

type QueryMatchedProcessor struct {
	raw     IndexedRunes
	matched []string
}

func (m *QueryMatchedProcessor) Init(str string, raw IndexedRunes) {
	m.raw = raw
}

func (m *QueryMatchedProcessor) Process(runes IndexedRunes, matched string) bool {
	rs := make([]rune, len(runes))
	for i, r := range runes {
		rs[i] = m.raw[r.Index].Rune
	}

	m.matched = append(m.matched, string(rs))
	return true
}

func (m *QueryMatchedProcessor) Result() []string {
	return m.matched
}

type QueryContainsProcessor struct {
	has bool
}

func (p *QueryContainsProcessor) Init(str string, raw IndexedRunes) {}

func (p *QueryContainsProcessor) Process(IndexedRunes, string) bool {
	p.has = true
	return false
}

func (p *QueryContainsProcessor) Result() bool { return p.has }

type ReplaceMatchedHandler struct {
	replacement rune

	str      string
	raw      IndexedRunes
	replaced bool
}

func NewReplaceMatchedHandler(replacement rune) *ReplaceMatchedHandler {
	return &ReplaceMatchedHandler{
		replacement: replacement,
	}
}

func (p *ReplaceMatchedHandler) Init(str string, raw IndexedRunes) {
	p.str = str
	p.raw = raw
}

func (p *ReplaceMatchedHandler) Process(runes IndexedRunes, matched string) bool {
	for _, r := range runes {
		p.raw[r.Index].Rune = p.replacement
		p.replaced = true
	}
	return true
}

func (p *ReplaceMatchedHandler) Result() string {
	if !p.replaced {
		return p.str
	}
	return p.raw.String()
}
