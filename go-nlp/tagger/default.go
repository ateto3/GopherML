package tagger

// DefaultTagger assigns a fixed tag to all tokens.
type DefaultTagger struct {
	TagValue string
}

// NewDefaultTagger creates a tagger that returns the same tag for every token.
func NewDefaultTagger(tag string) *DefaultTagger {
	if tag == "" {
		tag = "NN"
	}
	return &DefaultTagger{TagValue: tag}
}

// Tag tags every token with the configured default tag.
func (t *DefaultTagger) Tag(tokens []string) []TaggedToken {
	out := make([]TaggedToken, len(tokens))
	for i, tok := range tokens {
		out[i] = TaggedToken{Word: tok, Tag: t.TagValue}
	}
	return out
}
