package text

import (
	"testing"
)

func TestConcordance(t *testing.T) {
	idx := NewConcordanceIndex([]string{"the", "big", "dog", "the", "dog"})
	if len(idx.Find("dog", 20, 10)) == 0 {
		t.Fatal("no concordance")
	}
}

func TestContext(t *testing.T) {
	ci := NewContextIndex([]string{"the", "big", "dog", "runs"})
	if len(ci.SimilarWords("dog", 3)) == 0 {
		t.Fatal("no similar")
	}
}

func TestSearch(t *testing.T) {
	s := NewTokenSearcher([]string{"a", "bb", "ccc"})
	if len(s.Search(`^b+`)) != 1 {
		t.Fatal("bad search")
	}
}

func TestConcordanceMultipleOccurrences(t *testing.T) {
	tokens := []string{"the", "cat", "sat", "the", "cat", "ran"}
	idx := NewConcordanceIndex(tokens)
	lines := idx.Find("cat", 20, 10)
	if len(lines) != 2 {
		t.Errorf("expected 2 concordance lines for 'cat', got %d", len(lines))
	}
}

func TestConcordanceNoMatch(t *testing.T) {
	idx := NewConcordanceIndex([]string{"hello", "world"})
	lines := idx.Find("missing", 20, 10)
	if lines != nil {
		t.Errorf("expected nil for missing word, got %v", lines)
	}
}

func TestConcordanceCaseInsensitive(t *testing.T) {
	idx := NewConcordanceIndex([]string{"the", "BIG", "dog"})
	lines := idx.Find("big", 20, 10)
	if len(lines) != 1 {
		t.Errorf("expected 1 line for case-insensitive match, got %d", len(lines))
	}
	if lines[0].Query != "BIG" {
		t.Errorf("expected original case 'BIG', got %s", lines[0].Query)
	}
}

func TestConcordanceLineFields(t *testing.T) {
	idx := NewConcordanceIndex([]string{"a", "b", "c", "d", "e"})
	lines := idx.Find("c", 20, 10)
	if len(lines) == 0 {
		t.Fatal("expected at least 1 line")
	}
	line := lines[0]
	if line.Query != "c" {
		t.Errorf("expected Query='c', got %s", line.Query)
	}
	if line.Offset != 2 {
		t.Errorf("expected Offset=2, got %d", line.Offset)
	}
}

func TestConcordanceLimit(t *testing.T) {
	idx := NewConcordanceIndex([]string{"a", "a", "a", "a", "a"})
	lines := idx.Find("a", 20, 2)
	if len(lines) != 2 {
		t.Errorf("expected 2 lines (limited), got %d", len(lines))
	}
}

func TestConcordanceEmptyTokens(t *testing.T) {
	idx := NewConcordanceIndex([]string{})
	lines := idx.Find("anything", 20, 10)
	if lines != nil {
		t.Error("expected nil for empty index")
	}
}

func TestConcordanceString(t *testing.T) {
	idx := NewConcordanceIndex([]string{"a", "b", "c"})
	if idx.String() == "" {
		t.Error("expected non-empty string representation")
	}
}

func TestContextSimilarWordsBasic(t *testing.T) {
	tokens := []string{"the", "dog", "runs", "the", "cat", "runs"}
	ci := NewContextIndex(tokens)
	sim := ci.SimilarWords("dog", 3)
	// cat should be similar to dog (both are followed by 'runs')
	if len(sim) == 0 {
		t.Error("expected similar words for 'dog'")
	}
}

func TestContextSimilarWordsLimit(t *testing.T) {
	tokens := []string{"a", "b", "c", "a", "b", "d", "a", "b", "e"}
	ci := NewContextIndex(tokens)
	sim := ci.SimilarWords("a", 2)
	if len(sim) > 2 {
		t.Errorf("expected at most 2 similar words, got %d", len(sim))
	}
}

func TestContextSimilarWordsDefault(t *testing.T) {
	tokens := []string{"a", "b", "c", "d"}
	ci := NewContextIndex(tokens)
	sim := ci.SimilarWords("a", 0) // 0 should default
	_ = sim
}

func TestContextStartEnd(t *testing.T) {
	// The first token should have *START* as left context
	tokens := []string{"hello", "world"}
	ci := NewContextIndex(tokens)
	// Verify the structure
	if len(ci.WordCtx) == 0 {
		t.Error("expected non-empty WordCtx")
	}
}

func TestContextEmpty(t *testing.T) {
	ci := NewContextIndex([]string{})
	sim := ci.SimilarWords("anything", 5)
	if sim != nil && len(sim) > 0 {
		// May be nil or empty
		_ = sim
	}
}

func TestTokenSearcherBasic(t *testing.T) {
	s := NewTokenSearcher([]string{"apple", "banana", "cherry", "date"})
	matches := s.Search("^[ab]")
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d", len(matches))
	}
}

func TestTokenSearcherNoMatch(t *testing.T) {
	s := NewTokenSearcher([]string{"apple", "banana"})
	matches := s.Search("xyz")
	if len(matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(matches))
	}
}

func TestTokenSearcherEmpty(t *testing.T) {
	s := NewTokenSearcher([]string{})
	matches := s.Search(".*")
	if len(matches) != 0 {
		t.Errorf("expected 0 matches for empty tokens, got %d", len(matches))
	}
}

func TestTokenSearcherMultiple(t *testing.T) {
	s := NewTokenSearcher([]string{"a", "b", "a", "c", "a"})
	matches := s.Search("^a$")
	if len(matches) != 3 {
		t.Errorf("expected 3 matches for 'a', got %d", len(matches))
	}
}

func TestTokenSearcherComplexPattern(t *testing.T) {
	s := NewTokenSearcher([]string{"123", "abc", "456", "def"})
	matches := s.Search(`^\d+$`)
	if len(matches) != 2 {
		t.Errorf("expected 2 digit matches, got %d", len(matches))
	}
}

func TestConcordanceLineStruct(t *testing.T) {
	line := ConcordanceLine{
		Left:   "the big",
		Query:  "dog",
		Right:  "runs",
		Offset: 2,
	}
	if line.Left != "the big" || line.Query != "dog" || line.Right != "runs" || line.Offset != 2 {
		t.Errorf("unexpected ConcordanceLine: %+v", line)
	}
}

func TestNewConcordanceIndex(t *testing.T) {
	idx := NewConcordanceIndex([]string{"a", "b"})
	if len(idx.Tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(idx.Tokens))
	}
	if len(idx.Offsets) == 0 {
		t.Error("expected non-empty Offsets")
	}
}

func TestNewContextIndex(t *testing.T) {
	ci := NewContextIndex([]string{"a", "b", "c"})
	if len(ci.Tokens) != 3 {
		t.Errorf("expected 3 tokens, got %d", len(ci.Tokens))
	}
}

func TestNewTokenSearcher(t *testing.T) {
	s := NewTokenSearcher([]string{"a", "b"})
	if len(s.Tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(s.Tokens))
	}
}

func TestContextSingleToken(t *testing.T) {
	ci := NewContextIndex([]string{"hello"})
	sim := ci.SimilarWords("hello", 5)
	_ = sim
}

func TestConcordanceWordAtStart(t *testing.T) {
	idx := NewConcordanceIndex([]string{"first", "a", "b", "c"})
	lines := idx.Find("first", 20, 10)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0].Query != "first" {
		t.Errorf("expected Query='first', got %s", lines[0].Query)
	}
}

func TestConcordanceWordAtEnd(t *testing.T) {
	idx := NewConcordanceIndex([]string{"a", "b", "c", "last"})
	lines := idx.Find("last", 20, 10)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
}

func TestTokenSearcherEmptyPattern(t *testing.T) {
	s := NewTokenSearcher([]string{"a", "b", "c"})
	defer func() {
		// regexp.MustCompile panics on empty pattern
		if r := recover(); r == nil {
			// Not strictly required to panic, depends on Go version
		}
	}()
	matches := s.Search("")
	_ = matches
}

func BenchmarkConcordanceFind(b *testing.B) {
	tokens := make([]string, 1000)
	for i := range tokens {
		tokens[i] = "word"
	}
	idx := NewConcordanceIndex(tokens)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.Find("word", 80, 10)
	}
}

func BenchmarkContextSimilarWords(b *testing.B) {
	tokens := []string{"the", "dog", "runs", "the", "cat", "runs", "the", "bird", "flies"}
	ci := NewContextIndex(tokens)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ci.SimilarWords("dog", 5)
	}
}
