package tgrep

import (
	"testing"

	"github.com/ateto3/gopherml/go-nlp/tree"
)

func TestFind(t *testing.T) {
	tr := tree.New("S", tree.New("NP", "I"), tree.New("VP", "saw"))
	if len(Find(tr, "NP")) != 1 {
		t.Fatal("bad")
	}
}

func TestFindLabelPrefix(t *testing.T) {
	tr := tree.New("S", tree.New("NP", "I"), tree.New("VP", "saw"))
	if len(Find(tr, "label=VP")) != 1 {
		t.Fatal("expected 1 match for label=VP")
	}
}

func TestFindNoMatch(t *testing.T) {
	tr := tree.New("S", tree.New("NP", "I"), tree.New("VP", "saw"))
	matches := Find(tr, "PP")
	if len(matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(matches))
	}
}

func TestFindMultiple(t *testing.T) {
	tr := tree.New("S",
		tree.New("NP", "I"),
		tree.New("VP",
			tree.New("V", "saw"),
			tree.New("NP", "the dog")))
	matches := Find(tr, "NP")
	if len(matches) < 2 {
		t.Errorf("expected at least 2 NP matches, got %d", len(matches))
	}
}

func TestFindRootNode(t *testing.T) {
	tr := tree.New("S", tree.New("NP", "I"))
	matches := Find(tr, "S")
	if len(matches) != 1 {
		t.Errorf("expected 1 match for 'S', got %d", len(matches))
	}
}

func TestFindEmptyTree(t *testing.T) {
	tr := tree.New("S")
	matches := Find(tr, "S")
	if len(matches) != 1 {
		t.Errorf("expected 1 match for empty tree, got %d", len(matches))
	}
}

func TestFindLabelExactMatch(t *testing.T) {
	tr := tree.New("S",
		tree.New("NP-SBJ", "I"),
		tree.New("NP", "dog"))
	matches := Find(tr, "label=NP")
	if len(matches) != 1 {
		t.Errorf("expected 1 match for exact 'NP', got %d", len(matches))
	}
}

func TestFindDeepTree(t *testing.T) {
	tr := tree.New("S",
		tree.New("NP",
			tree.New("DET", "the"),
			tree.New("NN", "cat")))
	matches := Find(tr, "DET")
	if len(matches) != 1 {
		t.Errorf("expected 1 DET match, got %d", len(matches))
	}
}

func TestFindCaseSensitive(t *testing.T) {
	tr := tree.New("S", tree.New("NP", "I"))
	matchesUpper := Find(tr, "NP")
	matchesLower := Find(tr, "np")
	if len(matchesUpper) != 1 {
		t.Errorf("expected 1 match for 'NP', got %d", len(matchesUpper))
	}
	if len(matchesLower) != 0 {
		t.Errorf("expected 0 matches for 'np' (case-sensitive), got %d", len(matchesLower))
	}
}

func TestFindComplexTree(t *testing.T) {
	tr, _ := tree.FromString("(S (NP (DT the) (JJ big) (NN dog)) (VP (V saw) (NP (DT a) (NN cat)))")
	matches := Find(tr, "DT")
	if len(matches) != 2 {
		t.Errorf("expected 2 DT matches, got %d", len(matches))
	}
}

func TestMatchBasic(t *testing.T) {
	tr := tree.New("S", "I")
	if !match(tr, "S") {
		t.Error("expected match for 'S'")
	}
	if match(tr, "NP") {
		t.Error("expected no match for 'NP'")
	}
}

func TestMatchLabelPrefixSyntax(t *testing.T) {
	tr := tree.New("S", "I")
	if !match(tr, "label=S") {
		t.Error("expected match for 'label=S'")
	}
}

func TestFindAllSameLabel(t *testing.T) {
	tr := tree.New("S",
		tree.New("NP", "I"),
		tree.New("VP",
			tree.New("V", "saw"),
			tree.New("NP", "dog")))
	matches := Find(tr, "NP")
	if len(matches) != 2 {
		t.Errorf("expected 2 NP matches, got %d", len(matches))
	}
}

func TestFindUnicode(t *testing.T) {
	tr := tree.New("句子", tree.New("名词", "猫"))
	matches := Find(tr, "名词")
	if len(matches) != 1 {
		t.Errorf("expected 1 match for '名词', got %d", len(matches))
	}
}

func TestFindNonTerminalWithLeaves(t *testing.T) {
	// Tree with nested structure
	tr := tree.New("ROOT",
		tree.New("S",
			tree.New("NP", "I"),
			tree.New("VP", "saw")))
	matches := Find(tr, "S")
	if len(matches) != 1 {
		t.Errorf("expected 1 match for 'S', got %d", len(matches))
	}
}

func TestFindLeafLabels(t *testing.T) {
	tr := tree.New("ROOT", "hello")
	// "hello" is stored as a string in Children, not as a Tree node
	// Find only finds Tree nodes, so expecting 0 is correct
	matches := Find(tr, "hello")
	if len(matches) != 0 {
		t.Errorf("expected 0 matches for leaf string, got %d", len(matches))
	}
}

func BenchmarkFind(b *testing.B) {
	tr, _ := tree.FromString("(S (NP (DT the) (JJ big) (NN dog)) (VP (V saw) (NP (DT a) (NN cat)))")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Find(tr, "NP")
	}
}
