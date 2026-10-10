package translate

import "fmt"

func fmtS(f string, a ...any) string { return fmt.Sprintf(f, a...) }
func join(xs []string, sep string) string { out:=""; for i,x:=range xs { if i>0 { out+=sep }; out+=x }; return out }
func max(a,b int) int { if a>b { return a }; return b }
