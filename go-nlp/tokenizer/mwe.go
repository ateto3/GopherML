package tokenizer

import (
	"strings"
)

// MWETokenizer tokenizes text into multi-word expressions (MWEs).
// Equivalent to nltk.tokenize.mwe.MWETokenizer.
// An MWE tokenizer wraps a base tokenizer and merges certain tokens
// that are part of predefined multi-word expressions.
type MWETokenizer struct {
	baseTokenizer TokenizerI
	expressions   [][]string
	separator     string
}

// NewMWETokenizer creates a new MWETokenizer.
// baseTokenizer: the underlying tokenizer to use.
// expressions: list of multi-word expressions (each as a slice of tokens).
// separator: the separator to use when joining the tokens (e.g., " " or "_").
func NewMWETokenizer(baseTokenizer TokenizerI, expressions [][]string, separator string) *MWETokenizer {
	return &MWETokenizer{
		baseTokenizer: baseTokenizer,
		expressions:   expressions,
		separator:     separator,
	}
}

// Tokenize tokenizes the text and merges multi-word expressions.
func (mwe *MWETokenizer) Tokenize(text string) []string {
	tokens := mwe.baseTokenizer.Tokenize(text)
	return mwe.mergeExpressions(tokens)
}

// SpanTokenize returns token spans. For MWE tokenization, this
// delegates to the base tokenizer since span mapping with merged tokens
// is complex.
func (mwe *MWETokenizer) SpanTokenize(text string) []Span {
	// For MWEs, span tokenization is simplified
	return mwe.baseTokenizer.SpanTokenize(text)
}

// mergeExpressions merges consecutive tokens that form a multi-word expression.
func (mwe *MWETokenizer) mergeExpressions(tokens []string) []string {
	if len(mwe.expressions) == 0 {
		return tokens
	}

	// Build a trie for fast prefix matching of expressions
	trie := buildExpressionTrie(mwe.expressions)

	var result []string
	i := 0
	for i < len(tokens) {
		matched, matchLen := trieMatch(trie, tokens, i)
		if matched {
			// Merge the matched tokens
			merged := strings.Join(tokens[i:i+matchLen], mwe.separator)
			result = append(result, merged)
			i += matchLen
		} else {
			result = append(result, tokens[i])
			i++
		}
	}
	return result
}

// exprTrieNode is a node in the expression trie.
type exprTrieNode struct {
	children map[string]*exprTrieNode
	isEnd    bool
}

// buildExpressionTrie builds a trie from a list of expressions.
func buildExpressionTrie(expressions [][]string) *exprTrieNode {
	root := &exprTrieNode{children: make(map[string]*exprTrieNode)}
	for _, expr := range expressions {
		node := root
		for _, token := range expr {
			if _, ok := node.children[token]; !ok {
				node.children[token] = &exprTrieNode{children: make(map[string]*exprTrieNode)}
			}
			node = node.children[token]
		}
		node.isEnd = true
	}
	return root
}

// trieMatch checks if tokens starting at position i match any expression.
func trieMatch(root *exprTrieNode, tokens []string, start int) (bool, int) {
	node := root
	bestLen := 0
	for i := start; i < len(tokens); i++ {
		if next, ok := node.children[tokens[i]]; ok {
			node = next
			if node.isEnd {
				bestLen = i - start + 1
			}
		} else {
			break
		}
	}
	return bestLen > 0, bestLen
}
