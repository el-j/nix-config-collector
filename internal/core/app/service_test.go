package app_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/el-j/nix-config-collector/internal/core/app"
	"github.com/el-j/nix-config-collector/internal/core/models"
)

// --- Mock implementations ---

type mockScanner struct {
	packages []models.Package
	services []models.Service
	dotfiles []models.Dotfile
	prefs    []models.SystemPreference
	fonts    []models.Font
	shell    models.ShellConfig
	err      error
}

func (m *mockScanner) ScanPackages(_ context.Context) ([]models.Package, error) {
	return m.packages, m.err
}
func (m *mockScanner) ScanServices(_ context.Context) ([]models.Service, error) {
	return m.services, m.err
}
func (m *mockScanner) ScanDotfiles(_ context.Context, _ []string) ([]models.Dotfile, error) {
	return m.dotfiles, m.err
}
func (m *mockScanner) ScanPreferences(_ context.Context) ([]models.SystemPreference, error) {
	return m.prefs, m.err
}
func (m *mockScanner) ScanFonts(_ context.Context) ([]models.Font, error) {
	return m.fonts, m.err
}
func (m *mockScanner) ScanShellConfig(_ context.Context) (models.ShellConfig, error) {
	return m.shell, m.err
}

type mockWriter struct {
	written      map[string]string
	backups      []string
	writeErr     error
	ensureDirErr error
}

func newMockWriter() *mockWriter {
	return &mockWriter{written: make(map[string]string)}
}

func (m *mockWriter) WriteFile(path, content string) error {
	if m.writeErr != nil {
		return m.writeErr
	}
	m.written[path] = content
	return nil
}
func (m *mockWriter) EnsureDir(path string) error { return m.ensureDirErr }
func (m *mockWriter) BackupFile(path string) error {
	m.backups = append(m.backups, path)
	return nil
}
func (m *mockWriter) FileExists(path string) bool {
	_, ok := m.written[path]
	return ok
}

type mockNotifier struct {
	infos    []string
	warnings []string
	errors   []string
}

func (m *mockNotifier) Info(msg string)               { m.infos = append(m.infos, msg) }
func (m *mockNotifier) Success(msg string)             { m.infos = append(m.infos, msg) }
func (m *mockNotifier) Warning(msg string)             { m.warnings = append(m.warnings, msg) }
func (m *mockNotifier) Error(msg string)               { m.errors = append(m.errors, msg) }
func (m *mockNotifier) Progress(_, _ int, msg string)  { m.infos = append(m.infos, msg) }
func (m *mockNotifier) Confirm(_ string) bool          { return true }

type mockGenerator struct {
	darwinOut string
	homeOut   string
	flakeOut  string
	err       error
}

func (m *mockGenerator) GenerateDarwinConfig(_ models.SystemScan) (string, error) {
	return m.darwinOut, m.err
}
func (m *mockGenerator) GenerateHomeConfig(_ models.SystemScan) (string, error) {
	return m.homeOut, m.err
}
func (m *mockGenerator) GenerateFlakeConfig(_ models.SystemScan) (string, error) {
	return m.flakeOut, m.err
}

// --- Tests ---

func newTestService(scanner *mockScanner, writer *mockWriter, notifier *mockNotifier, gen *mockGenerator) *app.Service {
	return app.NewService(scanner, writer, notifier, gen)
}

func TestScan_Success(t *testing.T) {
	scanner := &mockScanner{
		packages: []models.Package{
			{Name: "git", Type: models.PackageTypeBrewFormula},
			{Name: "ripgrep", Type: models.PackageTypeBrewFormula},
		},
		services: []models.Service{{Name: "com.example.svc", Type: models.ServiceTypeLaunchAgent, Enabled: true}},
		dotfiles: []models.Dotfile{{Path: "~/.zshrc", Content: "export PATH=$PATH:/usr/local/bin"}},
		shell:    models.ShellConfig{Shell: "/bin/zsh"},
	}
	writer := newMockWriter()
	notifier := &mockNotifier{}
	gen := &mockGenerator{}

	svc := newTestService(scanner, writer, notifier, gen)
	result, err := svc.Scan(context.Background(), app.ScanOptions{})
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}
	if result == nil {
		t.Fatal("Scan() returned nil")
	}
	// git and ripgrep should be mapped to nix names
	var gitPkg *models.Package
	for i := range result.Packages {
		if result.Packages[i].Name == "git" {
			gitPkg = &result.Packages[i]
		}
	}
	if gitPkg == nil {
		t.Fatal("git package not found in scan results")
	}
	if gitPkg.NixName != "git" {
		t.Errorf("git NixName = %q, want %q", gitPkg.NixName, "git")
	}
	if len(result.Services) != 1 {
		t.Errorf("services count = %d, want 1", len(result.Services))
	}
	if result.ShellConfig.Shell != "/bin/zsh" {
		t.Errorf("shell = %q, want /bin/zsh", result.ShellConfig.Shell)
	}
}

func TestScan_ScannerError(t *testing.T) {
	scanner := &mockScanner{err: errors.New("brew not found")}
	svc := newTestService(scanner, newMockWriter(), &mockNotifier{}, &mockGenerator{})
	_, err := svc.Scan(context.Background(), app.ScanOptions{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "scanning packages") {
		t.Errorf("error %q does not mention 'scanning packages'", err.Error())
	}
}

func TestGenerate_NoOutputDir(t *testing.T) {
	scanner := &mockScanner{}
	writer := newMockWriter()
	gen := &mockGenerator{
		darwinOut: "# darwin config",
		homeOut:   "# home config",
		flakeOut:  "# flake config",
	}
	svc := newTestService(scanner, writer, &mockNotifier{}, gen)
	cfg, err := svc.Generate(context.Background(), app.GenerateOptions{})
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	if cfg.DarwinConfig != "# darwin config" {
		t.Errorf("DarwinConfig = %q", cfg.DarwinConfig)
	}
	if len(writer.written) != 0 {
		t.Errorf("expected no files written, got %d", len(writer.written))
	}
}

func TestGenerate_WithOutputDir(t *testing.T) {
	scanner := &mockScanner{}
	writer := newMockWriter()
	gen := &mockGenerator{
		darwinOut: "# darwin",
		homeOut:   "# home",
		flakeOut:  "# flake",
	}
	svc := newTestService(scanner, writer, &mockNotifier{}, gen)
	_, err := svc.Generate(context.Background(), app.GenerateOptions{OutputDir: "/tmp/nix-out"})
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}
	expectedFiles := []string{
		"/tmp/nix-out/darwin-configuration.nix",
		"/tmp/nix-out/home.nix",
		"/tmp/nix-out/flake.nix",
	}
	for _, f := range expectedFiles {
		if _, ok := writer.written[f]; !ok {
			t.Errorf("expected file %q to be written", f)
		}
	}
}

func TestGenerate_WriterError(t *testing.T) {
	scanner := &mockScanner{}
	writer := &mockWriter{
		written:  make(map[string]string),
		writeErr: errors.New("disk full"),
	}
	gen := &mockGenerator{darwinOut: "# d", homeOut: "# h", flakeOut: "# f"}
	svc := newTestService(scanner, writer, &mockNotifier{}, gen)
	_, err := svc.Generate(context.Background(), app.GenerateOptions{OutputDir: "/tmp/out"})
	if err == nil {
		t.Fatal("expected error from writer, got nil")
	}
}

func TestGenerate_GeneratorError(t *testing.T) {
	scanner := &mockScanner{}
	gen := &mockGenerator{err: errors.New("template error")}
	svc := newTestService(scanner, newMockWriter(), &mockNotifier{}, gen)
	_, err := svc.Generate(context.Background(), app.GenerateOptions{})
	if err == nil {
		t.Fatal("expected error from generator, got nil")
	}
	if !strings.Contains(err.Error(), "generating darwin config") {
		t.Errorf("error %q should mention 'generating darwin config'", err.Error())
	}
}
