package tree

import "strings"

func FromString(s string) (Tree, error) {
	toks := strings.Fields(strings.NewReplacer("(", " ( ", ")", " ) ").Replace(s))
	var stack []Tree
	for i := 0; i < len(toks); i++ {
		switch toks[i] {
		case "(":
			if i+1 >= len(toks) { return Tree{}, nil }
			stack = append(stack, Tree{Label: toks[i+1]}); i++
		case ")":
			if len(stack) == 0 { continue }
			t := stack[len(stack)-1]; stack = stack[:len(stack)-1]
			if len(stack) == 0 { return t, nil }
			stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, t)
		default:
			if len(stack) == 0 { continue }
			stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, toks[i])
		}
	}
	if len(stack) > 0 { return stack[0], nil }
	return Tree{}, nil
}
