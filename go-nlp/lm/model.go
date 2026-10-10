package lm

import "math/rand"

type MLE struct{ Counts *NgramCounter }

func NewMLE() *MLE { return &MLE{Counts: NewCounter()} }

func (m *MLE) Fit(text []Sentence) { m.Counts.Update(text) }

func (m *MLE) Score(word string, context ...string) float64 {
	if len(context) == 0 {
		n := 0
		for _, c := range m.Counts.Unigrams { n += c }
		if n == 0 { return 0 }
		return float64(m.Counts.Unigrams[word]) / float64(n)
	}
	ctx := context[len(context)-1]
	row := m.Counts.Bigrams[ctx]
	if row == nil { return 0 }
	tot := 0
	for _, c := range row { tot += c }
	if tot == 0 { return 0 }
	return float64(row[word]) / float64(tot)
}

func (m *MLE) Generate(seed string, n int) []string {
	out := make([]string, 0, n)
	cur := seed
	for i := 0; i < n; i++ {
		out = append(out, cur)
		row := m.Counts.Bigrams[cur]
		if len(row) == 0 { break }
		var best string
		bestN := -1
		for w, c := range row { if c > bestN { best, bestN = w, c } }
		cur = best
	}
	return out
}

type Laplace struct{ *MLE }

func NewLaplace() *Laplace { return &Laplace{NewMLE()} }

func (l *Laplace) Score(word string, context ...string) float64 {
	if len(context) == 0 {
		n := 0
		for _, c := range l.Counts.Unigrams { n += c }
		v := len(l.Counts.Unigrams)
		return float64(l.Counts.Unigrams[word]+1) / float64(n+v)
	}
	ctx := context[len(context)-1]
	row := l.Counts.Bigrams[ctx]
	tot := 0
	for _, c := range row { tot += c }
	v := len(l.Counts.Unigrams)
	return float64(row[word]+1) / float64(tot+v)
}

func RandomChoice(weights map[string]float64) string {
	var best string
	bestV := -1.0
	for k, v := range weights { if v > bestV { best, bestV = k, v } }
	return best
}

func ShuffleSentences(s []Sentence) []Sentence { return append([]Sentence(nil), s...) }
var _ = rand.Int
