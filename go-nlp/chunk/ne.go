package chunk

import "unicode"

// NamedEntityChunker heuristically groups capitalized tokens.
type NamedEntityChunker struct{}

// NewNamedEntityChunker creates a basic NE chunker.
func NewNamedEntityChunker() NamedEntityChunker { return NamedEntityChunker{} }

// Parse chunks consecutive title-cased tokens as PERSON.
func (NamedEntityChunker) Parse(tokens []TaggedToken) Tree {
	children := make([]any, 0, len(tokens))
	for i := 0; i < len(tokens); {
		if isNameStart(tokens[i].Word) {
			j := i + 1
			for j < len(tokens) && isNamePart(tokens[j].Word) {
				j++
			}
			ch := make([]any, 0, j-i)
			for k := i; k < j; k++ {
				ch = append(ch, tokens[k])
			}
			children = append(children, Tree{Label: "PERSON", Children: ch})
			i = j
			continue
		}
		children = append(children, tokens[i])
		i++
	}
	return Tree{Label: "S", Children: children}
}

func isNameStart(s string) bool {
	if s == "" {
		return false
	}
	r := []rune(s)
	return unicode.IsUpper(r[0])
}

func isNamePart(s string) bool {
	return isNameStart(s) || s == "." || s == "&"
}
