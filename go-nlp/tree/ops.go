package tree

func (t Tree) Height() int { h:=1; for _,c:=range t.Children { if tr,ok:=c.(Tree); ok { if x:=1+tr.Height(); x>h { h=x } } }; return h }
func (t Tree) Map(fn func(string) string) Tree { out:=Tree{Label: fn(t.Label)}; for _,c:=range t.Children { if tr,ok:=c.(Tree); ok { out.Children=append(out.Children, tr.Map(fn)) } else if s,ok:=c.(string); ok { out.Children=append(out.Children, fn(s)) } else { out.Children=append(out.Children, c) } }; return out }
