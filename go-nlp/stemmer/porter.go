package stemmer

import "strings"

// PorterStemmer is a compact suffix stripper.
type PorterStemmer struct{}

func NewPorterStemmer() PorterStemmer { return PorterStemmer{} }

func (PorterStemmer) Stem(word string) string {
	w := strings.ToLower(word)
	for _, suf := range []string{"ingly", "edly", "ing", "edly", "ed", "ly", "ies", "es", "s"} {
		if len(w) > 3 && strings.HasSuffix(w, suf) {
			w = strings.TrimSuffix(w, suf)
			if suf == "ies" {
				w += "y"
			}
			break
		}
	}
	if strings.HasSuffix(w, "at") || strings.HasSuffix(w, "bl") || strings.HasSuffix(w, "iz") {
		return w + "e"
	}
	return w
}
