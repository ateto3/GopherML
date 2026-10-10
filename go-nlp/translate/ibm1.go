package translate

import "math"

type IBMModel1 struct{ Table map[string]map[string]float64 }

func NewIBMModel1(corpus []AlignedSent, iters int) *IBMModel1 {
	m := &IBMModel1{Table: map[string]map[string]float64{}}
	srcV := map[string]struct{}{}
	trgV := map[string]struct{}{}
	for _, s := range corpus { for _, w := range s.Mots { srcV[w]=struct{}{} }; for _, w := range s.Words { trgV[w]=struct{}{} } }
	init := 1 / float64(max(1,len(srcV)))
	for t := range trgV { m.Table[t] = map[string]float64{}; for s := range srcV { m.Table[t][s]=init }; m.Table[t][""]=init }
	for i:=0;i<iters;i++ { m.train(corpus) }
	for i := range corpus { m.align(&corpus[i]) }
	return m
}

func (m *IBMModel1) train(corpus []AlignedSent) {
	count := map[string]map[string]float64{}
	tot := map[string]float64{}
	for _, sent := range corpus {
		src := append([]string{""}, sent.Mots...)
		for _, t := range sent.Words {
			z := 0.0
			for _, s := range src { z += m.Table[t][s] }
			for _, s := range src {
				c := m.Table[t][s] / z
				if count[t] == nil { count[t] = map[string]float64{} }
				count[t][s] += c; tot[s] += c
			}
		}
	}
	for t, row := range count { for s, c := range row { m.Table[t][s] = c / math.Max(tot[s], 1e-12) } }
}

func (m *IBMModel1) align(sent *AlignedSent) {
	al := make(Alignment, 0, len(sent.Words))
	for j, t := range sent.Words {
		best, bi := m.Table[t][""], -1
		for i, s := range sent.Mots { if p:=m.Table[t][s]; p >= best { best, bi = p, i } }
		al = append(al, [2]int{j, bi})
	}
	sent.Alignment = al
}
