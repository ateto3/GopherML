package downloader

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Downloader struct{ Dir string }

func New(dir string) *Downloader { return &Downloader{Dir: dir} }

func (d *Downloader) Download(url, name string, unzip bool) (string, error) {
	if err := os.MkdirAll(d.Dir, 0o755); err != nil { return "", err }
	dst := filepath.Join(d.Dir, name)
	if err := fetch(url, dst); err != nil { return "", err }
	if unzip && strings.HasSuffix(strings.ToLower(name), ".zip") { return dst, unzipFile(dst, d.Dir) }
	return dst, nil
}

func fetch(url, dst string) error {
	resp, err := http.Get(url)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { return fmt.Errorf("download failed: %s", resp.Status) }
	f, err := os.Create(dst)
	if err != nil { return err }
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func unzipFile(zipPath, dir string) error {
	zr, err := zip.OpenReader(zipPath)
	if err != nil { return err }
	defer zr.Close()
	for _, f := range zr.File {
		p := filepath.Join(dir, f.Name)
		if !strings.HasPrefix(filepath.Clean(p), filepath.Clean(dir)+string(os.PathSeparator)) { continue }
		if f.FileInfo().IsDir() { os.MkdirAll(p, 0o755); continue }
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil { return err }
		rc, err := f.Open()
		if err != nil { return err }
		w, err := os.Create(p)
		if err != nil { rc.Close(); return err }
		_, err = io.Copy(w, rc)
		w.Close(); rc.Close()
		if err != nil { return err }
	}
	return nil
}
