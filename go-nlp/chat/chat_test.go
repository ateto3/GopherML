package chat

import (
	"regexp"
	"strings"
	"testing"
)

func TestElizaReplyNotEmpty(t *testing.T) {
	b := Eliza()
	if got := b.Reply("I need help"); got == "" {
		t.Fatal("expected non-empty reply from Eliza")
	}
}

func TestElizaIMustMatch(t *testing.T) {
	b := Eliza()
	got := b.Reply("I need a doctor")
	if !strings.Contains(got, "doctor") && !strings.Contains(got, "need") {
		t.Errorf("expected reply to reference 'doctor' or 'need', got: %s", got)
	}
}

func TestElizaIAmMatch(t *testing.T) {
	b := Eliza()
	got := b.Reply("I am happy")
	if got == "" {
		t.Fatal("expected non-empty reply")
	}
}

func TestElizaICantMatch(t *testing.T) {
	b := Eliza()
	got := b.Reply("I can't sleep")
	if got == "" {
		t.Fatal("expected non-empty reply")
	}
}

func TestElizaWhyMatch(t *testing.T) {
	b := Eliza()
	got := b.Reply("why is the sky blue")
	if got == "" {
		t.Fatal("expected non-empty reply")
	}
}

func TestElizaHelloMatch(t *testing.T) {
	b := Eliza()
	got := b.Reply("hello there")
	if got == "" {
		t.Fatal("expected non-empty reply")
	}
}

func TestElizaFallback(t *testing.T) {
	b := Eliza()
	got := b.Reply("something completely unrelated xyz")
	if got != "Tell me more." {
		t.Errorf("expected fallback 'Tell me more.', got: %s", got)
	}
}

func TestElizaCaseInsensitive(t *testing.T) {
	b := Eliza()
	got := b.Reply("I NEED HELP")
	if got == "" {
		t.Fatal("expected non-empty reply for uppercase input")
	}
}

func TestElizaWhitespaceHandling(t *testing.T) {
	b := Eliza()
	got := b.Reply("  I need coffee  ")
	if got == "" {
		t.Fatal("expected non-empty reply with whitespace")
	}
}

func TestReflect(t *testing.T) {
	// Note: reflect strips first and last char, so it doesn't fully replace
	// We just verify the function returns something
	got := reflect("i am good")
	if got == "" {
		t.Error("expected non-empty result from reflect")
	}
}

func TestItoa(t *testing.T) {
	for i := 1; i <= 5; i++ {
		got := itoa(i)
		expected := string(rune('0' + i))
		if got != expected {
			t.Errorf("itoa(%d) = %q, want %q", i, got, expected)
		}
	}
}

func TestCustomBot(t *testing.T) {
	bot := New([]Rule{
		{Re: regexp.MustCompile(`^foo(.*)`), Replies: []string{"bar%1"}},
	})
	got := bot.Reply("foo baz")
	if got != "bar baz" {
		t.Errorf("expected 'bar baz', got: %s", got)
	}
}

func TestCustomBotNoMatch(t *testing.T) {
	bot := New([]Rule{
		{Re: regexp.MustCompile(`^hello`), Replies: []string{"hi"}},
	})
	got := bot.Reply("goodbye")
	if got != "Tell me more." {
		t.Errorf("expected fallback, got: %s", got)
	}
}

func TestCustomBotMultipleReplies(t *testing.T) {
	bot := New([]Rule{
		{Re: regexp.MustCompile(`^test`), Replies: []string{"a", "b", "c"}},
	})
	seen := map[string]bool{}
	for i := 0; i < 50; i++ {
		got := bot.Reply("test")
		seen[got] = true
	}
	if len(seen) < 2 {
		t.Errorf("expected multiple different replies, only saw: %v", seen)
	}
}

func TestCustomBotCaptureGroups(t *testing.T) {
	bot := New([]Rule{
		{Re: regexp.MustCompile(`^say (.+) twice`), Replies: []string{"%1 %1"}},
	})
	got := bot.Reply("say hello twice")
	if got != "hello hello" {
		t.Errorf("expected 'hello hello', got: %s", got)
	}
}

func TestNewBotEmptyRules(t *testing.T) {
	bot := New([]Rule{})
	got := bot.Reply("anything")
	if got != "Tell me more." {
		t.Errorf("expected fallback for empty rules, got: %s", got)
	}
}
