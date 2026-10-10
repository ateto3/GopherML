package featstruct

import (
	"fmt"
	"sort"
	"strings"
)

type FeatStruct map[string]any

func New(m map[string]any) FeatStruct { out:=FeatStruct{}; for k,v:=range m { out[k]=v }; return out }
func (f FeatStruct) Get(path ...string) any { cur:=any(f); for _,p:=range path { m,ok:=cur.(FeatStruct); if !ok { return nil }; cur=m[p] }; return cur }
func (f FeatStruct) Set(v any, path ...string) { if len(path)==0 { return }; cur:=f; for _,p:=range path[:len(path)-1] { nxt,ok:=cur[p].(FeatStruct); if !ok { nxt=FeatStruct{}; cur[p]=nxt }; cur=nxt }; cur[path[len(path)-1]]=v }
func (f FeatStruct) String() string { ks:=make([]string,0,len(f)); for k:=range f { ks=append(ks,k) }; sort.Strings(ks); parts:=make([]string,0,len(ks)); for _,k:=range ks { parts=append(parts, fmt.Sprintf("%s=%v", k, f[k])) }; return "["+strings.Join(parts," ")+"]" }
