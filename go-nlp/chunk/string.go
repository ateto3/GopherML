package chunk

import "strings"

// ChunkString is a lightweight IOB-like encoding.
type ChunkString struct {
	Tags    []string
	Chunks  []bool
	Root    string
	Pattern string
}

// NewChunkString creates an unchunked representation.
func NewChunkString(tokens []TaggedToken, root string) *ChunkString {
	tags := make([]string, len(tokens))
	for i, tok := range tokens {
		tags[i] = tok.Tag
	}
	return &ChunkString{Tags: tags, Chunks: make([]bool, len(tokens)), Root: root}
}

// String renders the chunk string in a compact tag/braces form.
func (cs *ChunkString) String() string {
	var b strings.Builder
	in := false
	for i, tag := range cs.Tags {
		if cs.Chunks[i] && !in {
			b.WriteByte('{')
			in = true
		}
		if !cs.Chunks[i] && in {
			b.WriteByte('}')
			in = false
		}
		b.WriteByte('<')
		b.WriteString(tag)
		b.WriteByte('>')
	}
	if in {
		b.WriteByte('}')
	}
	return b.String()
}

func (cs *ChunkString) mark(i, j int, value bool) {
	if i < 0 {
		i = 0
	}
	if j > len(cs.Chunks) {
		j = len(cs.Chunks)
	}
	for k := i; k < j; k++ {
		cs.Chunks[k] = value
	}
}
