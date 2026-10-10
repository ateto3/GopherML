package lm

import (
	"strings"
	"testing"
)

func TestLM(t *testing.T) {
	m := NewMLE()
	m.Fit([]Sentence{{"a", "b", "c"}, {"a", "b"}})
	if m.Score("b", "a") <= 0 {
		t.Fatal("score")
	}
	if len(m.Generate("a", 2)) == 0 {
		t.Fatal("gen")
	}
}

func TestNewMLE(t *testing.T) {
	m := NewMLE()
	if m == nil {
		t.Fatal("NewMLE returned nil")
	}
	if m.Counts == nil {
		t.Fatal("expected non-nil Counts")
	}
	if m.Counts.Unigrams == nil {
		t.Error("expected Unigrams map to be initialized")
	}
	if m.Counts.Bigrams == nil {
		t.Error("expected Bigrams map to be initialized")
	}
}

func TestFitPopulatesCounts(t *testing.T) {
	m := NewMLE()
	m.Fit([]Sentence{{"a", "b", "c"}})
	if m.Counts.Unigrams["a"] != 1 {
		t.Errorf("expected 'a' count 1, got %d", m.Counts.Unigrams["a"])
	}
}

func TestMLEScoreUnigram(t *testing.T) {
	m := NewMLE()
	m.Fit([]Sentence{{"a", "b", "a", "b"}})
	score := m.Score("a")
	if score != 0.5 {
		t.Errorf("expected score 0.5, got %f", score)
	}
}

func TestMLEScoreBigram(t *testing.T) {
	m := NewMLE()
	m.Fit([]Sentence{{"a", "b", "a", "b"}})
	score := m.Score("b", "a")
	if score != 1.0 {
		t.Errorf("expected score 1.0, got %f", score)
	}
}

func TestMLEScoreUnseenUnigram(t *testing.T) {
	m := NewMLE()
	m.Fit([]Sentence{{"a", "b"}})
	if m.Score("z") != 0 {
		t.Errorf("expected score 0 for unseen word, got %f", m.Score("z"))
	}
}

func TestMLEScoreUnseenBigram(t *testing.T) {
	m := NewMLE()
	m.Fit([]Sentence{{"a", "b"}})
	if m.Score("z", "a") != 0 {
		t.Errorf("expected score 0 for unseen bigram, got %f", m.Score("z", "a"))
	}
}

func TestMLEGenerate(t *testing.T) {
	m := NewMLE()
	m.Fit([]Sentence{{"a", "b", "c"}})
	got := m.Generate("a", 3)
	if len(got) != 3 {
		t.Errorf("expected 3 tokens, got %d", len(got))
	}
	if got[0] != "a" {
		t.Errorf("expected first token 'a', got %q", got[0])
	}
}

func TestMLEGenerateEmpty(t *testing.T) {
	m := NewMLE()
	got := m.Generate("a", 0)
	if len(got) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(got))
	}
}

func TestMLEGenerateUnknownSeed(t *testing.T) {
	m := NewMLE()
	m.Fit([]Sentence{{"a", "b"}})
	got := m.Generate("z", 5)
	if len(got) == 0 {
		t.Error("expected at least 1 token")
	}
	if got[0] != "z" {
		t.Errorf("expected first token 'z', got %q", got[0])
	}
}

func TestLaplaceScoreUnigram(t *testing.T) {
	l := NewLaplace()
	l.Fit([]Sentence{{"a", "b"}})
	score := l.Score("a")
	if score != 0.5 {
		t.Errorf("expected 0.5, got %f", score)
	}
}

func TestLaplaceScoreUnseenWord(t *testing.T) {
	l := NewLaplace()
	l.Fit([]Sentence{{"a", "b"}})
	score := l.Score("z")
	if score != 0.25 {
		t.Errorf("expected 0.25 for unseen word, got %f", score)
	}
}

func TestRandomChoice(t *testing.T) {
	weights := map[string]float64{"a": 0.5, "b": 0.3, "c": 0.1}
	got := RandomChoice(weights)
	if _, ok := weights[got]; !ok {
		t.Errorf("expected choice from weights, got %q", got)
	}
	if got != "a" {
		t.Errorf("expected argmax 'a', got %q", got)
	}
}

func TestRandomChoiceEmpty(t *testing.T) {
	got := RandomChoice(map[string]float64{})
	if got != "" {
		t.Errorf("expected empty string for empty weights, got %q", got)
	}
}

func TestShuffleSentences(t *testing.T) {
	original := []Sentence{{"a", "b"}, {"c", "d"}}
	shuffled := ShuffleSentences(original)
	if len(shuffled) != len(original) {
		t.Errorf("expected same length, got %d", len(shuffled))
	}
}

func TestNewCounter(t *testing.T) {
	c := NewCounter()
	if c == nil {
		t.Fatal("NewCounter returned nil")
	}
	if c.Unigrams == nil {
		t.Error("expected Unigrams initialized")
	}
}

func TestCounterUpdate(t *testing.T) {
	c := NewCounter()
	c.Update([]Sentence{{"a", "b", "c"}})
	if c.Unigrams["a"] != 1 || c.Unigrams["b"] != 1 || c.Unigrams["c"] != 1 {
		t.Errorf("unexpected unigram counts: %v", c.Unigrams)
	}
}

func TestMLEScoreWithEmptyCorpus(t *testing.T) {
	m := NewMLE()
	m.Fit(nil)
	if m.Score("a") != 0 {
		t.Errorf("expected 0 for empty corpus, got %f", m.Score("a"))
	}
}

func TestMLEProbabilitiesSumToOne(t *testing.T) {
	m := NewMLE()
	m.Fit([]Sentence{{"a", "b", "c"}})
	total := m.Score("a") + m.Score("b") + m.Score("c")
	if total < 0.99 || total > 1.01 {
		t.Errorf("expected probs to sum to ~1, got %f", total)
	}
}

func TestSentenceType(t *testing.T) {
	s := Sentence{"hello", "world"}
	if len(s) != 2 || s[0] != "hello" {
		t.Errorf("unexpected sentence: %v", s)
	}
	if !strings.Contains(s[0], "hello") {
		t.Error("expected 'hello' in sentence")
	}
}

func BenchmarkMLEScore(b *testing.B) {
	m := NewMLE()
	m.Fit([]Sentence{{"a", "b", "c", "a", "b", "c", "a", "b"}})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Score("b", "a")
	}
}

func BenchmarkMLEGenerate(b *testing.B) {
	m := NewMLE()
	m.Fit([]Sentence{{"a", "b", "c"}})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Generate("a", 10)
	}
}
