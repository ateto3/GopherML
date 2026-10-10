package draw

import (
	"strings"
	"testing"

	"github.com/ateto3/gopherml/go-nlp/tree"
)

func TestDrawBasic(t *testing.T) {
	tr := tree.New("S", tree.New("NP", "I"), tree.New("VP", "saw"))
	if ASCII(tr) == "" || DOT(tr) == "" {
		t.Fatal("draw")
	}
}

func TestASCIIContainsAllLabels(t *testing.T) {
	tr := tree.New("S", tree.New("NP", "I"), tree.New("VP", "saw"))
	out := ASCII(tr)
	for _, label := range []string{"S", "NP", "VP", "I", "saw"} {
		if !strings.Contains(out, label) {
			t.Errorf("expected ASCII output to contain %q, got:\n%s", label, out)
		}
	}
}

func TestDOTFormat(t *testing.T) {
	tr := tree.New("S", tree.New("NP", "I"), tree.New("VP", "saw"))
	out := DOT(tr)
	if !strings.HasPrefix(out, "digraph G {") {
		t.Errorf("expected DOT output to start with 'digraph G {', got: %s", out)
	}
	if !strings.HasSuffix(strings.TrimRight(out, "\n"), "}") {
		t.Error("expected DOT output to end with '}'")
	}
}

func TestDOTContainsAllLabels(t *testing.T) {
	tr := tree.New("S", tree.New("NP", "I"), tree.New("VP", "saw"))
	out := DOT(tr)
	for _, label := range []string{"S", "NP", "VP", "I", "saw"} {
		if !strings.Contains(out, label) {
			t.Errorf("expected DOT output to contain label %q, got:\n%s", label, out)
		}
	}
}

func TestASCIILeafNodes(t *testing.T) {
	tr := tree.New("S", "I", "saw")
	out := ASCII(tr)
	if !strings.Contains(out, "S") || !strings.Contains(out, "I") || !strings.Contains(out, "saw") {
		t.Errorf("expected all leaves, got:\n%s", out)
	}
}

func TestDOTLeafNodes(t *testing.T) {
	tr := tree.New("S", "I", "saw")
	out := DOT(tr)
	if !strings.Contains(out, "digraph G") {
		t.Errorf("expected digraph, got: %s", out)
	}
}

func TestASCIIEmptyChildren(t *testing.T) {
	tr := tree.New("S")
	out := ASCII(tr)
	if !strings.Contains(out, "S") {
		t.Errorf("expected 'S' in output, got: %s", out)
	}
}

func TestDOTEmptyChildren(t *testing.T) {
	tr := tree.New("S")
	out := DOT(tr)
	if !strings.Contains(out, "S") {
		t.Errorf("expected 'S' in DOT output, got: %s", out)
	}
}

func TestASCIIMultipleLevels(t *testing.T) {
	tr := tree.New("A",
		tree.New("B",
			tree.New("C", "leaf1"),
			"leaf2"),
		tree.New("D", "leaf3"))
	out := ASCII(tr)
	for _, l := range []string{"A", "B", "C", "D", "leaf1", "leaf2", "leaf3"} {
		if !strings.Contains(out, l) {
			t.Errorf("expected %q in output, got:\n%s", l, out)
		}
	}
}

func TestDOTHasArrows(t *testing.T) {
	tr := tree.New("S", tree.New("NP", "I"), tree.New("VP", "saw"))
	out := DOT(tr)
	if !strings.Contains(out, "->") {
		t.Errorf("expected '->' in DOT output, got: %s", out)
	}
}

func TestDOTNodeCount(t *testing.T) {
	// 1 root + 2 children + 2 leaves = 5 nodes
	tr := tree.New("S", tree.New("NP", "I"), tree.New("VP", "saw"))
	out := DOT(tr)
	// Count "n%d" patterns (excluding the digraph declaration)
	lines := strings.Split(out, "\n")
	nodeCount := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "n") && strings.Contains(line, "[label=") {
			nodeCount++
		}
	}
	if nodeCount != 5 {
		t.Errorf("expected 5 nodes in DOT output, got %d", nodeCount)
	}
}

func BenchmarkASCII(b *testing.B) {
	tr := tree.New("S",
		tree.New("NP", tree.New("DT", "the"), tree.New("NN", "cat")),
		tree.New("VP", tree.New("VBZ", "runs")))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ASCII(tr)
	}
}

func BenchmarkDOT(b *testing.B) {
	tr := tree.New("S",
		tree.New("NP", tree.New("DT", "the"), tree.New("NN", "cat")),
		tree.New("VP", tree.New("VBZ", "runs")))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DOT(tr)
	}
}
