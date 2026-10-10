package inference

import (
	"strings"
	"testing"

	"github.com/ateto3/gopherml/go-nlp/sem"
)

func TestProver(t *testing.T) {
	c := NewCommand(sem.Atom("q"), []sem.Expr{sem.Atom("p"), sem.Imp{A: sem.Atom("p"), B: sem.Atom("q")}})
	if !c.Prove() {
		t.Fatal("expected proof")
	}
}

func TestProverStoresResult(t *testing.T) {
	c := NewCommand(sem.Atom("q"), []sem.Expr{sem.Imp{A: sem.Atom("p"), B: sem.Atom("q")}, sem.Atom("p")})
	ok := c.Prove()
	if !ok {
		t.Error("expected prove to succeed")
	}
	if !c.Result {
		t.Error("expected Result field to be true")
	}
	if c.Proof == "" {
		t.Error("expected non-empty Proof")
	}
}

func TestProverFalseGoal(t *testing.T) {
	c := NewCommand(sem.Atom("r"), []sem.Expr{sem.Atom("p")})
	if c.Prove() {
		t.Fatal("expected proof to fail when goal isn't a known fact")
	}
}

func TestProverEmptyAssumptions(t *testing.T) {
	c := NewCommand(sem.Atom("p"), nil)
	if c.Prove() {
		t.Fatal("expected proof to fail with no assumptions and unknown goal")
	}
}

func TestProverModusPonensDirect(t *testing.T) {
	// From p and (p -> q), conclude q
	c := NewCommand(
		sem.Atom("q"),
		[]sem.Expr{sem.Atom("p"), sem.Imp{A: sem.Atom("p"), B: sem.Atom("q")}},
	)
	if !c.Prove() {
		t.Error("expected to prove q from p and (p -> q)")
	}
}

func TestProverConjunction(t *testing.T) {
	// From p and q, conclude (p and q)
	c := NewCommand(
		sem.And{A: sem.Atom("p"), B: sem.Atom("q")},
		[]sem.Expr{sem.Atom("p"), sem.Atom("q")},
	)
	if !c.Prove() {
		t.Error("expected to prove (p and q) from p and q")
	}
}

func TestProverConjunctionFailsPartial(t *testing.T) {
	// Cannot prove (p and q) from p alone
	c := NewCommand(
		sem.And{A: sem.Atom("p"), B: sem.Atom("q")},
		[]sem.Expr{sem.Atom("p")},
	)
	if c.Prove() {
		t.Error("expected proof to fail when only one of p and q is given")
	}
}

func TestProverChainedImplication(t *testing.T) {
	// From (p -> q) and (q -> r), conclude r
	c := NewCommand(
		sem.Atom("r"),
		[]sem.Expr{
			sem.Imp{A: sem.Atom("p"), B: sem.Atom("q")},
			sem.Imp{A: sem.Atom("q"), B: sem.Atom("r")},
			sem.Atom("p"),
		},
	)
	if !c.Prove() {
		t.Error("expected to prove r via chained implications")
	}
}

func TestProverImplicationGoal(t *testing.T) {
	// Prove (p -> q) given that p implies q via a fact
	// Note: this prover only does forward chaining, so it may not handle this case
	c := NewCommand(
		sem.Imp{A: sem.Atom("p"), B: sem.Atom("q")},
		[]sem.Expr{sem.Atom("p"), sem.Imp{A: sem.Atom("p"), B: sem.Atom("q")}},
	)
	// Goal is an implication — sat() returns false for non-atoms, so the prover should fail
	if c.Prove() {
		t.Error("expected proof to fail for implication goal in forward chaining")
	}
}

func TestResolutionProverDirect(t *testing.T) {
	prover := NewResolutionProver()
	ok, proof := prover.Prove(
		sem.Atom("q"),
		[]sem.Expr{sem.Atom("p"), sem.Imp{A: sem.Atom("p"), B: sem.Atom("q")}},
	)
	if !ok {
		t.Error("expected direct prove to succeed")
	}
	if !strings.Contains(proof, "prov") {
		t.Errorf("expected proof string, got: %s", proof)
	}
}

func TestResolutionProverFailure(t *testing.T) {
	prover := NewResolutionProver()
	ok, proof := prover.Prove(sem.Atom("r"), []sem.Expr{sem.Atom("p")})
	if ok {
		t.Error("expected prove to fail")
	}
	if proof != "unproved" {
		t.Errorf("expected 'unproved' proof, got: %s", proof)
	}
}

func TestNewCommand(t *testing.T) {
	c := NewCommand(sem.Atom("q"), []sem.Expr{sem.Atom("p")})
	if c == nil {
		t.Fatal("NewCommand returned nil")
	}
	if c.Goal == nil {
		t.Error("expected Goal to be set")
	}
	if len(c.Assumptions) != 1 {
		t.Errorf("expected 1 assumption, got %d", len(c.Assumptions))
	}
}

func TestSatAtom(t *testing.T) {
	facts := map[string]struct{}{"p": {}, "q": {}}
	if !sat(sem.Atom("p"), facts) {
		t.Error("expected 'p' to be satisfied")
	}
	if sat(sem.Atom("r"), facts) {
		t.Error("expected 'r' to not be satisfied")
	}
}

func TestSatAnd(t *testing.T) {
	facts := map[string]struct{}{"p": {}, "q": {}}
	if !sat(sem.And{A: sem.Atom("p"), B: sem.Atom("q")}, facts) {
		t.Error("expected (p and q) to be satisfied")
	}
	if sat(sem.And{A: sem.Atom("p"), B: sem.Atom("r")}, facts) {
		t.Error("expected (p and r) to not be satisfied")
	}
}

func TestSatOr(t *testing.T) {
	facts := map[string]struct{}{"p": {}}
	if !sat(sem.Or{A: sem.Atom("p"), B: sem.Atom("q")}, facts) {
		t.Error("expected (p or q) to be satisfied")
	}
}

func TestSatNot(t *testing.T) {
	facts := map[string]struct{}{"p": {}}
	if !sat(sem.Not{X: sem.Atom("r")}, facts) {
		t.Error("expected not(r) to be satisfied")
	}
	if sat(sem.Not{X: sem.Atom("p")}, facts) {
		t.Error("expected not(p) to not be satisfied")
	}
}

func TestProverInterface(t *testing.T) {
	var _ Prover = NewResolutionProver()
}

func BenchmarkProve(b *testing.B) {
	c := NewCommand(
		sem.Atom("r"),
		[]sem.Expr{
			sem.Imp{A: sem.Atom("p"), B: sem.Atom("q")},
			sem.Imp{A: sem.Atom("q"), B: sem.Atom("r")},
			sem.Atom("p"),
		},
	)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Prove()
	}
}
