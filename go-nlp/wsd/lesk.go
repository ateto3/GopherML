package wsd

import "strings"

func Lesk(context []string, word string, senses []Sense) *Sense {
	ctx := map[string]struct{}{}
	for _, w := range context { ctx[strings.ToLower(w)] = struct{}{} }
	var best *Sense; bestN := -1
	for i := range senses {
		n := 0
		for _, g := range senses[i].Gloss { if _, ok := ctx[strings.ToLower(g)]; ok { n++ } }
		if n > bestN { bestN, best = n, &senses[i] }
	}
	return best
}
