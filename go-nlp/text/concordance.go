package text

import (
	"fmt"
	"strings"
)

type ConcordanceIndex struct{ Tokens []string; Offsets map[string][]int; Key func(string) string }

func NewConcordanceIndex(tokens []string) ConcordanceIndex {
	idx := ConcordanceIndex{Tokens: tokens, Offsets: map[string][]int{}, Key: strings.ToLower}
	for i, t := range tokens { idx.Offsets[idx.Key(t)] = append(idx.Offsets[idx.Key(t)], i) }
	return idx
}

func (c ConcordanceIndex) Find(word string, width, lines int) []ConcordanceLine {
	pos := c.Offsets[c.Key(word)]
	if len(pos) == 0 { return nil }
	if lines > len(pos) { lines = len(pos) }
	out := make([]ConcordanceLine, 0, lines)
	for _, i := range pos[:lines] {
		l := max(0, i-width/4)
		r := min(len(c.Tokens), i+1+width/4)
		out = append(out, ConcordanceLine{
			Left: strings.Join(c.Tokens[l:i], " "),
			Query: c.Tokens[i],
			Right: strings.Join(c.Tokens[i+1:r], " "),
			Offset: i,
		})
	}
	return out
}

func (c ConcordanceIndex) String() string { return fmt.Sprintf("<ConcordanceIndex %d>", len(c.Tokens)) }
