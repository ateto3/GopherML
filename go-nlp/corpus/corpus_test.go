package corpus

import (
	"testing"
	"testing/fstest"
)

func TestFSReader(t *testing.T) {
	fsys := fstest.MapFS{
		"a.txt": {Data: []byte("hello world")},
		"b.txt": {Data: []byte("go nlp")},
	}
	r := NewFSReader(fsys, []string{"a.txt", "b.txt"})
	words, err := r.Words()
	if err != nil || len(words) != 4 {
		t.Fatalf("unexpected words: %v %v", words, err)
	}
}

func TestFSReaderWithConfig(t *testing.T) {
	fsys := fstest.MapFS{
		"a.txt": {Data: []byte("Hello World! Go NLP.")},
	}
	r := NewFSReader(fsys, []string{"a.txt"}, NewLowercase(), NewRemovePunctuation())
	words, err := r.Words()
	if err != nil || len(words) != 4 {
		t.Fatalf("unexpected words: %v %v", words, err)
	}

	expected := []string{"hello", "world", "go", "nlp"}
	for i, w := range words {
		if w != expected[i] {
			t.Errorf("expected '%s', got '%s'", expected[i], w)
		}
	}
}

func TestLazy(t *testing.T) {
	l := NewLazy(func() (Reader, error) {
		return NewFSReader(fstest.MapFS{"a.txt": {Data: []byte("hello")}}, []string{"a.txt"}), nil
	})
	words, err := l.Words()
	if err != nil || len(words) != 1 {
		t.Fatalf("unexpected lazy result: %v %v", words, err)
	}
}

func TestFSReaderRaw(t *testing.T) {
	fsys := fstest.MapFS{
		"a.txt": {Data: []byte("hello world")},
		"b.txt": {Data: []byte("go nlp")},
	}
	r := NewFSReader(fsys, []string{"a.txt", "b.txt"})
	raw, err := r.Raw()
	if err != nil || len(raw) == 0 {
		t.Fatalf("unexpected raw: %v %v", raw, err)
	}
}

func TestFSReaderVocab(t *testing.T) {
	fsys := fstest.MapFS{
		"a.txt": {Data: []byte("hello world hello")},
	}
	r := NewFSReader(fsys, []string{"a.txt"})
	vocab, err := r.Vocab()
	if err != nil || len(vocab) != 2 {
		t.Fatalf("unexpected vocab: %v %v", vocab, err)
	}
}

func TestFSReaderSize(t *testing.T) {
	fsys := fstest.MapFS{
		"a.txt": {Data: []byte("hello world go nlp")},
	}
	r := NewFSReader(fsys, []string{"a.txt"})
	size, err := r.Size()
	if err != nil || size != 4 {
		t.Fatalf("unexpected size: %v %v", size, err)
	}
}

func TestFSReaderRawCache(t *testing.T) {
	fsys := fstest.MapFS{
		"a.txt": {Data: []byte("hello world")},
	}
	r := NewFSReader(fsys, []string{"a.txt"})

	// First call
	raw1, err := r.Raw("a.txt")
	if err != nil {
		t.Fatalf("first raw failed: %v", err)
	}

	// Second call should use cache
	raw2, err := r.Raw("a.txt")
	if err != nil {
		t.Fatalf("second raw failed: %v", err)
	}

	if raw1 != raw2 {
		t.Error("raw cache mismatch")
	}
}
