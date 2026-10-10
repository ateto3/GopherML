package metrics

func Accuracy(ref, test []string) float64 { if len(ref)==0 || len(ref)!=len(test) { return 0 }; c:=0; for i:=range ref { if ref[i]==test[i] { c++ } }; return float64(c)/float64(len(ref)) }
func Precision(ref, test map[string]struct{}) float64 { if len(test)==0 { return 0 }; c:=0; for k:=range test { if _,ok:=ref[k]; ok { c++ } }; return float64(c)/float64(len(test)) }
func Recall(ref, test map[string]struct{}) float64 { if len(ref)==0 { return 0 }; c:=0; for k:=range ref { if _,ok:=test[k]; ok { c++ } }; return float64(c)/float64(len(ref)) }
func FMeasure(ref, test map[string]struct{}) float64 { p,r:=Precision(ref,test),Recall(ref,test); if p==0||r==0 { return 0 }; return 2*p*r/(p+r) }
