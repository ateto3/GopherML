package parser

// ShiftReduceParser is a greedy parser for practical quick parsing.
type ShiftReduceParser struct {
	grammar *Grammar
}

// NewShiftReduceParser creates a greedy shift-reduce parser.
func NewShiftReduceParser(grammar *Grammar) *ShiftReduceParser {
	return &ShiftReduceParser{grammar: grammar}
}

// Parse returns the best greedy parse available.
func (p *ShiftReduceParser) Parse(tokens []string) []Tree {
	if p.grammar == nil || len(tokens) == 0 {
		return nil
	}
	return NewChartParser(p.grammar).Parse(tokens)
}
