package grammar

import (
	"strings"
	"testing"
)

func TestCFG(t *testing.T) {
	g := ParseCFG("S -> NP VP\nNP -> I\nVP -> saw")
	if g.Start != "S" || len(g.Productions) != 3 {
		t.Fatal(g)
	}
}

func TestCFGMultipleProductions(t *testing.T) {
	g := ParseCFG("S -> NP VP\nS -> Aux NP VP\nNP -> Det N")
	if len(g.Productions) != 3 {
		t.Errorf("expected 3 productions, got %d", len(g.Productions))
	}
}

func TestCFGAlternatives(t *testing.T) {
	g := ParseCFG("NP -> Det N | NP PP")
	if len(g.Productions) != 2 {
		t.Errorf("expected 2 productions (alternatives), got %d", len(g.Productions))
	}
}

func TestCFGEmptyLines(t *testing.T) {
	g := ParseCFG("\nS -> NP VP\n\nNP -> I\n")
	if g.Start != "S" {
		t.Errorf("expected S, got %s", g.Start)
	}
	if len(g.Productions) != 2 {
		t.Errorf("expected 2 productions, got %d", len(g.Productions))
	}
}

func TestCFGFirstNonterminalIsStart(t *testing.T) {
	g := ParseCFG("NP -> Det N\nVP -> V NP")
	if g.Start != "NP" {
		t.Errorf("expected first nonterminal 'NP' as start, got %s", g.Start)
	}
}

func TestNewCFG(t *testing.T) {
	g := NewCFG("S", []Production{
		{LHS: NT("S"), RHS: []string{"NP", "VP"}},
		{LHS: NT("NP"), RHS: []string{"Det", "N"}},
	})
	if g.Start != "S" {
		t.Errorf("expected start 'S', got %s", g.Start)
	}
	if len(g.Productions) != 2 {
		t.Errorf("expected 2 productions, got %d", len(g.Productions))
	}
}

func TestNT(t *testing.T) {
	if NT("S") != "S" {
		t.Errorf("expected NT to wrap string, got %s", NT("S"))
	}
}

func TestNonterminalString(t *testing.T) {
	n := Nonterminal("VP")
	if n.String() != "VP" {
		t.Errorf("expected 'VP', got %s", n.String())
	}
}

func TestProductionString(t *testing.T) {
	p := Production{LHS: NT("S"), RHS: []string{"NP", "VP"}}
	got := p.String()
	if !strings.Contains(got, "S") || !strings.Contains(got, "NP") || !strings.Contains(got, "VP") {
		t.Errorf("expected production string to contain S, NP, VP, got: %s", got)
	}
}

func TestParseCFGInvalidLines(t *testing.T) {
	// Lines without "->" should be silently skipped
	g := ParseCFG("invalid line here\nS -> NP VP")
	if g.Start != "S" {
		t.Errorf("expected S, got %s", g.Start)
	}
	if len(g.Productions) != 1 {
		t.Errorf("expected 1 production, got %d", len(g.Productions))
	}
}

func TestParseCFGMultipleRHSAlternatives(t *testing.T) {
	g := ParseCFG("VP -> V NP | V | VP PP")
	if len(g.Productions) != 3 {
		t.Errorf("expected 3 productions, got %d", len(g.Productions))
	}
}

func TestParseCFGWhitespaceHandling(t *testing.T) {
	g := ParseCFG("  S   ->  NP VP  \n  NP -> Det N  ")
	if g.Start != "S" {
		t.Errorf("expected S, got %s", g.Start)
	}
	if len(g.Productions) != 2 {
		t.Errorf("expected 2 productions, got %d", len(g.Productions))
	}
}

func TestParseCFGEmpty(t *testing.T) {
	g := ParseCFG("")
	if len(g.Productions) != 0 {
		t.Errorf("expected 0 productions for empty input, got %d", len(g.Productions))
	}
}

func TestParseCFGSingleProduction(t *testing.T) {
	g := ParseCFG("S -> 'hello'")
	if g.Start != "S" {
		t.Errorf("expected start S, got %s", g.Start)
	}
	if len(g.Productions) != 1 {
		t.Errorf("expected 1 production, got %d", len(g.Productions))
	}
	if g.Productions[0].LHS != "S" {
		t.Errorf("expected LHS 'S', got %s", g.Productions[0].LHS)
	}
	if len(g.Productions[0].RHS) != 1 {
		t.Errorf("expected 1-token RHS, got %d", len(g.Productions[0].RHS))
	}
}
