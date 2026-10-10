package tokenizer

import (
	"regexp"
	"strings"
)

// TweetTokenizer tokenizes Twitter/social media text.
// Ported from nltk.tokenize.casual.TweetTokenizer.
type TweetTokenizer struct {
	preserveCase   bool
	reduceLen      bool
	stripHandles   bool
	matchHandles   bool
	matchHashtags  bool
	matchEmoticons bool
}

// NewTweetTokenizer creates a new TweetTokenizer.
// preserveCase: if false, convert to lowercase (default true)
// reduceLen: if true, reduce repeated characters (default false)
// stripHandles: if true, remove @username handles (default false)
// matchHandles: if true, keep @username as tokens (default true)
// matchHashtags: if true, keep #hashtag as tokens (default true)
// matchEmoticons: if true, keep emoticons as tokens (default true)
func NewTweetTokenizer(preserveCase, reduceLen, stripHandles, matchHandles, matchHashtags, matchEmoticons bool) *TweetTokenizer {
	return &TweetTokenizer{
		preserveCase:   preserveCase,
		reduceLen:      reduceLen,
		stripHandles:   stripHandles,
		matchHandles:   matchHandles,
		matchHashtags:  matchHashtags,
		matchEmoticons: matchEmoticons,
	}
}

// NewDefaultTweetTokenizer creates a TweetTokenizer with sensible defaults.
func NewDefaultTweetTokenizer() *TweetTokenizer {
	return NewTweetTokenizer(true, false, false, true, true, true)
}

// Tokenize tokenizes a tweet text.
func (tt *TweetTokenizer) Tokenize(text string) []string {
	if !tt.preserveCase {
		text = strings.ToLower(text)
	}

	// Build the regex pattern based on options
	pattern := tweetRegexPattern(tt)
	re := regexp.MustCompile(pattern)

	// Find all tokens
	tokens := re.FindAllString(text, -1)

	// Post-processing
	if tt.reduceLen {
		tokens = reduceRepeatedChars(tokens)
	}

	if tt.stripHandles {
		var filtered []string
		handleRe := regexp.MustCompile(`^@[\w_]+$`)
		for _, tok := range tokens {
			if !handleRe.MatchString(tok) {
				filtered = append(filtered, tok)
			}
		}
		tokens = filtered
	}

	return tokens
}

// SpanTokenize returns the spans of the tokens.
func (tt *TweetTokenizer) SpanTokenize(text string) []Span {
	tokens := tt.Tokenize(text)
	return alignTokens(tokens, text)
}

// tweetRegexPattern builds the regex for tweet tokenization.
func tweetRegexPattern(tt *TweetTokenizer) string {
	parts := []string{}

	// Emoticons
	if tt.matchEmoticons {
		parts = append(parts, emoticonPattern)
	}

	// URLs and email
	parts = append(parts, urlPatern)
	parts = append(parts, emailPattern)

	// Phone numbers, IP, dates
	parts = append(parts, phonePattern)
	parts = append(parts, ipPattern)
	parts = append(parts, datePattern)

	// Numbers
	parts = append(parts, numberPattern)

	// Hashtags
	if tt.matchHashtags {
		parts = append(parts, hashtagPattern)
	}

	// Handles
	if tt.matchHandles {
		parts = append(parts, handlePattern)
	}

	// Remaining words and punctuation
	parts = append(parts,
		decorativeDashPattern,
		apostrophePattern,
		remainderPattern,
	)

	return strings.Join(parts, "|")
}

// reduceRepeatedChars reduces repeated characters in tokens.
func reduceRepeatedChars(tokens []string) []string {
	re := regexp.MustCompile(`(.)\1{2,}`)
	result := make([]string, len(tokens))
	for i, tok := range tokens {
		result[i] = re.ReplaceAllString(tok, "${1}${1}")
	}
	return result
}

// Regex patterns for tweet tokenization
var (
	emoticonPattern       = `[<>]?[:;=X8][\-o*]?[\)\]\(\[dDpP/:}{@|,\\]|[:;=X8][\-o*]?[\)\]\(\[dDpP/:}{@|,\\]|[;:][\-o*][\|\\/]|[<]?[;:][\-o*][\(\[\)\]dDpP/:}{@|,\\\]]|\([-\+\|\\/][:;=X8]`
	urlPatern             = `https?:\/\/[^\s<>"']+|www\.[^\s<>"']+`
	emailPattern          = `[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+`
	phonePattern          = `[\+]?[0-9]{1,3}[-\.\s]?[0-9]{1,3}[-\.\s]?[0-9]{4,10}`
	ipPattern             = `\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`
	datePattern           = `\b\d{1,2}[-/.]\d{1,2}[-/.]\d{2,4}\b`
	numberPattern         = `[+\-]?\d+[,./]?\d*[%]?`
	hashtagPattern        = `#[a-zA-Z0-9_]+`
	handlePattern         = `@[\w_]+`
	decorativeDashPattern = `[<>]?[<>\-—–=]+`
	apostrophePattern     = `[a-zA-Z]+'(?:ll|ve|re|s|t|m|d|em|cause|em|n|im)`
	remainderPattern      = `[^\s<>"']+`

	// Standard casual tokenizer function
	casualTokenizeRE = regexp.MustCompile(
		`[<>]?[:;=X8][\-o*]?[\)\]\(\[dDpP/:}{@|,\\]|` +
			`https?:\/\/[^\s<>"']+|www\.[^\s<>"']+|` +
			`[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+|` +
			`[\+]?[0-9]{1,3}[-\.\s]?[0-9]{1,3}[-\.\s]?[0-9]{4,10}|` +
			`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b|` +
			`\b\d{1,2}[-/.]\d{1,2}[-/.]\d{2,4}\b|` +
			`[+\-]?\d+[,./]?\d*[%]?|` +
			`#[\w_]+|@[\w_]+|` +
			`[<>]?[<>\-—–=]+|` +
			`[a-zA-Z]+'(?:ll|ve|re|s|t|m|d|em|cause|n|im)|` +
			`[^\s<>"']+`,
	)
)

// CasualTokenize is the standard casual tokenizer function.
// Equivalent to nltk.tokenize.casual_tokenize().
func CasualTokenize(text string) []string {
	return casualTokenizeRE.FindAllString(text, -1)
}
