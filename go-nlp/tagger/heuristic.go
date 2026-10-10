package tagger

import (
	"strings"
	"unicode"
)

// HeuristicTagger is a small rule-based POS tagger.
type HeuristicTagger struct {
	defaultTag string
}

// NewHeuristicTagger creates a rule-based tagger with a sensible default.
func NewHeuristicTagger() *HeuristicTagger {
	return &HeuristicTagger{defaultTag: "NN"}
}

// Tag applies a set of simple POS heuristics.
func (t *HeuristicTagger) Tag(tokens []string) []TaggedToken {
	out := make([]TaggedToken, len(tokens))
	for i, tok := range tokens {
		out[i] = TaggedToken{Word: tok, Tag: t.tagToken(tok, i, tokens)}
	}
	return out
}

func (t *HeuristicTagger) tagToken(tok string, index int, tokens []string) string {
	lower := strings.ToLower(tok)
	switch {
	case tok == "":
		return t.defaultTag
	case isPunctuation(tok):
		return "."
	case isApostropheS(tok):
		return "POS"
	case lower == "a" || lower == "an" || lower == "the":
		return "DT"
	case lower == "and" || lower == "or" || lower == "but":
		return "CC"
	case lower == "to":
		return "TO"
	case lower == "in" || lower == "on" || lower == "at" || lower == "of" || lower == "for" || lower == "with" || lower == "by":
		return "IN"
	case lower == "is" || lower == "am" || lower == "are" || lower == "was" || lower == "were" || lower == "be" || lower == "been":
		return "VB"
	case lower == "has" || lower == "have" || lower == "had":
		return "VB"
	case strings.HasSuffix(lower, "n't"):
		return "RB"
	case strings.HasPrefix(tok, "@"):
		return "NNP"
	case strings.HasPrefix(tok, "#"):
		return "NN"
	case isNumber(tok):
		return "CD"
	case strings.HasSuffix(lower, "ing"):
		return "VBG"
	case strings.HasSuffix(lower, "ed"):
		return "VBD"
	case strings.HasSuffix(lower, "ly"):
		return "RB"
	case strings.HasSuffix(lower, "ous") || strings.HasSuffix(lower, "ful") || strings.HasSuffix(lower, "able") || strings.HasSuffix(lower, "ive"):
		return "JJ"
	case isTitle(tok):
		if index > 0 && index < len(tokens) {
			prev := strings.ToLower(tokens[index-1])
			if prev == "." || prev == "!" || prev == "?" {
				return "NNP"
			}
		}
		return "NNP"
	default:
		return t.defaultTag
	}
}

func isPunctuation(tok string) bool {
	if tok == "" {
		return false
	}
	for _, r := range tok {
		if !unicode.IsPunct(r) {
			return false
		}
	}
	return true
}

func isApostropheS(tok string) bool {
	return tok == "'s" || tok == "’s"
}

func isNumber(tok string) bool {
	hasDigit := false
	for _, r := range tok {
		if unicode.IsDigit(r) {
			hasDigit = true
			continue
		}
		if r == '.' || r == ',' || r == '-' {
			continue
		}
		return false
	}
	return hasDigit
}

func isTitle(tok string) bool {
	if tok == "" {
		return false
	}
	runes := []rune(tok)
	if !unicode.IsUpper(runes[0]) {
		return false
	}
	for _, r := range runes[1:] {
		if unicode.IsLetter(r) && !unicode.IsLower(r) {
			return false
		}
	}
	return true
}
