package collocations

import (
	"math"
	"testing"
)

func TestBigramFinderPMI(t *testing.T) {
	f := NewBigramFinder([]string{"new", "york", "new", "york"})
	scores := f.PMI()
	if len(scores) == 0 {
		t.Fatal("expected scores")
	}
}

func TestBigramFinderFrequentCollocation(t *testing.T) {
	f := NewBigramFinder([]string{"new", "york", "new", "york", "new", "york", "the", "cat"})
	scores := f.PMI()
	if len(scores) == 0 {
		t.Fatal("expected PMI scores")
	}
	found := false
	for _, s := range scores {
		if len(s.NGram) == 2 && s.NGram[0] == "new" && s.NGram[1] == "york" {
			found = true
			if s.Score <= 0 {
				t.Errorf("expected positive PMI for 'new york' collocation, got %f", s.Score)
			}
		}
	}
	if !found {
		t.Error("expected 'new york' bigram in PMI scores")
	}
}

func TestBigramFinderFilter(t *testing.T) {
	f := NewBigramFinder([]string{"a", "b", "a", "b", "c", "d"})
	f = f.Filter(2)
	scores := f.PMI()
	for _, s := range scores {
		if s.NGram[0] == "c" || s.NGram[0] == "d" {
			t.Errorf("expected 'c|d' to be filtered out, got %v", s.NGram)
		}
	}
}

func TestBigramFinderPMISorted(t *testing.T) {
	f := NewBigramFinder([]string{"x", "y", "x", "y", "x", "y", "a", "b"})
	scores := f.PMI()
	if len(scores) < 2 {
		t.Fatal("expected at least 2 scores")
	}
	for i := 1; i < len(scores); i++ {
		if scores[i].Score > scores[i-1].Score {
			t.Errorf("PMI scores not sorted in descending order at index %d", i)
		}
	}
}

func TestBigramFinderEmpty(t *testing.T) {
	f := NewBigramFinder([]string{})
	scores := f.PMI()
	if len(scores) != 0 {
		t.Errorf("expected no scores for empty input, got %d", len(scores))
	}
}

func TestBigramFinderSingleWord(t *testing.T) {
	f := NewBigramFinder([]string{"hello"})
	scores := f.PMI()
	if len(scores) != 0 {
		t.Errorf("expected no bigrams for single word, got %d", len(scores))
	}
}

func TestTrigramFinder(t *testing.T) {
	f := NewTrigramFinder([]string{"the", "quick", "brown", "fox", "the", "quick", "brown", "fox"})
	if f.T["the|quick|brown"] != 2 {
		t.Errorf("expected 'the quick brown' count 2, got %d", f.T["the|quick|brown"])
	}
}

func TestTrigramTop(t *testing.T) {
	f := NewTrigramFinder([]string{"a", "b", "c", "a", "b", "c", "d", "e", "f"})
	top := f.Top(2)
	if len(top) != 2 {
		t.Errorf("expected 2 top trigrams, got %d", len(top))
	}
	if top[0].Score < top[1].Score {
		t.Error("top trigrams should be sorted by score")
	}
}

func TestTrigramTopMoreThanAvailable(t *testing.T) {
	f := NewTrigramFinder([]string{"a", "b", "c"})
	top := f.Top(10)
	if len(top) != 1 {
		t.Errorf("expected 1 trigram, got %d", len(top))
	}
}

func TestTrigramTopEmpty(t *testing.T) {
	f := NewTrigramFinder([]string{})
	top := f.Top(5)
	if len(top) != 0 {
		t.Errorf("expected 0 trigrams for empty input, got %d", len(top))
	}
}

func TestTrigramFinderEmpty(t *testing.T) {
	f := NewTrigramFinder([]string{"hello"})
	if len(f.T) != 0 {
		t.Errorf("expected 0 trigrams for single word, got %d", len(f.T))
	}
}

func TestScoredNGram(t *testing.T) {
	s := ScoredNGram{NGram: []string{"a", "b"}, Score: 1.5}
	if s.NGram[0] != "a" || s.Score != 1.5 {
		t.Errorf("unexpected scored ngram: %v", s)
	}
}

func TestSortScored(t *testing.T) {
	s := []ScoredNGram{
		{NGram: []string{"c", "d"}, Score: 1.0},
		{NGram: []string{"a", "b"}, Score: 3.0},
		{NGram: []string{"e", "f"}, Score: 2.0},
	}
	sortScored(s)
	if s[0].Score != 3.0 || s[1].Score != 2.0 || s[2].Score != 1.0 {
		t.Errorf("scores not sorted in descending order: %v", s)
	}
}

func TestBigramFinderConcurrent(t *testing.T) {
	docs := [][]string{
		{"a", "b", "c", "a", "b"},
		{"b", "c", "d", "b", "c"},
		{"a", "b", "a", "b", "a"},
	}
	f := BigramFinderConcurrent(docs)
	if len(f.B) == 0 {
		t.Fatal("expected bigrams from concurrent finder")
	}
	if f.B["a|b"] != 4 {
		t.Errorf("expected 'a|b' count 4, got %d", f.B["a|b"])
	}
	if f.B["b|c"] != 3 {
		t.Errorf("expected 'b|c' count 3, got %d", f.B["b|c"])
	}
}

func TestBigramFinderConcurrentEmpty(t *testing.T) {
	f := BigramFinderConcurrent(nil)
	if len(f.B) != 0 {
		t.Errorf("expected empty bigrams, got %d", len(f.B))
	}
}

func TestSplit2(t *testing.T) {
	got := split2("a|b|c")
	if len(got) != 2 || got[0] != "a" || got[1] != "b|c" {
		t.Errorf("split2 failed: %v", got)
	}
}

func TestSplit3(t *testing.T) {
	got := split3("a|b|c|d")
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c|d" {
		t.Errorf("split3 failed: %v", got)
	}
}

func TestJoin(t *testing.T) {
	got := join([]string{"a", "b", "c"})
	if got != "a|b|c" {
		t.Errorf("join failed: %s", got)
	}
}

func TestLog2(t *testing.T) {
	if math.Abs(log2(2)-1) > 1e-10 {
		t.Errorf("log2(2) should be 1, got %f", log2(2))
	}
	if math.Abs(log2(4)-2) > 1e-10 {
		t.Errorf("log2(4) should be 2, got %f", log2(4))
	}
	if math.Abs(log2(1)-0) > 1e-10 {
		t.Errorf("log2(1) should be 0, got %f", log2(1))
	}
}

func TestPMIPositive(t *testing.T) {
	// PMI should be positive for a collocation vs random pairing
	words := []string{"new", "york", "new", "york", "new", "york", "the", "cat"}
	f := NewBigramFinder(words)
	scores := f.PMI()
	for _, s := range scores {
		if s.NGram[0] == "new" && s.NGram[1] == "york" {
			// PMI of "new york" relative to its marginal probabilities
			// should be positive (they co-occur more than chance)
			if s.Score <= 0 {
				t.Errorf("expected positive PMI, got %f", s.Score)
			}
			return
		}
	}
	t.Error("'new york' bigram not found")
}

func TestBigramFinderPreservesInput(t *testing.T) {
	words := []string{"a", "b", "c"}
	f := NewBigramFinder(words)
	if len(words) != 3 {
		t.Error("input words should not be modified")
	}
	if f.W["a"] != 1 {
		t.Errorf("expected 'a' count 1, got %d", f.W["a"])
	}
}
