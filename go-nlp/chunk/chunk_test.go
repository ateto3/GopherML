package chunk

import "testing"

func TestRegexpChunker(t *testing.T) {
	chunker := NewRegexpChunker("NP", `^(DT|JJ|NN.*)$`)
	tokens := []TaggedToken{
		{Word: "the", Tag: "DT"},
		{Word: "big", Tag: "JJ"},
		{Word: "dog", Tag: "NN"},
		{Word: "barks", Tag: "VBZ"},
	}
	tree := chunker.Parse(tokens)
	if len(tree.Children) != 4 {
		t.Fatalf("unexpected tree size: %d", len(tree.Children))
	}
}

func TestParseSents(t *testing.T) {
	chunker := NewRegexpChunker("NP", `^(DT|JJ|NN.*)$`)
	sents := [][]TaggedToken{{{Word: "the", Tag: "DT"}}, {{Word: "dog", Tag: "NN"}}}
	out := ParseSents(chunker, sents)
	if len(out) != 2 {
		t.Fatalf("expected 2 trees, got %d", len(out))
	}
}

func TestChunkStringAndRules(t *testing.T) {
	cs := NewChunkString([]TaggedToken{{Tag: "DT"}, {Tag: "JJ"}, {Tag: "NN"}, {Tag: "VB"}}, "S")
	NewChunkRule(`^<DT>$|^<JJ>$|^<NN.*>$`).Apply(cs)
	if cs.String() == "" {
		t.Fatal("expected chunk string")
	}
	NewStripRule(`^<JJ>$`).Apply(cs)
}

func TestNamedEntityChunker(t *testing.T) {
	ne := NewNamedEntityChunker()
	tree := ne.Parse([]TaggedToken{{Word: "John", Tag: "NNP"}, {Word: "Smith", Tag: "NNP"}, {Word: "runs", Tag: "VBZ"}})
	if len(tree.Children) == 0 {
		t.Fatal("expected chunks")
	}
}
