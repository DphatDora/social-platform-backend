package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// RotateWriter appends to path and renames the file when size exceeds maxBytes.
// All methods are safe for concurrent use from multiple goroutines.
type RotateWriter struct {
	path     string
	maxBytes int64
	mu       sync.Mutex
	file     *os.File
	size     int64
}

func NewRotateWriter(path string, maxBytes int64) (*RotateWriter, error) {
	if maxBytes <= 0 {
		return nil, fmt.Errorf("maxBytes must be positive")
	}
	w := &RotateWriter{path: path, maxBytes: maxBytes}
	if err := w.openExisting(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *RotateWriter) Path() string {
	return w.path
}

func (w *RotateWriter) openExisting() error {
	dir := filepath.Dir(w.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(w.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	st, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}
	if w.file != nil {
		_ = w.file.Close()
	}
	w.file = f
	w.size = st.Size()
	return nil
}

func (w *RotateWriter) rotate() error {
	if w.file != nil {
		_ = w.file.Close()
		w.file = nil
	}
	st, err := os.Stat(w.path)
	if err == nil && st.Size() > 0 {
		dir := filepath.Dir(w.path)
		base := strings.TrimSuffix(filepath.Base(w.path), filepath.Ext(w.path))
		t := time.Now().UTC()
		ts := strings.ReplaceAll(t.Format("2006-01-02T15-04-05"), ":", "-")
		rotated := filepath.Join(dir, fmt.Sprintf("%s-%s-%d.log", base, ts, t.UnixNano()))
		if err := os.Rename(w.path, rotated); err != nil {
			return err
		}
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	return w.openExisting()
}

func (w *RotateWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		if err := w.openExisting(); err != nil {
			return 0, err
		}
	}

	if w.size+int64(len(p)) > w.maxBytes {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = w.file.Write(p)
	w.size += int64(n)
	return n, err
}

func (w *RotateWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

func (w *RotateWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		return err
	}
	return nil
}
