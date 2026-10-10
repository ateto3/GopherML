package draw

import (
	"fmt"
	"strings"

	"github.com/ateto3/gopherml/go-nlp/tree"
)

func ASCII(t tree.Tree) string { return ascii(t, "", true) }

func ascii(t tree.Tree, pref string, last bool) string {
	out := pref
	if pref != "" {
		if last {
			out += "└─"
		} else {
			out += "├─"
		}
	}
	out += t.Label + "\n"
	childPref := pref
	if pref != "" {
		if last {
			childPref += "  "
		} else {
			childPref += "│ "
		}
	}
	for i, c := range t.Children {
		if tr, ok := c.(tree.Tree); ok {
			out += ascii(tr, childPref, i == len(t.Children)-1)
		} else {
			out += childPref
			if i == len(t.Children)-1 {
				out += "└─"
			} else {
				out += "├─"
			}
			out += fmt.Sprint(c) + "\n"
		}
	}
	return out
}

func DOT(t tree.Tree) string {
	var b strings.Builder
	b.WriteString("digraph G {\n")
	i := 0
	var walk func(tree.Tree, int) int
	walk = func(n tree.Tree, parent int) int {
		me := i
		i++
		b.WriteString(fmt.Sprintf("n%d [label=%q]\n", me, n.Label))
		if parent >= 0 {
			b.WriteString(fmt.Sprintf("n%d -> n%d\n", parent, me))
		}
		for _, c := range n.Children {
			if tr, ok := c.(tree.Tree); ok {
				walk(tr, me)
			} else {
				leaf := i
				i++
				b.WriteString(fmt.Sprintf("n%d [label=%q]\n", leaf, fmt.Sprint(c)))
				b.WriteString(fmt.Sprintf("n%d -> n%d\n", me, leaf))
			}
		}
		return me
	}
	walk(t, -1)
	b.WriteString("}\n")
	return b.String()
}
