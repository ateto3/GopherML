package chunk

// TaggedToken is a word/tag pair.
type TaggedToken struct {
	Word string
	Tag  string
}

// Tree is a shallow chunk tree.
type Tree struct {
	Label    string
	Children []any
}

// Leaves returns all token leaves in order.
func (t Tree) Leaves() []TaggedToken {
	var out []TaggedToken
	for _, child := range t.Children {
		switch v := child.(type) {
		case TaggedToken:
			out = append(out, v)
		case Tree:
			out = append(out, v.Leaves()...)
		}
	}
	return out
}

// ChunkParser defines the chunking interface.
type ChunkParser interface {
	Parse(tokens []TaggedToken) Tree
}
