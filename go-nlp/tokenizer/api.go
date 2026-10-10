package tokenizer

// TokenizerI defines the interface that all tokenizers must implement.
// It mirrors the nltk.tokenize.api.TokenizerI interface.
type TokenizerI interface {
	// Tokenize returns a tokenized copy of the input string s.
	Tokenize(s string) []string

	// SpanTokenize yields token spans as (start, end) integer offsets,
	// where s[start:end] is the corresponding token.
	SpanTokenize(s string) []Span
}

// Span represents a token span with start and end byte offsets.
type Span struct {
	Start int
	End   int
}

// NewSpan creates a new Span.
func NewSpan(start, end int) Span {
	return Span{Start: start, End: end}
}

// StringTokenizer is a base for tokenizers that split on a fixed delimiter string.
type StringTokenizer struct {
	Delim string
}

// NewStringTokenizer creates a StringTokenizer with the given delimiter.
func NewStringTokenizer(delim string) *StringTokenizer {
	return &StringTokenizer{Delim: delim}
}

// Tokenize splits the string on the delimiter.
func (st *StringTokenizer) Tokenize(s string) []string {
	return splitOnString(s, st.Delim)
}

// SpanTokenize yields spans for tokens split by the delimiter.
func (st *StringTokenizer) SpanTokenize(s string) []Span {
	return stringSpanTokenize(s, st.Delim)
}
