package translate

import (
	"strings"
	"testing"
)

func TestIBM1(t *testing.T) {
	c := []AlignedSent{{Words: []string{"the", "house"}, Mots: []string{"das", "haus"}}}
	m := NewIBMModel1(c, 2)
	if m == nil || len(c[0].Alignment) == 0 {
		t.Fatal("no align")
	}
}

func TestAlignmentInvert(t *testing.T) {
	a := Alignment{{0, 0}, {1, 1}, {2, 2}}
	inv := a.Invert()
	if len(inv) != len(a) {
		t.Errorf("expected same length after invert, got %d", len(inv))
	}
	if inv[0][0] != 0 || inv[0][1] != 0 {
		t.Errorf("invert failed: %v", inv)
	}
}

func TestAlignmentString(t *testing.T) {
	a := Alignment{{0, 1}, {1, 0}}
	got := a.String()
	if got == "" {
		t.Error("expected non-empty string")
	}
}

func TestAlignmentEmpty(t *testing.T) {
	a := Alignment{}
	if a.String() != "" {
		t.Errorf("expected empty string, got %s", a.String())
	}
}

func TestAlignedSentInvert(t *testing.T) {
	s := AlignedSent{Words: []string{"a", "b"}, Mots: []string{"x", "y"}, Alignment: Alignment{{0, 0}, {1, 1}}}
	s.Invert()
	if len(s.Words) != 2 || s.Words[0] != "x" {
		t.Errorf("expected Words[0]='x', got %s", s.Words[0])
	}
}

func TestNewIBMModel1Alignment(t *testing.T) {
	corpus := []AlignedSent{{Words: []string{"the", "house"}, Mots: []string{"das", "haus"}}}
	m := NewIBMModel1(corpus, 5)
	if m == nil {
		t.Fatal("NewIBMModel1 returned nil")
	}
	if len(corpus[0].Alignment) == 0 {
		t.Fatal("expected non-empty alignment")
	}
	if len(corpus[0].Alignment) != len(corpus[0].Words) {
		t.Errorf("alignment length %d != words length %d",
			len(corpus[0].Alignment), len(corpus[0].Words))
	}
}

func TestIBMModel1MultipleSentences(t *testing.T) {
	corpus := []AlignedSent{
		{Words: []string{"the", "house"}, Mots: []string{"das", "haus"}},
		{Words: []string{"a", "book"}, Mots: []string{"ein", "buch"}},
	}
	m := NewIBMModel1(corpus, 3)
	if m == nil {
		t.Fatal("NewIBMModel1 returned nil")
	}
	if len(corpus[0].Alignment) == 0 {
		t.Error("expected alignment for first sentence")
	}
	if len(corpus[1].Alignment) == 0 {
		t.Error("expected alignment for second sentence")
	}
}

func TestIBMModel1NoIterations(t *testing.T) {
	corpus := []AlignedSent{{Words: []string{"a"}, Mots: []string{"x"}}}
	m := NewIBMModel1(corpus, 0)
	if m == nil {
		t.Fatal("NewIBMModel1 with 0 iters returned nil")
	}
}

func TestIBMModel1EmptyCorpus(t *testing.T) {
	m := NewIBMModel1(nil, 2)
	if m == nil {
		t.Fatal("NewIBMModel1 with nil corpus returned nil")
	}
}

func TestIBMModel1TableInitialized(t *testing.T) {
	corpus := []AlignedSent{{Words: []string{"a", "b"}, Mots: []string{"x", "y"}}}
	m := NewIBMModel1(corpus, 1)
	if m.Table == nil {
		t.Error("expected non-nil Table")
	}
}

func TestIBMModel1AlignsValid(t *testing.T) {
	corpus := []AlignedSent{{Words: []string{"the", "house"}, Mots: []string{"das", "haus"}}}
	NewIBMModel1(corpus, 3)
	for _, p := range corpus[0].Alignment {
		if p[1] < -1 || p[1] >= len(corpus[0].Mots) {
			t.Errorf("invalid alignment position: %v", p)
		}
	}
}

func TestAlignmentStringSorted(t *testing.T) {
	a := Alignment{{2, 1}, {0, 0}, {1, 2}}
	got := a.String()
	parts := strings.Split(got, " ")
	for i := 1; i < len(parts); i++ {
		if parts[i-1] > parts[i] {
			t.Errorf("Alignment.String not sorted: %s", got)
		}
	}
}

func TestFmtS(t *testing.T) {
	got := fmtS("%d-%d", 1, 2)
	if got != "1-2" {
		t.Errorf("expected '1-2', got %s", got)
	}
}

func TestJoin(t *testing.T) {
	got := join([]string{"a", "b", "c"}, "-")
	if got != "a-b-c" {
		t.Errorf("expected 'a-b-c', got %s", got)
	}
}

func TestJoinEmpty(t *testing.T) {
	got := join([]string{}, "-")
	if got != "" {
		t.Errorf("expected empty string, got %s", got)
	}
}

func TestMax(t *testing.T) {
	if max(1, 2) != 2 {
		t.Error("expected max(1, 2) = 2")
	}
	if max(5, 3) != 5 {
		t.Error("expected max(5, 3) = 5")
	}
	if max(4, 4) != 4 {
		t.Error("expected max(4, 4) = 4")
	}
}

func TestAlignmentInvertTwice(t *testing.T) {
	a := Alignment{{0, 1}, {1, 0}}
	inv := a.Invert()
	doubleInv := inv.Invert()
	if len(doubleInv) != len(a) {
		t.Errorf("expected same length, got %d", len(doubleInv))
	}
	for i, p := range a {
		if doubleInv[i] != p {
			t.Errorf("at position %d: expected %v, got %v", i, p, doubleInv[i])
		}
	}
}

func TestIBMModel1LongSentence(t *testing.T) {
	corpus := []AlignedSent{{
		Words: []string{"this", "is", "a", "test", "sentence"},
		Mots:  []string{"das", "ist", "ein", "test", "satz"},
	}}
	m := NewIBMModel1(corpus, 2)
	if m == nil {
		t.Fatal("NewIBMModel1 returned nil")
	}
	if len(corpus[0].Alignment) != 5 {
		t.Errorf("expected 5 alignments, got %d", len(corpus[0].Alignment))
	}
}

func TestAlignedSentStruct(t *testing.T) {
	s := AlignedSent{
		Words:     []string{"a", "b"},
		Mots:      []string{"x", "y"},
		Alignment: Alignment{{0, 0}, {1, 1}},
	}
	if s.Words[0] != "a" || s.Mots[0] != "x" {
		t.Errorf("unexpected AlignedSent: %+v", s)
	}
}

func TestAlignmentType(t *testing.T) {
	var a Alignment = [][2]int{{0, 0}}
	if len(a) != 1 {
		t.Errorf("expected 1 element, got %d", len(a))
	}
}

func BenchmarkIBM1(b *testing.B) {
	corpus := []AlignedSent{
		{Words: []string{"the", "house"}, Mots: []string{"das", "haus"}},
		{Words: []string{"a", "book"}, Mots: []string{"ein", "buch"}},
		{Words: []string{"the", "cat"}, Mots: []string{"die", "katze"}},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewIBMModel1(corpus, 5)
	}
}
