package macos_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/el-j/nix-config-collector/internal/adapters/macos"
)

func TestNew(t *testing.T) {
	s := macos.New()
	if s == nil {
		t.Fatal("New() returned nil")
	}
}

func TestWithTimeout(t *testing.T) {
	s := macos.New().WithTimeout(5 * time.Second)
	if s == nil {
		t.Fatal("WithTimeout() returned nil")
	}
}

func TestScanPackages_GracefulOnMissingCommands(t *testing.T) {
	// On Linux/CI, brew and mas are not installed.
	// Scanner must return empty slice, not error.
	s := macos.New().WithTimeout(5 * time.Second)
	pkgs, err := s.ScanPackages(context.Background())
	if err != nil {
		t.Fatalf("ScanPackages() error = %v", err)
	}
	// On macOS with brew: non-nil slice; on Linux: nil/empty is fine
	_ = pkgs
}

func TestScanServices_GracefulOnMissingCommands(t *testing.T) {
	s := macos.New().WithTimeout(5 * time.Second)
	svcs, err := s.ScanServices(context.Background())
	if err != nil {
		t.Fatalf("ScanServices() error = %v", err)
	}
	_ = svcs
}

func TestScanDotfiles_DefaultPaths(t *testing.T) {
	// With nil paths, uses defaults. Most won't exist in CI, so result may be empty.
	s := macos.New()
	dotfiles, err := s.ScanDotfiles(context.Background(), nil)
	if err != nil {
		t.Fatalf("ScanDotfiles() error = %v", err)
	}
	_ = dotfiles
}

func TestScanDotfiles_ExistingFile(t *testing.T) {
	// Create a temp file and scan it
	f, err := os.CreateTemp("", "test-dotfile-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	content := "export PATH=$PATH:/usr/local/bin\nalias ll='ls -la'\n"
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()

	s := macos.New()
	dotfiles, err := s.ScanDotfiles(context.Background(), []string{f.Name()})
	if err != nil {
		t.Fatalf("ScanDotfiles() error = %v", err)
	}
	if len(dotfiles) != 1 {
		t.Fatalf("expected 1 dotfile, got %d", len(dotfiles))
	}
	if dotfiles[0].Content != content {
		t.Errorf("content mismatch: got %q, want %q", dotfiles[0].Content, content)
	}
}

func TestScanDotfiles_NonExistentFile(t *testing.T) {
	s := macos.New()
	dotfiles, err := s.ScanDotfiles(context.Background(), []string{"/nonexistent/path/.zshrc"})
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(dotfiles) != 0 {
		t.Errorf("expected 0 dotfiles for missing path, got %d", len(dotfiles))
	}
}

func TestScanDotfiles_Directory(t *testing.T) {
	// Create a temp dir with some text files
	dir, err := os.MkdirTemp("", "test-dotdir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Write a text file
	if err := os.WriteFile(dir+"/config.toml", []byte("[settings]\nkey = \"value\"\n"), 0600); err != nil {
		t.Fatal(err)
	}

	s := macos.New()
	dotfiles, err := s.ScanDotfiles(context.Background(), []string{dir})
	if err != nil {
		t.Fatalf("ScanDotfiles() error = %v", err)
	}
	if len(dotfiles) != 1 {
		t.Fatalf("expected 1 dotfile from dir, got %d", len(dotfiles))
	}
}

func TestScanPreferences_GracefulOnMissingDefaults(t *testing.T) {
	// 'defaults' command not available on Linux
	s := macos.New().WithTimeout(3 * time.Second)
	prefs, err := s.ScanPreferences(context.Background())
	if err != nil {
		t.Fatalf("ScanPreferences() error = %v", err)
	}
	_ = prefs
}

func TestScanFonts_GracefulOnMissingDirs(t *testing.T) {
	// Font directories don't exist on Linux
	s := macos.New()
	fonts, err := s.ScanFonts(context.Background())
	if err != nil {
		t.Fatalf("ScanFonts() error = %v", err)
	}
	_ = fonts
}

func TestScanShellConfig(t *testing.T) {
	s := macos.New()
	cfg, err := s.ScanShellConfig(context.Background())
	if err != nil {
		t.Fatalf("ScanShellConfig() error = %v", err)
	}
	if cfg.Shell == "" {
		t.Error("Shell should not be empty (falls back to /bin/zsh)")
	}
}

func TestScanPackages_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	s := macos.New()
	// Should not panic; may return error or empty results
	_, _ = s.ScanPackages(ctx)
}
