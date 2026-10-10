package translate

import "sort"

type Alignment [][2]int
type AlignedSent struct{ Words, Mots []string; Alignment Alignment }

func (a Alignment) Invert() Alignment { out:=make(Alignment,len(a)); for i,p:=range a { out[i]=[2]int{p[1],p[0]} }; return out }
func (a Alignment) String() string { out:=make([]string,len(a)); for i,p:=range a { out[i]=fmtS("%d-%d", p[0], p[1]) }; sort.Strings(out); return join(out," ") }
func (s *AlignedSent) Invert() { s.Words, s.Mots = s.Mots, s.Words; s.Alignment = s.Alignment.Invert() }
