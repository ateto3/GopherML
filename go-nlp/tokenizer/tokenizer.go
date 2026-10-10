// Package tokenizer provides tokenization algorithms for Natural Language Processing.
//
// This package is a port of NLTK's tokenize module to Go, with additional
// support for concurrent/parallel tokenization using goroutines and worker pools.
//
// Tokenizers divide strings into lists of substrings (tokens).
// For example, tokenizers can be used to find the words and punctuation in a string:
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/ateto3/gopherml/go-nlp/tokenizer"
//	)
//
//	func main() {
//	    s := "Good muffins cost $3.88\nin New York.  Please buy me\ntwo of them.\n\nThanks."
//	    tokens := tokenizer.TreebankWordTokenizer(NewTreebankWordTokenizer()).TokenizeTokens(s)
//	    fmt.Printf("%v\n", tokens)
//	    // ['Good', 'muffins', 'cost', '$', '3.88', 'in', 'New', 'York.', 'Please', 'buy', 'me', 'two', 'of', 'them.', 'Thanks', '.']
//	}
//
// Available tokenizers:
//   - TreebankWordTokenizer: Penn Treebank tokenizer (standard)
//   - NLTKWordTokenizer: Improved Treebank tokenizer (default NLTK word tokenizer)
//   - TweetTokenizer: Social media / Twitter tokenizer
//   - RegexpTokenizer: Custom regular expression tokenizer
//   - WhitespaceTokenizer: Tokenizes on whitespace
//   - WordPunctTokenizer: Splits into alphabetic and non-alphabetic tokens
//   - SpaceTokenizer: Splits on the space character
//   - TabTokenizer: Splits on the tab character
//   - CharTokenizer: Splits into individual characters
//   - LineTokenizer: Splits into lines
//   - BlanklineTokenizer: Splits on blank lines
//   - MWETokenizer: Merges multi-word expressions
//   - ParallelTokenizer: Concurrently tokenizes large texts using worker pools
//   - SentTokenizer: Simple sentence tokenizer
package tokenizer
