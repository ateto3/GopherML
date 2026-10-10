package sem

import "strings"

type Atom string
type Not struct{ X Expr }
type And struct{ A, B Expr }
type Or struct{ A, B Expr }
type Imp struct{ A, B Expr }
type All struct{ Var, Pred string; Body Expr }
type Exists struct{ Var, Pred string; Body Expr }

func (a Atom) Eval(m Model, env map[string]string) bool {
	if v, ok := m.Val[string(a)]; ok {
		switch x := v.(type) {
		case bool:
			return x
		case map[string]struct{}:
			return env != nil && x[env["x"]] == struct{}{}
		}
	}
	return false
}
func (n Not) Eval(m Model, e map[string]string) bool { return !n.X.Eval(m, e) }
func (a And) Eval(m Model, e map[string]string) bool { return a.A.Eval(m, e) && a.B.Eval(m, e) }
func (o Or) Eval(m Model, e map[string]string) bool  { return o.A.Eval(m, e) || o.B.Eval(m, e) }
func (i Imp) Eval(m Model, e map[string]string) bool { return !i.A.Eval(m, e) || i.B.Eval(m, e) }

func (a All) Eval(m Model, e map[string]string) bool {
	for d := range m.Dom {
		if e == nil { e = map[string]string{} }
		e[a.Var] = d
		if !a.Body.Eval(m, e) { return false }
	}
	return true
}
func (x Exists) Eval(m Model, e map[string]string) bool {
	for d := range m.Dom {
		if e == nil { e = map[string]string{} }
		e[x.Var] = d
		if x.Body.Eval(m, e) { return true }
	}
	return false
}

func ParseExpr(s string) Expr {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "not ") { return Not{ParseExpr(strings.TrimSpace(s[4:]))} }
	if i := strings.Index(s, " and "); i > 0 { return And{ParseExpr(s[:i]), ParseExpr(s[i+5:])} }
	if i := strings.Index(s, " or "); i > 0 { return Or{ParseExpr(s[:i]), ParseExpr(s[i+4:])} }
	if i := strings.Index(s, " -> "); i > 0 { return Imp{ParseExpr(s[:i]), ParseExpr(s[i+4:])} }
	if strings.HasPrefix(s, "all ") { return parseQuant(s, true) }
	if strings.HasPrefix(s, "exists ") { return parseQuant(s, false) }
	return Atom(s)
}

func parseQuant(s string, all bool) Expr {
	parts := strings.Fields(s)
	if len(parts) < 4 { return Atom(s) }
	varName := parts[1]
	body := strings.Join(parts[3:], " ")
	if all { return All{Var: varName, Body: ParseExpr(body)} }
	return Exists{Var: varName, Body: ParseExpr(body)}
}
