package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Writer implements ports.FileWriter for local filesystem
type Writer struct{}

// New creates a new filesystem Writer
func New() *Writer {
	return &Writer{}
}

func (w *Writer) WriteFile(path string, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	return os.WriteFile(path, []byte(content), 0640)
}

func (w *Writer) EnsureDir(path string) error {
	return os.MkdirAll(path, 0750)
}

func (w *Writer) BackupFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	backup := path + ".bak." + time.Now().Format("20060102150405")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading file for backup: %w", err)
	}
	return os.WriteFile(backup, data, 0640)
}

func (w *Writer) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
