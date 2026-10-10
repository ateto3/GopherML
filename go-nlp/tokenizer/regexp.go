package tokenizer

import (
	"regexp"
)

// RegexpTokenizer splits a string using a regular expression.
// The pattern can match either tokens or the gaps (delimiters) between tokens.
// Equivalent to nltk.tokenize.regexp.RegexpTokenizer.
type RegexpTokenizer struct {
	pattern      string
	regexp       *regexp.Regexp
	gaps         bool
	discardEmpty bool
}

// NewRegexpTokenizer creates a new RegexpTokenizer.
// If gaps is true, the pattern matches separators between tokens.
// If discardEmpty is true, empty tokens are discarded (only relevant when gaps=true).
// flags are standard regexp flags (see regexp package).
func NewRegexpTokenizer(pattern string, gaps, discardEmpty bool) *RegexpTokenizer {
	re := regexp.MustCompile(pattern)
	return &RegexpTokenizer{
		pattern:      pattern,
		regexp:       re,
		gaps:         gaps,
		discardEmpty: discardEmpty,
	}
}

// NewRegexpTokenizerFlags creates a new RegexpTokenizer with custom regexp flags.
// The flag parameter is a string of flags like "(?i)" for case-insensitive, "(?m)" for multiline, etc.
// These can be embedded in the pattern itself.
func NewRegexpTokenizerFlags(pattern string, gaps, discardEmpty bool, flags string) *RegexpTokenizer {
	fullPattern := flags + pattern
	re := regexp.MustCompile(fullPattern)
	return &RegexpTokenizer{
		pattern:      fullPattern,
		regexp:       re,
		gaps:         gaps,
		discardEmpty: discardEmpty,
	}
}

// Tokenize returns a tokenized copy of the text.
func (rt *RegexpTokenizer) Tokenize(text string) []string {
	if rt.gaps {
		// Split on the pattern (delimiter)
		parts := rt.regexp.Split(text, -1)
		if rt.discardEmpty {
			var result []string
			for _, p := range parts {
				if p != "" {
					result = append(result, p)
				}
			}
			return result
		}
		return parts
	}
	// Find all matches (tokens)
	return rt.regexp.FindAllString(text, -1)
}

// SpanTokenize yields token spans.
func (rt *RegexpTokenizer) SpanTokenize(text string) []Span {
	return regexpSpanTokenize(text, rt.regexp, rt.gaps)
}

// WhitespaceTokenizer tokenizes on whitespace.
// Equivalent to nltk.tokenize.WhitespaceTokenizer.
type WhitespaceTokenizer struct {
	RegexpTokenizer
}

// NewWhitespaceTokenizer creates a new WhitespaceTokenizer.
func NewWhitespaceTokenizer() *WhitespaceTokenizer {
	return &WhitespaceTokenizer{
		RegexpTokenizer: RegexpTokenizer{
			pattern:      `\s+`,
			regexp:       regexp.MustCompile(`\s+`),
			gaps:         true,
			discardEmpty: true,
		},
	}
}

// BlanklineTokenizer tokenizes on blank lines.
// Equivalent to nltk.tokenize.BlanklineTokenizer.
type BlanklineTokenizer struct {
	RegexpTokenizer
}

// NewBlanklineTokenizer creates a new BlanklineTokenizer.
func NewBlanklineTokenizer() *BlanklineTokenizer {
	return &BlanklineTokenizer{
		RegexpTokenizer: RegexpTokenizer{
			pattern:      `\s*\n\s*\n\s*`,
			regexp:       regexp.MustCompile(`\s*\n\s*\n\s*`),
			gaps:         true,
			discardEmpty: true,
		},
	}
}

// WordPunctTokenizer tokenizes into alphabetic and non-alphabetic characters.
// Equivalent to nltk.tokenize.WordPunctTokenizer.
type WordPunctTokenizer struct {
	RegexpTokenizer
}

// NewWordPunctTokenizer creates a new WordPunctTokenizer.
func NewWordPunctTokenizer() *WordPunctTokenizer {
	return &WordPunctTokenizer{
		RegexpTokenizer: RegexpTokenizer{
			pattern:      `\w+|[^\w\s]+`,
			regexp:       regexp.MustCompile(`\w+|[^\w\s]+`),
			gaps:         false,
			discardEmpty: true,
		},
	}
}

// Convenience functions

// RegexpTokenize tokenizes text using the given regexp pattern.
// Equivalent to nltk.tokenize.regexp_tokenize().
func RegexpTokenize(text, pattern string, gaps, discardEmpty bool) []string {
	t := NewRegexpTokenizer(pattern, gaps, discardEmpty)
	return t.Tokenize(text)
}

// WordpunctTokenize tokenizes text into alphabetic and non-alphabetic tokens.
// Equivalent to nltk.tokenize.wordpunct_tokenize().
func WordpunctTokenize(text string) []string {
	return NewWordPunctTokenizer().Tokenize(text)
}

// BlanklineTokenize tokenizes text on blank lines.
// Equivalent to nltk.tokenize.blankline_tokenize().
func BlanklineTokenize(text string) []string {
	return NewBlanklineTokenizer().Tokenize(text)
}
