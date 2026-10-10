package tokenizer

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// SpaceTokenizer tokenizes on the space character (" ").
// Equivalent to nltk.tokenize.simple.SpaceTokenizer.
type SpaceTokenizer struct {
	StringTokenizer
}

// NewSpaceTokenizer creates a new SpaceTokenizer.
func NewSpaceTokenizer() *SpaceTokenizer {
	return &SpaceTokenizer{
		StringTokenizer: StringTokenizer{Delim: " "},
	}
}

// TabTokenizer tokenizes on the tab character ("\t").
// Equivalent to nltk.tokenize.simple.TabTokenizer.
type TabTokenizer struct {
	StringTokenizer
}

// NewTabTokenizer creates a new TabTokenizer.
func NewTabTokenizer() *TabTokenizer {
	return &TabTokenizer{
		StringTokenizer: StringTokenizer{Delim: "\t"},
	}
}

// CharTokenizer tokenizes a string into individual characters.
// Equivalent to nltk.tokenize.simple.CharTokenizer.
type CharTokenizer struct{}

// NewCharTokenizer creates a new CharTokenizer.
func NewCharTokenizer() *CharTokenizer {
	return &CharTokenizer{}
}

// Tokenize splits the string into individual runes (characters).
func (ct *CharTokenizer) Tokenize(s string) []string {
	var tokens []string
	for _, r := range s {
		tokens = append(tokens, string(r))
	}
	return tokens
}

// SpanTokenize yields spans for each character.
func (ct *CharTokenizer) SpanTokenize(s string) []Span {
	var spans []Span
	pos := 0
	for _, r := range s {
		size := utf8.RuneLen(r)
		spans = append(spans, Span{Start: pos, End: pos + size})
		pos += size
	}
	return spans
}

// LineTokenizer tokenizes a string into its lines.
// Equivalent to nltk.tokenize.simple.LineTokenizer.
type LineTokenizer struct {
	BlankLines string // "discard", "keep", or "discard-eof"
}

// NewLineTokenizer creates a new LineTokenizer.
// blankLines controls how blank lines are handled: "discard", "keep", or "discard-eof".
func NewLineTokenizer(blankLines string) *LineTokenizer {
	return &LineTokenizer{BlankLines: blankLines}
}

// Tokenize splits the string into lines.
func (lt *LineTokenizer) Tokenize(s string) []string {
	lines := strings.Split(s, "\n")
	switch lt.BlankLines {
	case "discard":
		var filtered []string
		for _, l := range lines {
			if strings.TrimSpace(l) != "" {
				filtered = append(filtered, l)
			}
		}
		return filtered
	case "discard-eof":
		if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
			lines = lines[:len(lines)-1]
		}
		return lines
	case "keep":
		fallthrough
	default:
		return lines
	}
}

// SpanTokenize yields spans for each line.
func (lt *LineTokenizer) SpanTokenize(s string) []Span {
	if lt.BlankLines == "keep" {
		return stringSpanTokenize(s, "\n")
	}
	// For discard modes, use regexp to skip blank lines
	re := regexp.MustCompile(`\n(\s*\n)*`)
	return regexpSpanTokenize(s, re, true)
}
