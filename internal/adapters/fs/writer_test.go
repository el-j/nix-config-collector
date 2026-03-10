package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	fsadapter "github.com/el-j/nix-config-collector/internal/adapters/fs"
)

func TestWriteFile(t *testing.T) {
	w := fsadapter.New()
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "test.txt")

	if err := w.WriteFile(path, "hello world"); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("expected 'hello world', got %q", string(data))
	}
}

func TestEnsureDir(t *testing.T) {
	w := fsadapter.New()
	dir := t.TempDir()
	newDir := filepath.Join(dir, "a", "b", "c")

	if err := w.EnsureDir(newDir); err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}

	info, err := os.Stat(newDir)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory")
	}
}

func TestFileExists(t *testing.T) {
	w := fsadapter.New()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	if w.FileExists(path) {
		t.Error("expected file to not exist")
	}

	if err := os.WriteFile(path, []byte("data"), 0640); err != nil {
		t.Fatalf("creating test file: %v", err)
	}

	if !w.FileExists(path) {
		t.Error("expected file to exist")
	}
}

func TestBackupFile(t *testing.T) {
	w := fsadapter.New()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	// Backup of non-existent file should not error
	if err := w.BackupFile(path); err != nil {
		t.Fatalf("BackupFile on non-existent file: %v", err)
	}

	// Create file and backup
	if err := os.WriteFile(path, []byte("original"), 0640); err != nil {
		t.Fatalf("creating test file: %v", err)
	}

	if err := w.BackupFile(path); err != nil {
		t.Fatalf("BackupFile failed: %v", err)
	}

	// Backup file should exist
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading dir: %v", err)
	}
	if len(entries) < 2 {
		t.Error("expected backup file to be created")
	}
}
