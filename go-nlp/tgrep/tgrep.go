package tgrep

import (
	"strings"

	"github.com/ateto3/gopherml/go-nlp/tree"
)

func Find(t tree.Tree, q string) []tree.Tree {
	out := []tree.Tree{}
	walk(t, q, &out)
	return out
}

func walk(t tree.Tree, q string, out *[]tree.Tree) {
	if match(t, q) {
		*out = append(*out, t)
	}
	for _, c := range t.Children {
		if tr, ok := c.(tree.Tree); ok {
			walk(tr, q, out)
		}
	}
}

func match(t tree.Tree, q string) bool {
	if strings.HasPrefix(q, "label=") {
		return t.Label == strings.TrimPrefix(q, "label=")
	}
	return t.Label == q
}
