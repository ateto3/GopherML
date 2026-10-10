package sem

import (
	"testing"
)

func TestSem(t *testing.T) {
	m := Model{Dom: map[string]struct{}{"a": {}}, Val: map[string]any{"p": true}}
	if !Evaluate(ParseExpr("p"), m) {
		t.Fatal("atom")
	}
	if Evaluate(ParseExpr("not p"), m) {
		t.Fatal("not")
	}
}

func TestAtomEvalTrue(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": true}}
	if !Evaluate(Atom("p"), m) {
		t.Error("expected p to be true")
	}
}

func TestAtomEvalFalse(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": false}}
	if Evaluate(Atom("p"), m) {
		t.Error("expected p to be false")
	}
}

func TestAtomEvalUnknown(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{}}
	if Evaluate(Atom("unknown"), m) {
		t.Error("expected unknown atom to be false")
	}
}

func TestNotTrue(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": false}}
	if !Evaluate(Not{X: Atom("p")}, m) {
		t.Error("expected not(p) to be true when p is false")
	}
}

func TestNotFalse(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": true}}
	if Evaluate(Not{X: Atom("p")}, m) {
		t.Error("expected not(p) to be false when p is true")
	}
}

func TestAndBothTrue(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": true, "q": true}}
	if !Evaluate(And{A: Atom("p"), B: Atom("q")}, m) {
		t.Error("expected p and q to be true")
	}
}

func TestAndOneFalse(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": true, "q": false}}
	if Evaluate(And{A: Atom("p"), B: Atom("q")}, m) {
		t.Error("expected p and q to be false")
	}
}

func TestAndBothFalse(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": false, "q": false}}
	if Evaluate(And{A: Atom("p"), B: Atom("q")}, m) {
		t.Error("expected p and q to be false")
	}
}

func TestOrFirstTrue(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": true, "q": false}}
	if !Evaluate(Or{A: Atom("p"), B: Atom("q")}, m) {
		t.Error("expected p or q to be true")
	}
}

func TestOrSecondTrue(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": false, "q": true}}
	if !Evaluate(Or{A: Atom("p"), B: Atom("q")}, m) {
		t.Error("expected p or q to be true")
	}
}

func TestOrBothFalse(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": false, "q": false}}
	if Evaluate(Or{A: Atom("p"), B: Atom("q")}, m) {
		t.Error("expected p or q to be false")
	}
}

func TestImplicationTrue(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": false, "q": false}}
	if !Evaluate(Imp{A: Atom("p"), B: Atom("q")}, m) {
		t.Error("expected p->q to be true when p is false")
	}
}

func TestImplicationFalse(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": true, "q": false}}
	if Evaluate(Imp{A: Atom("p"), B: Atom("q")}, m) {
		t.Error("expected p->q to be false when p=true, q=false")
	}
}

func TestParseExprAtom(t *testing.T) {
	expr := ParseExpr("p")
	if _, ok := expr.(Atom); !ok {
		t.Errorf("expected Atom, got %T", expr)
	}
}

func TestParseExprNot(t *testing.T) {
	expr := ParseExpr("not p")
	if _, ok := expr.(Not); !ok {
		t.Errorf("expected Not, got %T", expr)
	}
}

func TestParseExprAnd(t *testing.T) {
	expr := ParseExpr("p and q")
	if _, ok := expr.(And); !ok {
		t.Errorf("expected And, got %T", expr)
	}
}

func TestParseExprOr(t *testing.T) {
	expr := ParseExpr("p or q")
	if _, ok := expr.(Or); !ok {
		t.Errorf("expected Or, got %T", expr)
	}
}

func TestParseExprImplication(t *testing.T) {
	expr := ParseExpr("p -> q")
	if _, ok := expr.(Imp); !ok {
		t.Errorf("expected Imp, got %T", expr)
	}
}

func TestEvaluateNotAtom(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": true}}
	if Evaluate(Not{X: Atom("p")}, m) {
		t.Error("expected not(p) to be false (since p is true)")
	}
}

func TestEvaluateAnd(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": true, "q": true}}
	if !Evaluate(And{A: Atom("p"), B: Atom("q")}, m) {
		t.Error("expected p and q to be true")
	}
}

func TestEvaluateOr(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": false, "q": true}}
	if !Evaluate(Or{A: Atom("p"), B: Atom("q")}, m) {
		t.Error("expected p or q to be true")
	}
}

func TestEvaluateImplication(t *testing.T) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": true, "q": true}}
	if !Evaluate(Imp{A: Atom("p"), B: Atom("q")}, m) {
		t.Error("expected p->q to be true")
	}
}

func TestExprInterface(t *testing.T) {
	var _ Expr = Atom("p")
	var _ Expr = Not{X: Atom("p")}
	var _ Expr = And{A: Atom("p"), B: Atom("q")}
	var _ Expr = Or{A: Atom("p"), B: Atom("q")}
	var _ Expr = Imp{A: Atom("p"), B: Atom("q")}
}

func BenchmarkEvaluateSimple(b *testing.B) {
	m := Model{Dom: map[string]struct{}{}, Val: map[string]any{"p": true, "q": true}}
	expr := And{A: Atom("p"), B: Atom("q")}
	for i := 0; i < b.N; i++ {
		Evaluate(expr, m)
	}
}
