package grammar

import "strings"

type Nonterminal string
type Production struct{ LHS Nonterminal; RHS []string }
type CFG struct{ Start Nonterminal; Productions []Production }

func NT(s string) Nonterminal { return Nonterminal(s) }
func (n Nonterminal) String() string { return string(n) }
func (p Production) String() string { return string(p.LHS)+" -> "+strings.Join(p.RHS," ") }
