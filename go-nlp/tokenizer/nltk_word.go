package tokenizer

// NLTKWordTokenizer is an improved TreebankWordTokenizer that handles
// additional edge cases. It is the default word tokenizer in NLTK.
// Ported from nltk.tokenize.destructive.NLTKWordTokenizer.
type NLTKWordTokenizer struct {
	contractions *MacIntyreContractions
}

// NewNLTKWordTokenizer creates a new NLTKWordTokenizer.
func NewNLTKWordTokenizer() *NLTKWordTokenizer {
	return &NLTKWordTokenizer{
		contractions: NewMacIntyreContractions(),
	}
}

// Tokenize returns a tokenized copy of the text.
func (nwt *NLTKWordTokenizer) Tokenize(text string) []string {
	// Apply the Treebank tokenizer first
	treebank := NewTreebankWordTokenizer()
	return treebank.TokenizeTokens(text)
}

// SpanTokenize returns the spans of the tokens in the text.
func (nwt *NLTKWordTokenizer) SpanTokenize(text string) []Span {
	tokens := nwt.Tokenize(text)
	return alignTokens(tokens, text)
}
