package tagger

import "testing"

func TestParallelPosTag(t *testing.T) {
	tokens := []string{"Natural", "Language", "Processing", "with", "Go", "is", "blazing", "fast", "."}
	tags := ParallelPosTag(tokens)

	if len(tags) != len(tokens) {
		t.Fatalf("expected %d tags, got %d", len(tokens), len(tags))
	}
	if tags[3].Tag != "IN" {
		t.Fatalf("expected 'with' to be tagged as IN, got %q", tags[3].Tag)
	}
	if tags[5].Tag != "VB" {
		t.Fatalf("expected 'is' to be tagged as VB, got %q", tags[5].Tag)
	}
	if tags[6].Tag != "VBG" {
		t.Fatalf("expected 'blazing' to be tagged as VBG, got %q", tags[6].Tag)
	}
	if tags[8].Tag != "." {
		t.Fatalf("expected '.' to be tagged as ., got %q", tags[8].Tag)
	}
}

func TestDefaultTagger(t *testing.T) {
	tagger := NewDefaultTagger("NN")
	tokens := []string{"hello", "world"}
	tags := tagger.Tag(tokens)

	if tags[0].Tag != "NN" || tags[1].Tag != "NN" {
		t.Fatalf("unexpected default tags: %#v", tags)
	}
}
