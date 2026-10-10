package inference

import "github.com/ateto3/gopherml/go-nlp/sem"

type ResolutionProver struct{}

func NewResolutionProver() ResolutionProver { return ResolutionProver{} }

func (ResolutionProver) Prove(goal sem.Expr, assumptions []sem.Expr) (bool, string) {
	facts := map[string]struct{}{}
	rules := []sem.Imp{}
	for _, a := range assumptions {
		collect(a, facts, &rules)
	}
	changed := true
	for changed {
		changed = false
		for _, r := range rules {
			if sat(r.A, facts) && addFact(r.B, facts) {
				changed = true
			}
		}
	}
	ok := sat(goal, facts)
	if ok {
		return true, "proved"
	}
	return false, "unproved"
}

func collect(e sem.Expr, facts map[string]struct{}, rules *[]sem.Imp) {
	switch x := e.(type) {
	case sem.Atom:
		facts[string(x)] = struct{}{}
	case sem.And:
		collect(x.A, facts, rules)
		collect(x.B, facts, rules)
	case sem.Imp:
		*rules = append(*rules, x)
	}
}

func sat(e sem.Expr, facts map[string]struct{}) bool {
	switch x := e.(type) {
	case sem.Atom:
		_, ok := facts[string(x)]
		return ok
	case sem.And:
		return sat(x.A, facts) && sat(x.B, facts)
	case sem.Or:
		return sat(x.A, facts) || sat(x.B, facts)
	case sem.Not:
		return !sat(x.X, facts)
	default:
		return false
	}
}

func addFact(e sem.Expr, facts map[string]struct{}) bool {
	if a, ok := e.(sem.Atom); ok {
		k := string(a)
		if _, seen := facts[k]; !seen {
			facts[k] = struct{}{}
			return true
		}
	}
	return false
}
