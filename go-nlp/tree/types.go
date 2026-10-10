package tree

import (
	"fmt"
	"strings"
)

type Tree struct{ Label string; Children []any }
type ImmutableTree struct{ Tree }

func New(label string, children ...any) Tree { return Tree{Label: label, Children: children} }
func (t Tree) String() string { if len(t.Children)==0 { return "("+t.Label+")" }; parts:=make([]string,len(t.Children)); for i,c:=range t.Children { parts[i]=fmt.Sprint(c) }; return "("+t.Label+" "+strings.Join(parts," ")+")" }
func (t Tree) Leaves() []any { out:=[]any{}; for _,c:=range t.Children { if tr,ok:=c.(Tree); ok { out=append(out,tr.Leaves()...) } else { out=append(out,c) } }; return out }
