package downloader

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestNewDownloader(t *testing.T) {
	d := New("/tmp/nlp-data")
	if d.Dir == "" {
		t.Fatal("dir should not be empty")
	}
	if d.Dir != "/tmp/nlp-data" {
		t.Errorf("expected dir '/tmp/nlp-data', got %q", d.Dir)
	}
}

func TestDownloaderFields(t *testing.T) {
	d := New("/some/dir")
	if d == nil {
		t.Fatal("New returned nil")
	}
}

func TestDownloadPlainFile(t *testing.T) {
	content := []byte("hello world")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))
	defer server.Close()

	dir := t.TempDir()
	d := New(dir)
	dst, err := d.Download(server.URL+"/file.txt", "file.txt", false)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}
	if dst == "" {
		t.Fatal("expected non-empty destination")
	}
	data, err := os.ReadFile(filepath.Join(dir, "file.txt"))
	if err != nil {
		t.Fatalf("could not read downloaded file: %v", err)
	}
	if !bytes.Equal(data, content) {
		t.Errorf("downloaded content mismatch: got %q, want %q", data, content)
	}
}

func TestDownloadUnzip(t *testing.T) {
	// Build a zip file in memory
	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)
	f, err := w.Create("hello.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(f, "hello"); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	zipData := buf.Bytes()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(zipData)
	}))
	defer server.Close()

	dir := t.TempDir()
	d := New(dir)
	_, err = d.Download(server.URL+"/test.zip", "test.zip", true)
	if err != nil {
		t.Fatalf("Download (with unzip) failed: %v", err)
	}
	// Verify the unzipped file exists
	extracted := filepath.Join(dir, "hello.txt")
	if _, err := os.Stat(extracted); err != nil {
		t.Errorf("expected extracted file %q, got error: %v", extracted, err)
	}
}

func TestDownloadHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	dir := t.TempDir()
	d := New(dir)
	_, err := d.Download(server.URL+"/missing.txt", "missing.txt", false)
	if err == nil {
		t.Fatal("expected error for HTTP 404")
	}
}

func TestDownloadUnreachableURL(t *testing.T) {
	dir := t.TempDir()
	d := New(dir)
	_, err := d.Download("http://127.0.0.1:1/never", "x.txt", false)
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}

func TestDownloadNonZipNoUnzip(t *testing.T) {
	content := []byte("plain text")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))
	defer server.Close()

	dir := t.TempDir()
	d := New(dir)
	dst, err := d.Download(server.URL+"/data.txt", "data.txt", true) // unzip=true but not .zip
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}
	if dst == "" {
		t.Fatal("expected non-empty destination")
	}
}

func TestFetchSuccess(t *testing.T) {
	content := []byte("test content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))
	defer server.Close()

	dst := filepath.Join(t.TempDir(), "out.txt")
	if err := fetch(server.URL, dst); err != nil {
		t.Fatalf("fetch failed: %v", err)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, content) {
		t.Errorf("content mismatch: got %q, want %q", data, content)
	}
}

func TestFetchHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	dst := filepath.Join(t.TempDir(), "out.txt")
	err := fetch(server.URL, dst)
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func TestUnzipFile(t *testing.T) {
	// Create a zip file with multiple entries
	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)
	for _, name := range []string{"a.txt", "b.txt", "subdir/c.txt"} {
		f, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(f, "content of "+name); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	zipPath := filepath.Join(dir, "test.zip")
	if err := os.WriteFile(zipPath, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(dir, "extracted")
	if err := unzipFile(zipPath, target); err != nil {
		t.Fatalf("unzipFile failed: %v", err)
	}

	for _, name := range []string{"a.txt", "b.txt", filepath.Join("subdir", "c.txt")} {
		path := filepath.Join(target, name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected %q to exist: %v", path, err)
		}
	}
}

func TestUnzipFileInvalidPath(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "invalid.zip")
	if err := os.WriteFile(zipPath, []byte("not a zip"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "extracted")
	if err := unzipFile(zipPath, target); err == nil {
		t.Fatal("expected error for invalid zip")
	}
}

func TestDownloadCreatesDir(t *testing.T) {
	content := []byte("test")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))
	defer server.Close()

	parent := t.TempDir()
	nested := filepath.Join(parent, "nested", "subdir")
	d := New(nested)
	_, err := d.Download(server.URL+"/f.txt", "f.txt", false)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}
	if _, err := os.Stat(nested); err != nil {
		t.Errorf("nested dir not created: %v", err)
	}
}
