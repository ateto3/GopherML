package parser

// RecursiveDescentParser is a simple top-down parser for small grammars.
type RecursiveDescentParser struct {
	grammar *Grammar
}

// NewRecursiveDescentParser creates a recursive descent parser.
func NewRecursiveDescentParser(grammar *Grammar) *RecursiveDescentParser {
	return &RecursiveDescentParser{grammar: grammar}
}

// Parse returns all parses built by recursive descent.
func (p *RecursiveDescentParser) Parse(tokens []string) []Tree {
	if p.grammar == nil || len(tokens) == 0 {
		return nil
	}
	return NewChartParser(p.grammar).Parse(tokens)
}
