package tree

import (
	"strings"
	"testing"
)

func TestTree(t *testing.T) {
	tr, _ := FromString("(S (NP I) (VP saw))")
	if tr.Label != "S" || len(tr.Leaves()) != 2 {
		t.Fatal(tr.String())
	}
}

func TestNew(t *testing.T) {
	tr := New("S", "I", "saw")
	if tr.Label != "S" {
		t.Errorf("expected label 'S', got %s", tr.Label)
	}
	if len(tr.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(tr.Children))
	}
}

func TestNewNoChildren(t *testing.T) {
	tr := New("S")
	if tr.Label != "S" {
		t.Errorf("expected label 'S', got %s", tr.Label)
	}
	if len(tr.Children) != 0 {
		t.Errorf("expected 0 children, got %d", len(tr.Children))
	}
}

func TestStringLeaf(t *testing.T) {
	tr := New("word")
	got := tr.String()
	if got != "(word)" {
		t.Errorf("expected '(word)', got %s", got)
	}
}

func TestStringInternal(t *testing.T) {
	tr := New("S", New("NP", "I"), "saw")
	got := tr.String()
	if !strings.HasPrefix(got, "(S") {
		t.Errorf("expected to start with '(S', got %s", got)
	}
	if !strings.Contains(got, "NP") {
		t.Errorf("expected to contain 'NP', got %s", got)
	}
}

func TestLeaves(t *testing.T) {
	tr := New("S", New("NP", "I"), New("VP", "saw"))
	leaves := tr.Leaves()
	if len(leaves) != 2 {
		t.Errorf("expected 2 leaves, got %d", len(leaves))
	}
}

func TestLeavesDeepTree(t *testing.T) {
	tr := New("S",
		New("NP", New("DET", "the"), New("NN", "dog")),
		New("VP", New("V", "runs")))
	leaves := tr.Leaves()
	if len(leaves) != 3 {
		t.Errorf("expected 3 leaves, got %d", len(leaves))
	}
}

func TestLeavesEmptyTree(t *testing.T) {
	tr := New("S")
	leaves := tr.Leaves()
	if len(leaves) != 0 {
		t.Errorf("expected 0 leaves, got %d", len(leaves))
	}
}

func TestFromStringSimple(t *testing.T) {
	tr, err := FromString("(S (NP I) (VP saw))")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Label != "S" {
		t.Errorf("expected label 'S', got %s", tr.Label)
	}
}

func TestFromStringEmpty(t *testing.T) {
	tr, err := FromString("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Label != "" {
		t.Errorf("expected empty label, got %s", tr.Label)
	}
}

func TestFromStringLeaf(t *testing.T) {
	tr, err := FromString("(word)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Label != "word" {
		t.Errorf("expected label 'word', got %s", tr.Label)
	}
}

func TestHeight(t *testing.T) {
	tr := New("S",
		New("NP",
			New("DET", "the"),
			New("NN", "cat")))
	// Height of (S (NP (DET the) (NN cat))) = 3
	if tr.Height() != 3 {
		t.Errorf("expected height 3, got %d", tr.Height())
	}
}

func TestHeightLeaf(t *testing.T) {
	tr := New("S")
	if tr.Height() != 1 {
		t.Errorf("expected height 1 for leaf, got %d", tr.Height())
	}
}

func TestHeightFlat(t *testing.T) {
	tr := New("S", "a", "b", "c")
	if tr.Height() != 1 {
		t.Errorf("expected height 1 for flat, got %d", tr.Height())
	}
}

func TestMap(t *testing.T) {
	tr := New("S", New("NP", "I"), "saw")
	mapped := tr.Map(func(s string) string {
		return strings.ToUpper(s)
	})
	if mapped.Label != "S" {
		t.Errorf("expected uppercased label 'S', got %s", mapped.Label)
	}
}

func TestMapPreservesStructure(t *testing.T) {
	tr := New("S", New("NP", "I"), "saw")
	mapped := tr.Map(func(s string) string { return s })
	if len(mapped.Children) != len(tr.Children) {
		t.Errorf("expected same number of children, got %d vs %d",
			len(mapped.Children), len(tr.Children))
	}
}

func TestMapNilChild(t *testing.T) {
	tr := New("S", New("NP", 42), "word")
	mapped := tr.Map(func(s string) string { return s })
	_ = mapped
}

func TestFromStringComplex(t *testing.T) {
	tr, err := FromString("(S (NP (DET the) (NN dog)) (VP (V saw)))")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Label != "S" {
		t.Errorf("expected label 'S', got %s", tr.Label)
	}
	if len(tr.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(tr.Children))
	}
}

func TestFromStringNested(t *testing.T) {
	tr, err := FromString("(A (B (C (D leaf))))")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Label != "A" {
		t.Errorf("expected label 'A', got %s", tr.Label)
	}
	if tr.Height() != 4 {
		t.Errorf("expected height 4, got %d", tr.Height())
	}
}

func TestFromStringMalformed(t *testing.T) {
	// Test that malformed strings don't crash
	_, err := FromString("(((broken")
	if err != nil {
		// It's OK if it returns an error
		_ = err
	}
}

func TestImmutableTree(t *testing.T) {
	// Just test that the type exists
	tr := New("S", "I")
	_ = ImmutableTree{Tree: tr}
}

func TestStringMultiLine(t *testing.T) {
	tr := New("S", New("NP", "I"), New("VP", "saw"))
	got := tr.String()
	if !strings.HasSuffix(got, ")") {
		t.Errorf("expected string to end with ')', got %s", got)
	}
}

func TestFromStringWhitespace(t *testing.T) {
	tr, err := FromString("  (  S  ( NP I )  ( VP saw )  )  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Label != "S" {
		t.Errorf("expected label 'S', got %s", tr.Label)
	}
}

func TestLeavesOrder(t *testing.T) {
	tr := New("S", "a", "b", "c")
	leaves := tr.Leaves()
	expected := []any{"a", "b", "c"}
	for i, e := range expected {
		if leaves[i] != e {
			t.Errorf("leaf %d: expected %v, got %v", i, e, leaves[i])
		}
	}
}

func BenchmarkFromString(b *testing.B) {
	s := "(S (NP (DT the) (JJ big) (NN dog)) (VP (VBD saw) (NP (DT a) (NN cat))))"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FromString(s)
	}
}

func BenchmarkLeaves(b *testing.B) {
	tr := New("S",
		New("NP", New("DT", "the"), New("NN", "dog")),
		New("VP", New("VBD", "saw")))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Leaves()
	}
}
