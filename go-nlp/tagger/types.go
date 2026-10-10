package tagger

// TaggedToken represents a token with its POS tag.
type TaggedToken struct {
	Word string
	Tag  string
}

// Tagger defines the minimal interface for POS taggers.
type Tagger interface {
	Tag(tokens []string) []TaggedToken
}
