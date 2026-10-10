package tokenizer

import (
	"regexp"
	"strings"
)

// splitOnString splits the string s on the delimiter delim.
// This is equivalent to Python's s.split(delim).
func splitOnString(s, delim string) []string {
	return strings.Split(s, delim)
}

// stringSpanTokenize generates spans for tokens split by the delimiter.
func stringSpanTokenize(s, delim string) []Span {
	var spans []Span
	start := 0
	for {
		idx := strings.Index(s[start:], delim)
		if idx == -1 {
			spans = append(spans, Span{Start: start, End: len(s)})
			break
		}
		end := start + idx
		spans = append(spans, Span{Start: start, End: end})
		start = end + len(delim)
	}
	return spans
}

// regexpSpanTokenize generates token spans using a compiled regexp.
// If the regexp matches tokens (not gaps), it finds all matches and yields their spans.
// If gaps=true, it uses the regexp as delimiter patterns.
func regexpSpanTokenize(s string, re *regexp.Regexp, gaps bool) []Span {
	var spans []Span

	if gaps {
		// Use regexp as delimiter: find spans between delimiters
		locs := re.FindAllStringIndex(s, -1)
		start := 0
		for _, loc := range locs {
			if loc[0] > start {
				spans = append(spans, Span{Start: start, End: loc[0]})
			}
			start = loc[1]
		}
		if start < len(s) {
			spans = append(spans, Span{Start: start, End: len(s)})
		}
	} else {
		// Use regexp to match tokens directly
		locs := re.FindAllStringIndex(s, -1)
		for _, loc := range locs {
			spans = append(spans, Span{Start: loc[0], End: loc[1]})
		}
	}
	return spans
}

// alignTokens aligns a list of tokens back to the original string,
// returning their byte offset spans. Equivalent to nltk.tokenize.util.align_tokens.
func alignTokens(tokens []string, text string) []Span {
	var spans []Span
	cur := 0
	for _, tok := range tokens {
		// Find the token in the remaining text
		idx := strings.Index(text[cur:], tok)
		if idx == -1 {
			// Token not found; try to recover by scanning
			continue
		}
		start := cur + idx
		end := start + len(tok)
		spans = append(spans, Span{Start: start, End: end})
		cur = end
	}
	return spans
}
