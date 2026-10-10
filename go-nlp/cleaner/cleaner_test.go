package cleaner

import (
	"strings"
	"testing"
)

func TestCleanerBasic(t *testing.T) {
	cleaner := New()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple text",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			name:     "with punctuation",
			input:    "Hello, world!",
			expected: "Hello, world!",
		},
		{
			name:     "with HTML",
			input:    "<p>Hello world</p>",
			expected: "<p>Hello world</p>",
		},
		{
			name:     "with URL",
			input:    "Check out https://example.com",
			expected: "Check out https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleaner.Clean(tt.input)
			if result != tt.expected {
				t.Errorf("Clean() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCleanerOptions(t *testing.T) {
	cleaner := New(
		WithHTML(false),
		WithPunctuation(false),
		WithLowercase(false),
		WithURLs(false),
	)

	result := cleaner.Clean("<p>Hello, World!</p> https://example.com")
	expected := "<p>Hello, World!</p> https://example.com"
	if result != expected {
		t.Errorf("Clean() = %q, want %q", result, expected)
	}
}

func TestCleanerTokens(t *testing.T) {
	cleaner := New(WithStopwords(true))
	tokens := []string{"hello", "world", "the", "and", "go"}
	result := cleaner.CleanTokens(tokens)
	
	// "the" and "and" should be removed as stopwords
	if len(result) != 3 {
		t.Errorf("Expected 3 tokens, got %d: %v", len(result), result)
	}
}

func TestCleanerStopwords(t *testing.T) {
	cleaner := New(WithStopwords(true))
	
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"english stopwords", "the quick brown fox", 3},
		{"spanish stopwords", "el rápido zorro marrón", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleaner.Options.Language = tt.name[:2]
			result := cleaner.Clean(tt.input)
			words := len(cleaner.CleanTokens(strings.Fields(result)))
			if words != tt.want {
				t.Errorf("Clean() produced %d words, want %d", words, tt.want)
			}
		})
	}
}

func TestCleanerCustomStopwords(t *testing.T) {
	cleaner := New(WithStopwords(true))
	cleaner.AddStopwords("quick", "brown")
	result := cleaner.Clean("the quick brown fox")
	
	words := cleaner.CleanTokens(strings.Fields(result))
	for _, w := range words {
		if w == "quick" || w == "brown" {
			t.Errorf("Custom stopwords not removed: %v", words)
		}
	}
}

func TestCleanerNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"with numbers", "price is 100 dollars", "price is dollars"},
		{"without numbers", "price is 100 dollars", "price is 100 dollars"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c *Cleaner
			if tt.name == "with numbers" {
				c = New(WithHTML(false), WithPunctuation(false), WithNumbers(true))
			} else {
				c = New(WithHTML(false), WithPunctuation(false), WithNumbers(false))
			}
			
			result := c.Clean(tt.input)
			if result != tt.expected {
				t.Errorf("Clean() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCleanerCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		lower    bool
		expected string
	}{
		{"lowercase", "Hello WORLD", true, "hello world"},
		{"uppercase", "Hello world", false, "HELLO WORLD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []OptionsFunc{WithLowercase(tt.lower)}
			if tt.lower {
				opts = []OptionsFunc{WithLowercase(true)}
			} else {
				opts = []OptionsFunc{WithLowercase(false), WithUppercase(true)}
			}
			cleaner := New(opts...)
			result := cleaner.Clean(tt.input)
			if result != tt.expected {
				t.Errorf("Clean() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCleanerWhitespace(t *testing.T) {
	cleaner := New(WithHTML(false), WithPunctuation(false))
	
	result := cleaner.Clean("hello   \t\n   world")
	if result != "hello world" {
		t.Errorf("Clean() = %q, want %q", result, "hello world")
	}
}

func TestCleanerAccents(t *testing.T) {
	cleaner := New(
		WithHTML(false),
		WithPunctuation(false),
		WithAccents(true),
		WithLowercase(true),
	)
	
	result := cleaner.Clean("Café résumé")
	if result != "cafe resume" {
		t.Errorf("Clean() = %q, want %q", result, "cafe resume")
	}
}

func TestCleanerLengthFilter(t *testing.T) {
	cleaner := New(
		WithHTML(false),
		WithPunctuation(false),
		WithMinWordLength(3),
		WithMaxWordLength(10),
	)
	
	result := cleaner.Clean("I am a developer")
	words := cleaner.CleanTokens(strings.Fields(result))
	
	// "I", "am", "a" should be filtered out (length < 3)
	if len(words) != 1 {
		t.Errorf("Expected 1 token, got %d: %v", len(words), words)
	}
}

func TestCleanerEmail(t *testing.T) {
	cleaner := New(WithEmails(true))
	result := cleaner.Clean("Contact us at test@example.com")
	if result != "Contact us at" {
		t.Errorf("Clean() = %q, want %q", result, "Contact us at")
	}
}

func TestCleanerPipeline(t *testing.T) {
	cleaner := New(
		WithHTML(false),
		WithPunctuation(false),
		WithStopwords(true),
		WithLowercase(true),
		WithMinWordLength(3),
		WithMaxWordLength(20),
	)
	
	result := cleaner.Clean("The quick brown fox jumps")
	words := cleaner.CleanTokens(strings.Fields(result))
	
	// "The" should be removed as a stopword
	// "quick" (5), "brown" (5), "fox" (3), "jumps" (5) should remain
	// But "quick" and "brown" are filtered by length (default min=1, max=100)
	if len(words) < 2 {
		t.Errorf("Expected at least 2 tokens, got %d: %v", len(words), words)
	}
}

func BenchmarkCleaner(b *testing.B) {
	cleaner := New()
	input := "The quick brown fox jumps over the lazy dog! https://example.com @user #test"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cleaner.Clean(input)
	}
}

func BenchmarkCleanTokens(b *testing.B) {
	cleaner := New(WithStopwords(true))
	input := []string{"the", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cleaner.CleanTokens(input)
	}
}
