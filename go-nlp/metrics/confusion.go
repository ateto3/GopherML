package metrics

import "sort"

type ConfusionMatrix struct{ Values []string; idx map[string]int; M [][]int }

func NewConfusionMatrix(ref, test []string) ConfusionMatrix {
	set:=map[string]struct{}{}; for _,v:=range ref { set[v]=struct{}{} }; for _,v:=range test { set[v]=struct{}{} }
	vals:=make([]string,0,len(set)); for v:=range set { vals=append(vals,v) }; sort.Strings(vals)
	idx:=map[string]int{}; for i,v:=range vals { idx[v]=i }
	m:=make([][]int,len(vals)); for i:=range m { m[i]=make([]int,len(vals)) }
	for i:=range ref { m[idx[ref[i]]][idx[test[i]]]++ }
	return ConfusionMatrix{Values: vals, idx: idx, M: m}
}

func (c ConfusionMatrix) Count(a,b string) int { return c.M[c.idx[a]][c.idx[b]] }
func (c ConfusionMatrix) Precision(v string) float64 { tp:=c.Count(v,v); tot:=0; for _,x:=range c.Values { tot+=c.Count(x,v) }; if tot==0 { return 0 }; return float64(tp)/float64(tot) }
func (c ConfusionMatrix) Recall(v string) float64 { tp:=c.Count(v,v); tot:=0; for _,x:=range c.Values { tot+=c.Count(v,x) }; if tot==0 { return 0 }; return float64(tp)/float64(tot) }
func (c ConfusionMatrix) FMeasure(v string) float64 { p,r:=c.Precision(v),c.Recall(v); if p==0||r==0 { return 0 }; return 2*p*r/(p+r) }
