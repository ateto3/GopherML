package parser

import (
	"testing"
)

func TestChartParser(t *testing.T) {
	grammar := NewGrammar("S", []Rule{
		{LHS: "S", RHS: []string{"NP", "VP"}},
		{LHS: "NP", RHS: []string{"Det", "N"}},
		{LHS: "VP", RHS: []string{"V", "NP"}},
		{LHS: "Det", RHS: []string{"the"}},
		{LHS: "N", RHS: []string{"dog"}},
		{LHS: "N", RHS: []string{"cat"}},
		{LHS: "V", RHS: []string{"sees"}},
	})

	p := NewChartParser(grammar)
	trees := p.Parse([]string{"the", "dog", "sees", "the", "cat"})

	if len(trees) == 0 {
		t.Fatal("expected at least one parse tree")
	}
	if trees[0].Label != "S" {
		t.Fatalf("expected root S, got %s", trees[0].Label)
	}
}

func TestChartParserSimpleSentence(t *testing.T) {
	grammar := NewGrammar("S", []Rule{
		{LHS: "S", RHS: []string{"NP", "VP"}},
		{LHS: "NP", RHS: []string{"I"}},
		{LHS: "VP", RHS: []string{"saw"}},
	})
	p := NewChartParser(grammar)
	trees := p.Parse([]string{"I", "saw"})
	if len(trees) == 0 {
		t.Fatal("expected at least one parse tree")
	}
	if trees[0].Label != "S" {
		t.Errorf("expected root S, got %s", trees[0].Label)
	}
}

func TestChartParserFailsOnUnknownWords(t *testing.T) {
	grammar := NewGrammar("S", []Rule{
		{LHS: "S", RHS: []string{"NP", "VP"}},
		{LHS: "NP", RHS: []string{"I"}},
		{LHS: "VP", RHS: []string{"saw"}},
	})
	p := NewChartParser(grammar)
	trees := p.Parse([]string{"unknown", "word"})
	if len(trees) != 0 {
		t.Errorf("expected no parse for unknown words, got %d", len(trees))
	}
}

func TestChartParserEmpty(t *testing.T) {
	grammar := NewGrammar("S", []Rule{
		{LHS: "S", RHS: []string{"NP", "VP"}},
	})
	p := NewChartParser(grammar)
	trees := p.Parse([]string{})
	if len(trees) != 0 {
		t.Errorf("expected no parse for empty input, got %d", len(trees))
	}
}

func TestNewGrammar(t *testing.T) {
	g := NewGrammar("S", []Rule{
		{LHS: "S", RHS: []string{"NP", "VP"}},
		{LHS: "NP", RHS: []string{"Det", "N"}},
	})
	if g.Start != "S" {
		t.Errorf("expected start 'S', got %s", g.Start)
	}
	if len(g.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(g.Rules))
	}
}

func TestChartParserTreeStructure(t *testing.T) {
	grammar := NewGrammar("S", []Rule{
		{LHS: "S", RHS: []string{"NP", "VP"}},
		{LHS: "NP", RHS: []string{"Det", "N"}},
		{LHS: "VP", RHS: []string{"V", "NP"}},
		{LHS: "Det", RHS: []string{"the"}},
		{LHS: "N", RHS: []string{"dog"}},
		{LHS: "V", RHS: []string{"sees"}},
	})
	p := NewChartParser(grammar)
	trees := p.Parse([]string{"the", "dog", "sees", "the", "dog"})
	if len(trees) == 0 {
		t.Fatal("expected at least one parse tree")
	}
	if len(trees[0].Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(trees[0].Children))
	}
}

func TestChartParserSingleWord(t *testing.T) {
	grammar := NewGrammar("S", []Rule{
		{LHS: "S", RHS: []string{"hello"}},
	})
	p := NewChartParser(grammar)
	trees := p.Parse([]string{"hello"})
	if len(trees) == 0 {
		t.Fatal("expected at least one parse tree for single word")
	}
}

func TestChartParserMultipleParses(t *testing.T) {
	// S -> NP VP | VP NP (two ways to parse)
	grammar := NewGrammar("S", []Rule{
		{LHS: "S", RHS: []string{"NP", "VP"}},
		{LHS: "S", RHS: []string{"VP", "NP"}},
		{LHS: "NP", RHS: []string{"I"}},
		{LHS: "VP", RHS: []string{"ran"}},
	})
	p := NewChartParser(grammar)
	trees := p.Parse([]string{"I", "ran"})
	if len(trees) == 0 {
		t.Fatal("expected at least one parse tree")
	}
	for _, tree := range trees {
		if tree.Label != "S" {
			t.Errorf("expected root S, got %s", tree.Label)
		}
	}
}

func TestChartParserLeftAssoc(t *testing.T) {
	// NP -> NP PP (left associative binary)
	grammar := NewGrammar("NP", []Rule{
		{LHS: "NP", RHS: []string{"NP", "PP"}},
		{LHS: "NP", RHS: []string{"I"}},
		{LHS: "PP", RHS: []string{"home"}},
	})
	p := NewChartParser(grammar)
	trees := p.Parse([]string{"I", "home"})
	if len(trees) == 0 {
		t.Fatal("expected at least one parse tree")
	}
}

func TestRuleString(t *testing.T) {
	r := Rule{LHS: "S", RHS: []string{"NP", "VP"}}
	if r.LHS != "S" {
		t.Errorf("expected LHS 'S', got %s", r.LHS)
	}
	if len(r.RHS) != 2 {
		t.Errorf("expected 2 RHS, got %d", len(r.RHS))
	}
}

func TestSymbolString(t *testing.T) {
	s := Symbol("VP")
	if string(s) != "VP" {
		t.Errorf("expected 'VP', got %s", string(s))
	}
}

func TestTreeString(t *testing.T) {
	tree := Tree{
		Label:    "S",
		Children: []Tree{{Label: "NP", Word: "I"}, {Label: "VP", Word: "ran"}},
	}
	got := tree.String()
	if got == "" {
		t.Error("expected non-empty tree string")
	}
	if got[0] != '(' {
		t.Errorf("expected tree string to start with '(', got: %s", got)
	}
}

func TestParserInterface(t *testing.T) {
	var p Parser = NewChartParser(NewGrammar("S", []Rule{
		{LHS: "S", RHS: []string{"hello"}},
	}))
	trees := p.Parse([]string{"hello"})
	if len(trees) == 0 {
		t.Error("expected at least one parse via interface")
	}
}

func TestChartParserLexicalOnly(t *testing.T) {
	grammar := NewGrammar("Word", []Rule{
		{LHS: "Word", RHS: []string{"hello"}},
	})
	p := NewChartParser(grammar)
	trees := p.Parse([]string{"hello"})
	if len(trees) == 0 {
		t.Fatal("expected parse for lexical rule")
	}
	if trees[0].Label != "Word" {
		t.Errorf("expected root 'Word', got %s", trees[0].Label)
	}
}

func BenchmarkChartParser(b *testing.B) {
	grammar := NewGrammar("S", []Rule{
		{LHS: "S", RHS: []string{"NP", "VP"}},
		{LHS: "NP", RHS: []string{"Det", "N"}},
		{LHS: "VP", RHS: []string{"V", "NP"}},
		{LHS: "Det", RHS: []string{"the"}},
		{LHS: "N", RHS: []string{"dog"}},
		{LHS: "V", RHS: []string{"sees"}},
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := NewChartParser(grammar)
		p.Parse([]string{"the", "dog", "sees", "the", "dog"})
	}
}
