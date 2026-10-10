package text

import "regexp"

type TokenSearcher struct{ Tokens []string }

func NewTokenSearcher(tokens []string) TokenSearcher { return TokenSearcher{Tokens: tokens} }

func (s TokenSearcher) Search(pattern string) []int {
	re := regexp.MustCompile(pattern)
	out := []int{}
	for i, t := range s.Tokens { if re.MatchString(t) { out = append(out, i) } }
	return out
}
