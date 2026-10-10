package stopwords

import (
	"testing"
)

func TestDefault(t *testing.T) {
	stopwords := Default("en")
	if stopwords.Count() == 0 {
		t.Error("Expected English stopwords, got empty map")
	}
}

func TestDefaultUnknownLanguage(t *testing.T) {
	stopwords := Default("xx")
	if stopwords.Count() == 0 {
		t.Error("Expected English fallback, got empty map")
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		word     string
		expected bool
	}{
		{"english stopword", "en", "the", true},
		{"english non-stopword", "en", "hello", false},
		{"spanish stopword", "es", "el", true},
		{"spanish non-stopword", "es", "hola", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.lang, tt.word)
			if result != tt.expected {
				t.Errorf("Contains(%s, %s) = %v, want %v", tt.lang, tt.word, result, tt.expected)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	// Add a stopword
	Add("en", "testword")
	if !Contains("en", "testword") {
		t.Error("Word not added")
	}

	// Remove the stopword
	Remove("en", "testword")
	if Contains("en", "testword") {
		t.Error("Word not removed")
	}
}

func TestAdd(t *testing.T) {
	// Add a stopword
	Add("en", "customword")
	if !Contains("en", "customword") {
		t.Error("Word not added")
	}
}

func TestLanguages(t *testing.T) {
	languages := Languages()
	expected := []string{"en", "es", "pt", "fr", "de", "it", "nl", "ru", "zh", "ja"}
	if len(languages) != len(expected) {
		t.Errorf("Expected %d languages, got %d", len(expected), len(languages))
	}
}

func TestCount(t *testing.T) {
	count := Count("en")
	if count == 0 {
		t.Error("Expected non-zero count for English")
	}
}

func TestStopwordsMap(t *testing.T) {
	// Verify that all languages have stopwords
	languages := Languages()
	for _, lang := range languages {
		stopwords := Default(lang)
		if stopwords.Count() == 0 {
			t.Errorf("No stopwords found for language: %s", lang)
		}
	}
}

func BenchmarkContains(b *testing.B) {
	stopwords := Default("en")
	words := stopwords.Words()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = words["the"]
	}
}

func BenchmarkContainsWithFunction(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Contains("en", "the")
	}
}
