package stemmer

import "regexp"

// RegexpStemmer strips suffixes matching a regex.
type RegexpStemmer struct {
	re  *regexp.Regexp
	min int
}

func NewRegexpStemmer(pattern string, min int) *RegexpStemmer {
	return &RegexpStemmer{re: regexp.MustCompile(pattern), min: min}
}

func (s *RegexpStemmer) Stem(word string) string {
	if len(word) < s.min {
		return word
	}
	return s.re.ReplaceAllString(word, "")
}
