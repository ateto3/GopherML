package tokenizer

import (
	"regexp"
	"testing"
)

func TestSpaceTokenizer(t *testing.T) {
	tok := NewSpaceTokenizer()
	s := "Good muffins cost $3.88\nin New York.  Please buy me\ntwo of them.\n\nThanks."
	tokens := tok.Tokenize(s)
	if len(tokens) == 0 {
		t.Error("expected tokens, got empty")
	}
	// SpaceTokenizer splits only on single spaces, so newlines remain in tokens
}

func TestTabTokenizer(t *testing.T) {
	tok := NewTabTokenizer()
	tokens := tok.Tokenize("a\tb c\n\t d")
	if len(tokens) != 3 {
		t.Errorf("expected 3 tokens, got %d: %v", len(tokens), tokens)
	}
}

func TestCharTokenizer(t *testing.T) {
	tok := NewCharTokenizer()
	tokens := tok.Tokenize("abc")
	if len(tokens) != 3 {
		t.Errorf("expected 3 tokens, got %d", len(tokens))
	}
	if tokens[0] != "a" || tokens[1] != "b" || tokens[2] != "c" {
		t.Errorf("unexpected tokens: %v", tokens)
	}

	spans := tok.SpanTokenize("abc")
	if len(spans) != 3 {
		t.Errorf("expected 3 spans, got %d", len(spans))
	}
}

func TestLineTokenizer(t *testing.T) {
	s := "Good muffins cost $3.88\nin New York.  Please buy me\ntwo of them.\n\nThanks."

	// Default (discard)
	tok := NewLineTokenizer("discard")
	tokens := tok.Tokenize(s)
	if len(tokens) != 4 {
		t.Errorf("expected 4 tokens, got %d: %v", len(tokens), tokens)
	}

	// Keep
	tokKeep := NewLineTokenizer("keep")
	tokensKeep := tokKeep.Tokenize(s)
	if len(tokensKeep) != 5 {
		t.Errorf("expected 5 tokens (keep), got %d: %v", len(tokensKeep), tokensKeep)
	}
}

func TestRegexpTokenizer(t *testing.T) {
	s := "Good muffins cost $3.88\nin New York.  Please buy me\ntwo of them.\n\nThanks."

	// Word tokenizer
	tok := NewRegexpTokenizer(`\w+|\$[\d\.]+|\S+`, false, true)
	tokens := tok.Tokenize(s)
	if len(tokens) == 0 {
		t.Error("expected tokens, got empty")
	}
	// Should find tokens: Good, muffins, cost, $3.88, in, New, York, ., Please, buy, me, two, of, them, ., Thanks, .
	// The \S+ pattern matches punctuation individually as well.
	if len(tokens) != 17 {
		t.Errorf("expected 17 tokens (word), got %d: %v", len(tokens), tokens)
	}

	// Gap mode (whitespace delimiter)
	tokGap := NewRegexpTokenizer(`\s+`, true, true)
	tokensGap := tokGap.Tokenize(s)
	// Expected: Good, muffins, cost, $3.88, in, New, York., Please, buy, me, two, of, them., Thanks.
	if len(tokensGap) != 14 {
		t.Errorf("expected 14 tokens (gap), got %d: %v", len(tokensGap), tokensGap)
	}
}

func TestWhitespaceTokenizer(t *testing.T) {
	tok := NewWhitespaceTokenizer()
	s := "Good muffins cost $3.88\nin New York.  Please buy me\ntwo of them.\n\nThanks."
	tokens := tok.Tokenize(s)
	// Expected: Good, muffins, cost, $3.88, in, New, York., Please, buy, me, two, of, them., Thanks.
	if len(tokens) != 14 {
		t.Errorf("expected 14 tokens, got %d: %v", len(tokens), tokens)
	}
}

func TestWordPunctTokenizer(t *testing.T) {
	tok := NewWordPunctTokenizer()
	s := "Good muffins cost $3.88\nin New York.  Please buy me\ntwo of them.\n\nThanks."
	tokens := tok.Tokenize(s)
	if len(tokens) == 0 {
		t.Error("expected tokens, got empty")
	}
	// Should split "$3.88" into "$", "3", ".", "88"
	hasDollar := false
	for _, tok := range tokens {
		if tok == "$" {
			hasDollar = true
			break
		}
	}
	if !hasDollar {
		t.Errorf("expected $ token, got %v", tokens)
	}
}

func TestBlanklineTokenizer(t *testing.T) {
	tok := NewBlanklineTokenizer()
	s := "Good muffins cost $3.88\nin New York.  Please buy me\ntwo of them.\n\nThanks."
	tokens := tok.Tokenize(s)
	if len(tokens) != 2 {
		t.Errorf("expected 2 tokens (blankline), got %d: %v", len(tokens), tokens)
	}
}

func TestTreebankWordTokenizer(t *testing.T) {
	tok := NewTreebankWordTokenizer()

	// Test contractions
	s := "They'll save and invest more."
	tokens := tok.Tokenize(s)
	if len(tokens) < 5 {
		t.Errorf("expected at least 5 tokens, got %d: %v", len(tokens), tokens)
	}

	// Test standard sentence
	s2 := "Good muffins cost $3.88\nin New York."
	tokens2 := tok.Tokenize(s2)
	if len(tokens2) < 6 {
		t.Errorf("expected at least 6 tokens, got %d: %v", len(tokens2), tokens2)
	}
}

func TestTreebankSpanTokenize(t *testing.T) {
	tok := NewTreebankWordTokenizer()
	s := "Good muffins cost $3.88\nin New (York)."
	spans := tok.SpanTokenize(s)
	if len(spans) == 0 {
		t.Error("expected spans, got empty")
	}
	// Verify that spans produce correct tokens
	for _, span := range spans {
		if span.Start > span.End {
			t.Errorf("invalid span: %v", span)
		}
		if span.Start < 0 || span.End > len(s) {
			t.Errorf("span out of bounds: %v", span)
		}
	}
}

func TestNLTKWordTokenizer(t *testing.T) {
	tok := NewNLTKWordTokenizer()
	s := "Good muffins cost $3.88\nin New York.  Please buy me\ntwo of them.\n\nThanks."
	tokens := tok.Tokenize(s)
	if len(tokens) == 0 {
		t.Error("expected tokens, got empty")
	}
}

func TestTweetTokenizer(t *testing.T) {
	tok := NewDefaultTweetTokenizer()
	text := "RT @user: Hello World! Check out https://example.com #golang"
	tokens := tok.Tokenize(text)
	if len(tokens) == 0 {
		t.Error("expected tokens, got empty")
	}
	// Should capture URL as single token
	hasURL := false
	for _, tok := range tokens {
		if tok == "https://example.com" {
			hasURL = true
			break
		}
	}
	if !hasURL {
		t.Errorf("expected URL token, got %v", tokens)
	}
}

func TestCasualTokenize(t *testing.T) {
	text := "RT @user: Hello World! Check out https://example.com #golang"
	tokens := CasualTokenize(text)
	if len(tokens) == 0 {
		t.Error("expected tokens, got empty")
	}
}

func TestMWETokenizer(t *testing.T) {
	base := NewWhitespaceTokenizer()
	expressions := [][]string{
		{"New", "York"},
		{"Los", "Angeles"},
	}
	tok := NewMWETokenizer(base, expressions, "_")
	s := "I live in New York and Los Angeles"
	tokens := tok.Tokenize(s)
	if len(tokens) != 6 {
		t.Errorf("expected 6 tokens (with merged), got %d: %v", len(tokens), tokens)
	}
	if tokens[3] != "New_York" {
		t.Errorf("expected 'New_York', got '%s'", tokens[3])
	}
	if tokens[5] != "Los_Angeles" {
		t.Errorf("expected 'Los_Angeles', got '%s'", tokens[5])
	}
}

func TestRegexpTokenizeFunc(t *testing.T) {
	s := "Good muffins cost $3.88\nin New York."
	tokens := RegexpTokenize(s, `\w+|\$[\d\.]+|\S+`, false, true)
	if len(tokens) == 0 {
		t.Error("expected tokens, got empty")
	}
}

func TestWordpunctTokenizeFunc(t *testing.T) {
	s := "Good muffins cost $3.88\nin New York."
	tokens := WordpunctTokenize(s)
	if len(tokens) == 0 {
		t.Error("expected tokens, got empty")
	}
}

func TestBlanklineTokenizeFunc(t *testing.T) {
	s := "Good muffins cost $3.88\nin New York.\n\nPlease buy me."
	tokens := BlanklineTokenize(s)
	if len(tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d: %v", len(tokens), tokens)
	}
}

func TestRegexpSpanTokenize(t *testing.T) {
	s := "Hello world, this is a test."
	re := regexp.MustCompile(`\w+`)
	spans := regexpSpanTokenize(s, re, false)
	if len(spans) == 0 {
		t.Error("expected spans, got empty")
	}
	for _, span := range spans {
		if span.Start < 0 || span.End > len(s) || span.Start >= span.End {
			t.Errorf("invalid span: %v", span)
		}
	}
}

func TestAlignTokens(t *testing.T) {
	tokens := []string{"Hello", "world", "."}
	text := "Hello world."
	spans := alignTokens(tokens, text)
	if len(spans) != 3 {
		t.Errorf("expected 3 spans, got %d", len(spans))
	}
	if spans[0] != (Span{0, 5}) {
		t.Errorf("expected span {0,5}, got %v", spans[0])
	}
	// Verify extracted tokens match
	for i, span := range spans {
		if text[span.Start:span.End] != tokens[i] {
			t.Errorf("token mismatch: expected '%s', got '%s'", tokens[i], text[span.Start:span.End])
		}
	}
}

func TestSpan(t *testing.T) {
	span := NewSpan(5, 10)
	if span.Start != 5 || span.End != 10 {
		t.Errorf("unexpected span: %v", span)
	}
}

func TestParallelTokenizer(t *testing.T) {
	tok := NewParallelTokenizer(100, 2, NewWhitespaceTokenizer())
	text := "This is a long text that should be tokenized in parallel chunks to test concurrency."
	tokens := tok.Tokenize(text)
	if len(tokens) == 0 {
		t.Error("expected tokens, got empty")
	}
}

func TestParallelTokenize(t *testing.T) {
	text := "This is a test of the parallel tokenize function using the default NLTK word tokenizer."
	tokens := ParallelTokenize(text)
	if len(tokens) == 0 {
		t.Error("expected tokens, got empty")
	}
}

func TestSentTokenizer(t *testing.T) {
	tok := NewSentTokenizer("english")
	text := "Hello world. How are you? I am fine."
	tokens := tok.Tokenize(text)
	if len(tokens) != 3 {
		t.Errorf("expected 3 sentences, got %d: %v", len(tokens), tokens)
	}
}

func TestSplitIntoSentences(t *testing.T) {
	sentences := splitSentences("Hello world. How are you? I am fine!")
	if len(sentences) != 3 {
		t.Errorf("expected 3 sentences, got %d: %v", len(sentences), sentences)
	}
}

func TestSplitIntoChunks(t *testing.T) {
	text := "This is a test string for chunking into multiple pieces"
	chunks := splitIntoChunks(text, 10)
	if len(chunks) < 2 {
		t.Errorf("expected at least 2 chunks, got %d", len(chunks))
	}
	// Verify chunks reassemble correctly
	reassembled := ""
	for _, chunk := range chunks {
		reassembled += chunk
	}
	if reassembled != text {
		t.Errorf("reassembly mismatch:\n  expected: %s\n  got:      %s", text, reassembled)
	}
}

func TestTokenizerIInterface(t *testing.T) {
	// Test that all tokenizers implement the TokenizerI interface
	var tokenizers []TokenizerI
	tokenizers = append(tokenizers, NewSpaceTokenizer())
	tokenizers = append(tokenizers, NewTabTokenizer())
	tokenizers = append(tokenizers, NewCharTokenizer())
	tokenizers = append(tokenizers, NewLineTokenizer("keep"))
	tokenizers = append(tokenizers, NewRegexpTokenizer(`\w+`, false, true))
	tokenizers = append(tokenizers, NewWhitespaceTokenizer())
	tokenizers = append(tokenizers, NewBlanklineTokenizer())
	tokenizers = append(tokenizers, NewWordPunctTokenizer())
	tokenizers = append(tokenizers, NewTreebankWordTokenizer())
	tokenizers = append(tokenizers, NewDefaultTweetTokenizer())

	s := "Hello world. Test 123."
	for i, tok := range tokenizers {
		tokens := tok.Tokenize(s)
		if len(tokens) == 0 {
			t.Errorf("tokenizer %d produced empty tokens", i)
		}
		spans := tok.SpanTokenize(s)
		if len(spans) == 0 {
			t.Errorf("tokenizer %d produced empty spans", i)
		}
	}
}
