package wsd

import (
	"testing"
)

func TestLesk(t *testing.T) {
	s := Lesk(ContextWords("I deposited money at the bank"), "bank", []Sense{
		{Name: "river", Gloss: []string{"river", "water"}},
		{Name: "finance", Gloss: []string{"money", "deposit"}},
	})
	if s == nil || s.Name != "finance" {
		t.Fatal(s)
	}
}

func TestLeskContextWords(t *testing.T) {
	tokens := ContextWords("Hello World Test")
	if len(tokens) != 3 {
		t.Errorf("expected 3 tokens, got %d", len(tokens))
	}
	if tokens[0] != "hello" || tokens[1] != "world" || tokens[2] != "test" {
		t.Errorf("expected lowercase tokens, got %v", tokens)
	}
}

func TestContextWordsCaseInsensitive(t *testing.T) {
	tokens := ContextWords("HELLO World")
	if tokens[0] != "hello" {
		t.Errorf("expected lowercase, got %s", tokens[0])
	}
}

func TestContextWordsSingle(t *testing.T) {
	tokens := ContextWords("hello")
	if len(tokens) != 1 {
		t.Errorf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0] != "hello" {
		t.Errorf("expected 'hello', got %s", tokens[0])
	}
}

func TestContextWordsEmpty(t *testing.T) {
	tokens := ContextWords("")
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(tokens))
	}
}

func TestContextWordsPunctuation(t *testing.T) {
	tokens := ContextWords("Hello, world! How are you?")
	// Should split on whitespace
	if len(tokens) != 5 {
		t.Errorf("expected 5 tokens, got %d", len(tokens))
	}
}

func TestContextWordsMultipleSpaces(t *testing.T) {
	tokens := ContextWords("hello   world")
	// strings.Fields handles multiple spaces
	if len(tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(tokens))
	}
}

func TestLeskBankFinancial(t *testing.T) {
	senses := []Sense{
		{Name: "river", Gloss: []string{"river", "water", "shore"}},
		{Name: "finance", Gloss: []string{"money", "deposit", "account"}},
	}
	s := Lesk(ContextWords("I deposited money at the bank"), "bank", senses)
	if s == nil {
		t.Fatal("expected non-nil result")
	}
	if s.Name != "finance" {
		t.Errorf("expected 'finance' sense, got %s", s.Name)
	}
}

func TestLeskBankRiver(t *testing.T) {
	senses := []Sense{
		{Name: "river", Gloss: []string{"river", "water", "shore"}},
		{Name: "finance", Gloss: []string{"money", "deposit", "account"}},
	}
	s := Lesk(ContextWords("the boat is by the bank"), "bank", senses)
	if s == nil {
		t.Fatal("expected non-nil result")
	}
	if s.Name != "river" {
		t.Errorf("expected 'river' sense, got %s", s.Name)
	}
}

func TestLeskNoOverlap(t *testing.T) {
	senses := []Sense{
		{Name: "sense1", Gloss: []string{"apple", "fruit"}},
		{Name: "sense2", Gloss: []string{"car", "vehicle"}},
	}
	// Context has no overlap with either sense
	s := Lesk(ContextWords("xyz unknown words"), "word", senses)
	if s == nil {
		t.Fatal("expected non-nil result")
	}
	// Either sense could be returned (no overlap)
	_ = s
}

func TestLeskSingleSense(t *testing.T) {
	senses := []Sense{{Name: "only", Gloss: []string{"unique", "single"}}}
	s := Lesk(ContextWords("unique single item"), "item", senses)
	if s == nil {
		t.Fatal("expected non-nil result")
	}
	if s.Name != "only" {
		t.Errorf("expected 'only', got %s", s.Name)
	}
}

func TestLeskEmptySenses(t *testing.T) {
	s := Lesk(ContextWords("hello world"), "word", []Sense{})
	if s != nil {
		t.Errorf("expected nil for empty senses, got %v", s)
	}
}

func TestLeskEmptyContext(t *testing.T) {
	senses := []Sense{
		{Name: "a", Gloss: []string{"alpha"}},
		{Name: "b", Gloss: []string{"beta"}},
	}
	s := Lesk(ContextWords(""), "word", senses)
	if s == nil {
		t.Fatal("expected non-nil result")
	}
	// With no context overlap, first sense might be returned
}

func TestLeskMultipleWordOverlap(t *testing.T) {
	senses := []Sense{
		{Name: "weak", Gloss: []string{"a", "b"}},
		{Name: "strong", Gloss: []string{"a", "b", "c", "d"}},
	}
	s := Lesk(ContextWords("a b c d"), "word", senses)
	if s == nil {
		t.Fatal("expected non-nil result")
	}
	if s.Name != "strong" {
		t.Errorf("expected 'strong' (more overlap), got %s", s.Name)
	}
}

func TestSenseStruct(t *testing.T) {
	s := Sense{Name: "test", Gloss: []string{"a", "b"}}
	if s.Name != "test" {
		t.Errorf("expected name 'test', got %s", s.Name)
	}
	if len(s.Gloss) != 2 {
		t.Errorf("expected 2 gloss words, got %d", len(s.Gloss))
	}
}

func TestLeskTieBreaking(t *testing.T) {
	// Both senses have exactly one overlapping word
	senses := []Sense{
		{Name: "first", Gloss: []string{"a"}},
		{Name: "second", Gloss: []string{"a"}},
	}
	s := Lesk(ContextWords("a"), "word", senses)
	if s == nil {
		t.Fatal("expected non-nil result")
	}
	// First sense wins in case of tie
	if s.Name != "first" {
		t.Errorf("expected 'first' in case of tie, got %s", s.Name)
	}
}

func TestLeskCaseInsensitive(t *testing.T) {
	senses := []Sense{
		{Name: "a", Gloss: []string{"HELLO"}},
		{Name: "b", Gloss: []string{"goodbye"}},
	}
	// Context uses lowercase, gloss is uppercase
	s := Lesk(ContextWords("hello world"), "word", senses)
	if s == nil {
		t.Fatal("expected non-nil result")
	}
	if s.Name != "a" {
		t.Errorf("expected 'a' (case-insensitive match), got %s", s.Name)
	}
}

func TestLeskEmptyGloss(t *testing.T) {
	senses := []Sense{
		{Name: "empty", Gloss: []string{}},
		{Name: "full", Gloss: []string{"some"}},
	}
	s := Lesk(ContextWords("some text"), "text", senses)
	if s == nil {
		t.Fatal("expected non-nil result")
	}
	// The "full" sense should match "some"
	if s.Name != "full" {
		t.Errorf("expected 'full', got %s", s.Name)
	}
}

func TestLeskMultiWordGloss(t *testing.T) {
	senses := []Sense{
		{Name: "plain", Gloss: []string{"one", "two", "three"}},
		{Name: "rich", Gloss: []string{"a", "b", "c", "d", "e"}},
	}
	// Context overlaps more with "rich"
	s := Lesk(ContextWords("a b c d e f"), "f", senses)
	if s == nil {
		t.Fatal("expected non-nil result")
	}
	if s.Name != "rich" {
		t.Errorf("expected 'rich', got %s", s.Name)
	}
}

func TestContextWordsNewlines(t *testing.T) {
	tokens := ContextWords("hello\nworld")
	// strings.Fields treats \n as whitespace
	if len(tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(tokens))
	}
}

func TestContextWordsTabs(t *testing.T) {
	tokens := ContextWords("hello\tworld")
	if len(tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(tokens))
	}
}

func BenchmarkLesk(b *testing.B) {
	senses := []Sense{
		{Name: "river", Gloss: []string{"river", "water", "shore"}},
		{Name: "finance", Gloss: []string{"money", "deposit", "account"}},
	}
	ctx := ContextWords("I deposited money at the bank")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Lesk(ctx, "bank", senses)
	}
}

func BenchmarkContextWords(b *testing.B) {
	text := "I deposited money at the bank yesterday in the morning"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ContextWords(text)
	}
}
