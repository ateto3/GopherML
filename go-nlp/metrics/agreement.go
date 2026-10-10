package metrics

func CohenKappa(a, b []string) float64 {
	if len(a)==0 || len(a)!=len(b) { return 0 }
	cm:=NewConfusionMatrix(a,b); n:=float64(len(a))
	po:=0.0; pe:=0.0
	for _,v:=range cm.Values { po += float64(cm.Count(v,v)) / n }
	for _,v:=range cm.Values {
		r:=0; t:=0
		for _,x:=range cm.Values { r+=cm.Count(v,x); t+=cm.Count(x,v) }
		pe += (float64(r)/n)*(float64(t)/n)
	}
	if pe==1 { return 1 }
	return (po-pe)/(1-pe)
}
