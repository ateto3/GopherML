package grammar

import "strings"

func NewCFG(start Nonterminal, prods []Production) CFG { return CFG{Start: start, Productions: prods} }

func ParseCFG(s string) CFG {
	lines:=strings.Split(s,"\n")
	var start Nonterminal
	prods:=[]Production{}
	for _,line:=range lines {
		line=strings.TrimSpace(line)
		if line=="" { continue }
		parts:=strings.Split(line,"->")
		if len(parts)!=2 { continue }
		lhs:=NT(strings.TrimSpace(parts[0]))
		if start=="" { start=lhs }
		for _,alt:=range strings.Split(parts[1],"|") {
			rhs:=strings.Fields(strings.TrimSpace(alt))
			prods=append(prods, Production{LHS: lhs, RHS: rhs})
		}
	}
	return NewCFG(start, prods)
}
