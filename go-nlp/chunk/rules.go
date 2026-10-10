package chunk

import "regexp"

// ChunkRule chunks tokens whose tag sequence matches the regex.
type ChunkRule struct{ Re *regexp.Regexp }

// NewChunkRule creates a ChunkRule.
func NewChunkRule(pattern string) ChunkRule { return ChunkRule{Re: regexp.MustCompile(pattern)} }

// Apply marks matching tokens as chunked.
func (r ChunkRule) Apply(cs *ChunkString) {
	for i, tag := range cs.Tags {
		if r.Re.MatchString("<" + tag + ">") {
			cs.mark(i, i+1, true)
		}
	}
}

// StripRule removes chunk marks for matching tags.
type StripRule struct{ Re *regexp.Regexp }

// NewStripRule creates a StripRule.
func NewStripRule(pattern string) StripRule { return StripRule{Re: regexp.MustCompile(pattern)} }

// Apply removes chunk marks.
func (r StripRule) Apply(cs *ChunkString) {
	for i, tag := range cs.Tags {
		if r.Re.MatchString("<" + tag + ">") {
			cs.mark(i, i+1, false)
		}
	}
}

// UnChunkRule aliases StripRule for API parity.
type UnChunkRule = StripRule
