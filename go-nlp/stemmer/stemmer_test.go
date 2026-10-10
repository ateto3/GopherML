package stemmer

import "testing"

func TestRegexpStemmer(t *testing.T) {
	s := NewRegexpStemmer(`(ing|s|e)$`, 4)
	if got := s.Stem("cars"); got != "car" {
		t.Fatalf("got %s", got)
	}
}

func TestPorterStemmer(t *testing.T) {
	s := NewPorterStemmer()
	if got := s.Stem("running"); got != "runn" {
		t.Fatalf("got %s", got)
	}
}

func TestRegexpStemmerBasicSuffixes(t *testing.T) {
	s := NewRegexpStemmer(`(ing|s)$`, 4)
	tests := []struct {
		input, expected string
	}{
		{"cars", "car"},
		{"running", "runn"},
	}
	for _, tt := range tests {
		if got := s.Stem(tt.input); got != tt.expected {
			t.Errorf("RegexpStemmer(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestRegexpStemmerMinLength(t *testing.T) {
	s := NewRegexpStemmer(`(s)$`, 5)
	if got := s.Stem("cars"); got != "cars" {
		t.Errorf("expected 'cars' unchanged (len 4 < min 5), got %q", got)
	}
}

func TestRegexpStemmerNoMatch(t *testing.T) {
	s := NewRegexpStemmer(`(xyz)$`, 4)
	if got := s.Stem("hello"); got != "hello" {
		t.Errorf("expected 'hello' unchanged, got %q", got)
	}
}

func TestRegexpStemmerEmpty(t *testing.T) {
	s := NewRegexpStemmer(`(s)$`, 4)
	if got := s.Stem(""); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestPorterStemmerRunning(t *testing.T) {
	s := NewPorterStemmer()
	if got := s.Stem("running"); got != "runn" {
		t.Errorf("expected 'runn', got %q", got)
	}
}

func TestPorterStemmerCaring(t *testing.T) {
	s := NewPorterStemmer()
	if got := s.Stem("caring"); got != "car" {
		t.Errorf("expected 'car', got %q", got)
	}
}

func TestPorterStemmerIES(t *testing.T) {
	s := NewPorterStemmer()
	if got := s.Stem("ponies"); got != "pony" {
		t.Errorf("expected 'pony', got %q", got)
	}
}

func TestPorterStemmerEmpty(t *testing.T) {
	s := NewPorterStemmer()
	if got := s.Stem(""); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestPorterStemmerShortWord(t *testing.T) {
	s := NewPorterStemmer()
	if got := s.Stem("is"); got != "is" {
		t.Errorf("expected 'is' unchanged, got %q", got)
	}
}

func TestPorterStemmerCaseInsensitive(t *testing.T) {
	s := NewPorterStemmer()
	got1 := s.Stem("Running")
	got2 := s.Stem("running")
	if got1 != got2 {
		t.Errorf("expected case-insensitive: %q vs %q", got1, got2)
	}
}

func TestPorterStemmerDeterministic(t *testing.T) {
	s := NewPorterStemmer()
	got1 := s.Stem("running")
	got2 := s.Stem("running")
	if got1 != got2 {
		t.Errorf("expected deterministic result: %q vs %q", got1, got2)
	}
}

func TestNewPorterStemmer(t *testing.T) {
	s := NewPorterStemmer()
	_ = s
}

func TestNewRegexpStemmer(t *testing.T) {
	s := NewRegexpStemmer(`(s)$`, 4)
	if s == nil {
		t.Fatal("NewRegexpStemmer returned nil")
	}
}

func TestStemmerInterface(t *testing.T) {
	var _ Stemmer = NewPorterStemmer()
	var _ Stemmer = NewRegexpStemmer(`(s)$`, 4)
}

func TestRegexpStemmerVariousPatterns(t *testing.T) {
	s := NewRegexpStemmer(`(s)$`, 4)
	if got := s.Stem("dogs"); got != "dog" {
		t.Errorf("expected 'dog', got %q", got)
	}
}

func TestPorterStemmerED(t *testing.T) {
	s := NewPorterStemmer()
	if got := s.Stem("walked"); got != "walk" {
		t.Errorf("expected 'walk', got %q", got)
	}
}

func TestPorterStemmerEndingS(t *testing.T) {
	s := NewPorterStemmer()
	if got := s.Stem("dogs"); got != "dog" {
		t.Errorf("expected 'dog', got %q", got)
	}
}

func BenchmarkPorterStemmer(b *testing.B) {
	s := NewPorterStemmer()
	words := []string{"running", "cats", "ponies", "playing", "walked", "caring"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Stem(words[i%len(words)])
	}
}

func BenchmarkRegexpStemmer(b *testing.B) {
	s := NewRegexpStemmer(`(ing|s|ed|ies)$`, 4)
	words := []string{"running", "cats", "ponies", "playing", "walked"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Stem(words[i%len(words)])
	}
}
