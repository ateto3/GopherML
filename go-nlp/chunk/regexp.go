package chunk

import (
	"regexp"
	"sync"
)

// RegexpChunker chunks tagged tokens by matching tag patterns.
type RegexpChunker struct {
	Label   string
	Pattern *regexp.Regexp
}

// NewRegexpChunker creates a regex chunker for a chunk label.
func NewRegexpChunker(label, pattern string) *RegexpChunker {
	return &RegexpChunker{Label: label, Pattern: regexp.MustCompile(pattern)}
}

// Parse groups contiguous tagged tokens whose tags match the pattern.
func (c *RegexpChunker) Parse(tokens []TaggedToken) Tree {
	children := make([]any, 0, len(tokens))
	for _, tok := range tokens {
		if c.Pattern.MatchString(tok.Tag) {
			children = append(children, Tree{Label: c.Label, Children: []any{tok}})
		} else {
			children = append(children, tok)
		}
	}
	return Tree{Label: "S", Children: children}
}

// ParseSents chunks many tagged sentences in parallel.
func ParseSents(chunker ChunkParser, sents [][]TaggedToken) []Tree {
	if len(sents) == 0 {
		return nil
	}
	out := make([]Tree, len(sents))
	var wg sync.WaitGroup
	for i := range sents {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			out[i] = chunker.Parse(sents[i])
		}()
	}
	wg.Wait()
	return out
}
