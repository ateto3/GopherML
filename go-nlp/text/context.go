package text

import "strings"

type ContextIndex struct{ Tokens []string; WordCtx map[string]map[string]int; CtxWord map[string]map[string]int; Key func(string) string }

func NewContextIndex(tokens []string) ContextIndex {
	ci := ContextIndex{Tokens: tokens, WordCtx: map[string]map[string]int{}, CtxWord: map[string]map[string]int{}, Key: strings.ToLower}
	for i, t := range tokens {
		l := "*START*"; r := "*END*"
		if i > 0 { l = ci.Key(tokens[i-1]) }
		if i+1 < len(tokens) { r = ci.Key(tokens[i+1]) }
		ctx := l + "|" + r
		w := ci.Key(t)
		if ci.WordCtx[w] == nil { ci.WordCtx[w] = map[string]int{} }
		if ci.CtxWord[ctx] == nil { ci.CtxWord[ctx] = map[string]int{} }
		ci.WordCtx[w][ctx]++
		ci.CtxWord[ctx][w]++
	}
	return ci
}

func (c ContextIndex) SimilarWords(word string, n int) []string {
	if n <= 0 { n = 20 }
	scores := map[string]int{}
	for ctx := range c.WordCtx[c.Key(word)] {
		for w, c2 := range c.CtxWord[ctx] { scores[w] += c2 }
	}
	out := topN(scores, n)
	return out
}
