package tokenizer

import (
	"regexp"
	"strings"
)

// TreebankWordTokenizer tokenizes text as in the Penn Treebank.
// Ported from nltk.tokenize.treebank.TreebankWordTokenizer.
type TreebankWordTokenizer struct {
	contractions *MacIntyreContractions
}

// NewTreebankWordTokenizer creates a new TreebankWordTokenizer.
func NewTreebankWordTokenizer() *TreebankWordTokenizer {
	return &TreebankWordTokenizer{
		contractions: NewMacIntyreContractions(),
	}
}

// Tokenize returns a tokenized copy of the text as a slice of strings.
// Implements the TokenizerI interface.
func (tw *TreebankWordTokenizer) Tokenize(text string) []string {
	return tw.TokenizeTokens(text)
}

// TokenizeToString returns a tokenized copy of the text using Penn Treebank conventions
// as a single space-separated string.
func (tw *TreebankWordTokenizer) TokenizeToString(text string) string {
	// Starting quotes
	text = startingQuotesRE.ReplaceAllString(text, " `` ")

	// Punctuation
	text = punctuationRE1.ReplaceAllString(text, "${1} ${2}")
	text = punctuationRE2.ReplaceAllString(text, "${1} ")
	text = ellipsisRE.ReplaceAllString(text, " ... ")
	text = punctSymbolsRE.ReplaceAllString(text, " ${1} ")
	text = finalPeriodRE.ReplaceAllString(text, "${1} ${2}${3} ")
	text = questionExclaimRE.ReplaceAllString(text, " ${1} ")
	text = singleQuoteRE.ReplaceAllString(text, "${1} ' ")

	// Parentheses and brackets
	text = parensBracketsRE.ReplaceAllString(text, " ${1} ")

	// Double dashes
	text = doubleDashRE.ReplaceAllString(text, " -- ")

	// Add extra space
	text = " " + text + " "

	// Ending quotes
	text = endingQuotesRE1.ReplaceAllString(text, " '' ")
	text = endingQuotesRE2.ReplaceAllString(text, " '' ")
	text = endingQuotesRE3.ReplaceAllString(text, "${1} ${2} ")
	text = endingQuotesRE4.ReplaceAllString(text, "${1} ${2} ")

	// Contractions (2-part)
	for _, pat := range tw.contractions.CONTRACTIONS2 {
		re := regexp.MustCompile(pat)
		text = re.ReplaceAllString(text, " ${1} ${2} ")
	}
	// Contractions (3-part)
	for _, pat := range tw.contractions.CONTRACTIONS3 {
		re := regexp.MustCompile(pat)
		text = re.ReplaceAllString(text, " ${1} ${2} ")
	}

	return strings.TrimSpace(text)
}

// TokenizeTokens returns the tokenized text as a slice of strings.
func (tw *TreebankWordTokenizer) TokenizeTokens(text string) []string {
	result := tw.TokenizeToString(text)
	if result == "" {
		return nil
	}
	return strings.Fields(result)
}

// SpanTokenize returns the spans of the tokens in the text.
func (tw *TreebankWordTokenizer) SpanTokenize(text string) []Span {
	tokens := tw.TokenizeTokens(text)
	// Handle quote conversion for span matching
	if strings.Contains(text, "\"") || strings.Contains(text, "''") {
		tokens = convertQuoteSpans(tokens, text)
	}
	return alignTokens(tokens, text)
}

// convertQuoteSpans handles quote conversion for accurate span matching.
func convertQuoteSpans(tokens []string, text string) []string {
	re := regexp.MustCompile("``|'{2}|\"")
	matched := re.FindAllString(text, -1)
	var result []string
	quoteIdx := 0
	for _, tok := range tokens {
		if tok == "\"" || tok == "``" || tok == "''" {
			if quoteIdx < len(matched) {
				result = append(result, matched[quoteIdx])
				quoteIdx++
			} else {
				result = append(result, tok)
			}
		} else {
			result = append(result, tok)
		}
	}
	return result
}

// TreebankWordDetokenizer reverses the TreebankWordTokenizer.
// Ported from nltk.tokenize.treebank.TreebankWordDetokenizer.
type TreebankWordDetokenizer struct {
	contractions *MacIntyreContractions
}

// NewTreebankWordDetokenizer creates a new TreebankWordDetokenizer.
func NewTreebankWordDetokenizer() *TreebankWordDetokenizer {
	return &TreebankWordDetokenizer{
		contractions: NewMacIntyreContractions(),
	}
}

// Detokenize joins tokens back into text, reversing the tokenization rules.
func (twd *TreebankWordDetokenizer) Detokenize(tokens []string) string {
	text := " " + strings.Join(tokens, " ") + " "

	// Reverse contractions (3-part first)
	for _, pat := range twd.contractions.CONTRACTIONS3 {
		re := regexp.MustCompile(patReverser(pat))
		text = re.ReplaceAllString(text, "${1}${2}")
	}
	for _, pat := range twd.contractions.CONTRACTIONS2 {
		re := regexp.MustCompile(patReverser(pat))
		text = re.ReplaceAllString(text, "${1}${2}")
	}

	// Reverse ending quotes
	text = detokEndingQuotesRE1.ReplaceAllString(text, "${1}${2} ")
	text = detokEndingQuotesRE2.ReplaceAllString(text, "${1}${2} ")
	text = detokEndingQuotesRE3.ReplaceAllString(text, "${1}${2}")
	text = detokEndingQuotesRE4.ReplaceAllString(text, "${1}${2}")
	text = detokEndingQuotesRE5.ReplaceAllString(text, "${1}${2}")
	text = detokEndingQuotesRE6.ReplaceAllString(text, "\"")
	text = detokEndingQuotesRE7.ReplaceAllString(text, `${1}${2}"`)

	// Reverse double dashes
	text = detokDoubleDashRE.ReplaceAllString(text, "--")

	// Reverse parentheses/brackets padding
	text = detokParensBracketsRE1.ReplaceAllString(text, "${1}")
	text = detokParensBracketsRE2.ReplaceAllString(text, "${1}")
	text = detokParensBracketsRE3.ReplaceAllString(text, "${1}${2}")

	// Reverse punctuation
	text = detokPunctRE1.ReplaceAllString(text, "${1}' ")
	text = detokPunctRE2.ReplaceAllString(text, "${1}")
	text = detokPunctRE3.ReplaceAllString(text, "${1}${2}${3}")
	text = detokPunctRE4.ReplaceAllString(text, "${1}")
	text = detokPunctRE5.ReplaceAllString(text, "${1}")
	text = detokPunctRE6.ReplaceAllString(text, "...")
	text = detokPunctRE7.ReplaceAllString(text, "${1}")

	// Reverse starting quotes
	text = detokStartQuotesRE1.ReplaceAllString(text, "${1}``")
	text = detokStartQuotesRE2.ReplaceAllString(text, "${1}")
	text = detokStartQuotesRE3.ReplaceAllString(text, "\"")

	return strings.TrimSpace(text)
}

// patReverser converts a tokenization contraction pattern to its detokenization form.
func patReverser(pat string) string {
	// Replace capturing groups with space-separated versions
	re := regexp.MustCompile(`\\(\w+)`)
	return re.ReplaceAllString(pat, " $1 ")
}

// Precompiled regex patterns for Treebank tokenizer
var (
	startingQuotesRE  = regexp.MustCompile(`(^")|(\x60\x60)|([ \(\[{<])("|'{2})`)
	punctuationRE1    = regexp.MustCompile(`([:,])([^\d])`)
	punctuationRE2    = regexp.MustCompile(`([:,])$`)
	ellipsisRE        = regexp.MustCompile(`\.\.\.`)
	punctSymbolsRE    = regexp.MustCompile(`[;@#$%&]`)
	finalPeriodRE     = regexp.MustCompile(`([^\.])(\.)([\]\)}>"']*)\s*$`)
	questionExclaimRE = regexp.MustCompile(`[?!]`)
	singleQuoteRE     = regexp.MustCompile(`([^'])' `)
	parensBracketsRE  = regexp.MustCompile(`[\]\[\(\)\{\}\<\>]`)
	doubleDashRE      = regexp.MustCompile(`--`)

	endingQuotesRE1 = regexp.MustCompile(`''`)
	endingQuotesRE2 = regexp.MustCompile(`"`)
	endingQuotesRE3 = regexp.MustCompile(`([^' ])('[sS]|'[mM]|'[dD]|') `)
	endingQuotesRE4 = regexp.MustCompile(`([^' ])('ll|'LL|'re|'RE|'ve|'VE|n't|N'T) `)

	// Detokenization patterns
	detokEndingQuotesRE1 = regexp.MustCompile(`([^' ])\s('ll|'LL|'re|'RE|'ve|'VE|n't|N'T) `)
	detokEndingQuotesRE2 = regexp.MustCompile(`([^' ])\s('[sS]|'[mM]|'[dD]|') `)
	detokEndingQuotesRE3 = regexp.MustCompile(`([^'\s])\s('')`)
	detokEndingQuotesRE4 = regexp.MustCompile(`([,.;:!?'])\s+("|'')`)
	detokEndingQuotesRE5 = regexp.MustCompile(`('')\s([.,:)\]>};%])`)
	detokEndingQuotesRE6 = regexp.MustCompile(`''`)
	detokEndingQuotesRE7 = regexp.MustCompile(`([,.;:!?])"(')`)

	detokDoubleDashRE = regexp.MustCompile(` -- `)

	detokParensBracketsRE1 = regexp.MustCompile(`([\[\(\{\<])\s`)
	detokParensBracketsRE2 = regexp.MustCompile(`\s([\]\)\}\>])`)
	detokParensBracketsRE3 = regexp.MustCompile(`([\]\)\}\>])\s([:;,.])`)

	detokPunctRE1 = regexp.MustCompile(`([^'])\s'\s`)
	detokPunctRE2 = regexp.MustCompile(`\s([?!])`)
	detokPunctRE3 = regexp.MustCompile(`([^\.])\s(\.)([\]\)}>"']*)\s*$`)
	detokPunctRE4 = regexp.MustCompile(`([#$])\s`)
	detokPunctRE5 = regexp.MustCompile(`\s([;%])`)
	detokPunctRE6 = regexp.MustCompile(`\s\.\.\.\s`)
	detokPunctRE7 = regexp.MustCompile(`\s([:,])`)

	detokStartQuotesRE1 = regexp.MustCompile(`([ (\[{<])\s\x60\x60`)
	detokStartQuotesRE2 = regexp.MustCompile(`(\x60\x60)\s`)
	detokStartQuotesRE3 = regexp.MustCompile(`\x60\x60`)
)
