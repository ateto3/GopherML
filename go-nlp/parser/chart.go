package parser

import (
	"sort"
	"sync"
)

// ChartParser is a compact CKY-style parser for lexical + binary rules.
type ChartParser struct {
	grammar *Grammar
}

// NewChartParser creates a new parser for a grammar.
func NewChartParser(grammar *Grammar) *ChartParser {
	return &ChartParser{grammar: grammar}
}

// Parse returns all parse trees found for the input tokens.
func (p *ChartParser) Parse(tokens []string) []Tree {
	if len(tokens) == 0 || p.grammar == nil {
		return nil
	}
	lexical, binary := splitRules(p.grammar.Rules)
	n := len(tokens)
	chart := make([][]map[Symbol][]Tree, n)
	for i := range chart {
		chart[i] = make([]map[Symbol][]Tree, n+1)
	}
	for i, tok := range tokens {
		for _, rule := range lexical[tok] {
			cell := ensureCell(chart, i, i+1)
			cell[rule.LHS] = append(cell[rule.LHS], Tree{Label: rule.LHS, Word: tok})
		}
	}
	for span := 2; span <= n; span++ {
		var wg sync.WaitGroup
		var mu sync.Mutex
		for start := 0; start+span <= n; start++ {
			start := start
			span := span
			wg.Add(1)
			go func() {
				defer wg.Done()
				end := start + span
				local := make(map[Symbol][]Tree)
				for mid := start + 1; mid < end; mid++ {
					left := chart[start][mid]
					right := chart[mid][end]
					if left == nil || right == nil {
						continue
					}
					for _, rule := range binary {
						lhsTrees, lok := left[Symbol(rule.RHS[0])]
						rhsTrees, rok := right[Symbol(rule.RHS[1])]
						if !lok || !rok {
							continue
						}
						for _, lt := range lhsTrees {
							for _, rt := range rhsTrees {
								local[rule.LHS] = append(local[rule.LHS], Tree{Label: rule.LHS, Children: []Tree{lt, rt}})
							}
						}
					}
				}
				if len(local) > 0 {
					mu.Lock()
					cell := ensureCell(chart, start, end)
					for sym, trees := range local {
						cell[sym] = append(cell[sym], trees...)
					}
					mu.Unlock()
				}
			}()
		}
		wg.Wait()
	}
	result := chart[0][n][p.grammar.Start]
	sort.SliceStable(result, func(i, j int) bool { return result[i].String() < result[j].String() })
	return result
}

// ParseOne returns the best parse, or nil if no parse is available.
func (p *ChartParser) ParseOne(tokens []string) *Tree {
	trees := p.Parse(tokens)
	if len(trees) == 0 {
		return nil
	}
	return &trees[0]
}

func splitRules(rules []Rule) (map[string][]Rule, []Rule) {
	lexical := make(map[string][]Rule)
	binary := make([]Rule, 0, len(rules))
	for _, rule := range rules {
		switch len(rule.RHS) {
		case 1:
			lexical[rule.RHS[0]] = append(lexical[rule.RHS[0]], rule)
		case 2:
			binary = append(binary, rule)
		}
	}
	return lexical, binary
}

func ensureCell(chart [][]map[Symbol][]Tree, i, j int) map[Symbol][]Tree {
	if chart[i][j] == nil {
		chart[i][j] = make(map[Symbol][]Tree)
	}
	return chart[i][j]
}
