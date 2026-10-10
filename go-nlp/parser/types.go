package parser

import "fmt"

// Symbol identifies terminals and non-terminals.
type Symbol string

// Tree is a lightweight parse tree node.
type Tree struct {
	Label    Symbol
	Children []Tree
	Word     string
}

// String returns a bracketed tree representation.
func (t Tree) String() string {
	if len(t.Children) == 0 {
		if t.Word != "" {
			return fmt.Sprintf("(%s %s)", t.Label, t.Word)
		}
		return fmt.Sprintf("(%s)", t.Label)
	}
	out := "(" + string(t.Label)
	for _, c := range t.Children {
		out += " " + c.String()
	}
	return out + ")"
}

// Rule represents either a lexical rule (A -> word) or a binary rule (A -> B C).
type Rule struct {
	LHS    Symbol
	RHS    []string
	Weight float64
}

// Grammar stores a compact CFG-like grammar.
type Grammar struct {
	Start Symbol
	Rules []Rule
}

// NewGrammar creates a grammar with the provided start symbol and rules.
func NewGrammar(start Symbol, rules []Rule) *Grammar {
	return &Grammar{Start: start, Rules: rules}
}

// Parser defines the parsing interface inspired by nltk.parse.api.ParserI.
type Parser interface {
	Parse(tokens []string) []Tree
}
