package featstruct

import (
	"strings"
	"testing"
)

func TestUnifyBasic(t *testing.T) {
	a := New(map[string]any{"cat": "np", "agr": New(map[string]any{"num": "sg"})})
	b := New(map[string]any{"agr": New(map[string]any{"pers": "3"})})
	u, err := Unify(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Get("agr", "num") != "sg" {
		t.Fatalf("expected num=sg, got %v", u)
	}
}

func TestUnifyEmptyStructures(t *testing.T) {
	a := New(map[string]any{})
	b := New(map[string]any{})
	u, err := Unify(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestUnifyDisjoint(t *testing.T) {
	a := New(map[string]any{"x": 1})
	b := New(map[string]any{"y": 2})
	u, err := Unify(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Get("x") != 1 || u.Get("y") != 2 {
		t.Errorf("expected both x and y, got %v", u)
	}
}

func TestUnifyOverlapping(t *testing.T) {
	a := New(map[string]any{"x": 1, "y": 2})
	b := New(map[string]any{"y": 2, "z": 3})
	u, err := Unify(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Get("x") != 1 || u.Get("y") != 2 || u.Get("z") != 3 {
		t.Errorf("expected all keys, got %v", u)
	}
}

func TestUnifyClash(t *testing.T) {
	a := New(map[string]any{"x": 1})
	b := New(map[string]any{"x": 2})
	_, err := Unify(a, b)
	if err == nil {
		t.Fatal("expected error for feature clash")
	}
	if !strings.Contains(err.Error(), "clash") {
		t.Errorf("expected 'clash' in error, got: %v", err)
	}
}

func TestUnifyDeepClash(t *testing.T) {
	a := New(map[string]any{"f": New(map[string]any{"x": 1})})
	b := New(map[string]any{"f": New(map[string]any{"x": 2})})
	_, err := Unify(a, b)
	if err == nil {
		t.Fatal("expected error for deep feature clash")
	}
}

func TestUnifyDeepMerge(t *testing.T) {
	a := New(map[string]any{"f": New(map[string]any{"x": 1})})
	b := New(map[string]any{"f": New(map[string]any{"y": 2})})
	u, err := Unify(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Get("f", "x") != 1 || u.Get("f", "y") != 2 {
		t.Errorf("expected nested merge, got %v", u)
	}
}

func TestNew(t *testing.T) {
	f := New(map[string]any{"a": 1, "b": 2})
	if f.Get("a") != 1 || f.Get("b") != 2 {
		t.Errorf("expected New to copy map, got %v", f)
	}
}

func TestNewEmpty(t *testing.T) {
	f := New(map[string]any{})
	if len(f) != 0 {
		t.Errorf("expected empty struct, got %v", f)
	}
}

func TestGetNested(t *testing.T) {
	f := New(map[string]any{
		"a": New(map[string]any{
			"b": New(map[string]any{
				"c": "deep",
			}),
		}),
	})
	if f.Get("a", "b", "c") != "deep" {
		t.Errorf("expected 'deep', got %v", f.Get("a", "b", "c"))
	}
}

func TestGetMissing(t *testing.T) {
	f := New(map[string]any{"a": 1})
	if f.Get("nonexistent") != nil {
		t.Error("expected nil for missing key")
	}
}

func TestGetThroughNonMapping(t *testing.T) {
	f := New(map[string]any{"a": "string-value"})
	if f.Get("a", "b") != nil {
		t.Error("expected nil when traversing through non-mapping value")
	}
}

func TestSet(t *testing.T) {
	f := New(map[string]any{})
	f.Set("value", "key")
	if f.Get("key") != "value" {
		t.Errorf("expected 'value', got %v", f.Get("key"))
	}
}

func TestSetNested(t *testing.T) {
	f := New(map[string]any{})
	f.Set("deep", "a", "b", "c")
	if f.Get("a", "b", "c") != "deep" {
		t.Errorf("expected 'deep', got %v", f.Get("a", "b", "c"))
	}
}

func TestSetEmpty(t *testing.T) {
	f := New(map[string]any{})
	f.Set("v")
	if len(f) != 0 {
		t.Error("Set with no path should not modify")
	}
}

func TestSetCreatesIntermediate(t *testing.T) {
	f := New(map[string]any{})
	f.Set("v", "a", "b")
	a := f.Get("a")
	nested, ok := a.(FeatStruct)
	if !ok {
		t.Fatalf("expected FeatStruct, got %T", a)
	}
	if nested.Get("b") != "v" {
		t.Errorf("expected 'v', got %v", nested.Get("b"))
	}
}

func TestSetOverwritesPrimitive(t *testing.T) {
	f := New(map[string]any{"a": "string"})
	f.Set(42, "a", "b")
	if f.Get("a", "b") != 42 {
		t.Errorf("expected to overwrite, got %v", f.Get("a", "b"))
	}
}

func TestString(t *testing.T) {
	f := New(map[string]any{"b": 2, "a": 1})
	s := f.String()
	if !strings.HasPrefix(s, "[") || !strings.HasSuffix(s, "]") {
		t.Errorf("expected bracketed output, got: %s", s)
	}
	if !strings.Contains(s, "a=1") || !strings.Contains(s, "b=2") {
		t.Errorf("expected both entries, got: %s", s)
	}
	aIdx := strings.Index(s, "a=")
	bIdx := strings.Index(s, "b=")
	if aIdx > bIdx {
		t.Errorf("expected keys sorted alphabetically, got: %s", s)
	}
}

func TestStringEmpty(t *testing.T) {
	f := New(map[string]any{})
	if f.String() != "[]" {
		t.Errorf("expected '[]', got: %s", f.String())
	}
}

func TestParseEmpty(t *testing.T) {
	f := Parse("")
	if len(f) != 0 {
		t.Errorf("expected empty struct, got %v", f)
	}
}

func TestParseBrackets(t *testing.T) {
	f := Parse("[]")
	if len(f) != 0 {
		t.Errorf("expected empty struct, got %v", f)
	}
}

func TestParseSimple(t *testing.T) {
	f := Parse("a=1 b=hello")
	if f.Get("a") != "1" {
		t.Errorf("expected a=1, got %v", f.Get("a"))
	}
	if f.Get("b") != "hello" {
		t.Errorf("expected b=hello, got %v", f.Get("b"))
	}
}

func TestParseNestedAsString(t *testing.T) {
	f := Parse("f=[x=1]")
	if f.Get("f") != "x=1" {
		t.Errorf("expected f='x=1', got %v", f.Get("f"))
	}
}

func TestParseTrimsWhitespace(t *testing.T) {
	f := Parse("  a=1   b=2  ")
	if f.Get("a") != "1" || f.Get("b") != "2" {
		t.Errorf("expected trimmed parse, got %v", f)
	}
}

func TestUnifyImmutableOriginal(t *testing.T) {
	a := New(map[string]any{"x": 1})
	b := New(map[string]any{"y": 2})
	Unify(a, b)
	if a.Get("y") != nil {
		t.Error("Unify should not modify original 'a'")
	}
	if b.Get("x") != nil {
		t.Error("Unify should not modify original 'b'")
	}
}

func TestUnifyNestedClone(t *testing.T) {
	a := New(map[string]any{"f": New(map[string]any{"x": 1})})
	b := New(map[string]any{"f": New(map[string]any{"y": 2})})
	u, err := Unify(a, b)
	if err != nil {
		t.Fatal(err)
	}
	u.Set(99, "f", "x")
	if a.Get("f", "x") != 1 {
		t.Error("mutating result should not affect input")
	}
}

func BenchmarkUnify(b *testing.B) {
	a := New(map[string]any{"f": New(map[string]any{"x": 1, "y": 2})})
	c := New(map[string]any{"f": New(map[string]any{"z": 3})})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Unify(a, c)
	}
}
