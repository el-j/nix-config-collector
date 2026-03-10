package generator_test

import (
	"strings"
	"testing"
	"time"

	"github.com/el-j/nix-config-collector/internal/core/generator"
	"github.com/el-j/nix-config-collector/internal/core/models"
)

func testScan() models.SystemScan {
	return models.SystemScan{
		ScannedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Hostname:  "test-host",
		Packages: []models.Package{
			{Name: "git", Type: models.PackageTypeBrewFormula, NixName: "git"},
			{Name: "neovim", Type: models.PackageTypeBrewFormula, NixName: "neovim"},
			{Name: "unknown-brew", Type: models.PackageTypeBrewFormula},
			{Name: "SomeApp", Type: models.PackageTypeBrewCask},
		},
		ShellConfig: models.ShellConfig{
			Shell:   "/bin/zsh",
			RCFiles: []string{"~/.zshrc"},
		},
		Dotfiles: []models.Dotfile{
			{Path: "~/.gitconfig", Content: "[user]\n  name = Test User\n  email = test@example.com\n"},
		},
	}
}

func TestGenerateDarwinConfig(t *testing.T) {
	g := generator.New()
	config, err := g.GenerateDarwinConfig(testScan())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(config, "darwin-configuration.nix") {
		t.Error("expected darwin config header")
	}
	if !strings.Contains(config, "homebrew") {
		t.Error("expected homebrew section")
	}
	if !strings.Contains(config, "git") {
		t.Error("expected git in nix packages")
	}
	if !strings.Contains(config, "SomeApp") {
		t.Error("expected SomeApp in casks")
	}
}

func TestGenerateHomeConfig(t *testing.T) {
	g := generator.New()
	config, err := g.GenerateHomeConfig(testScan())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(config, "home.nix") {
		t.Error("expected home.nix header")
	}
	if !strings.Contains(config, "programs.zsh") {
		t.Error("expected zsh program")
	}
	if !strings.Contains(config, "programs.git") {
		t.Error("expected git program")
	}
}

func TestGenerateDarwinConfig_MASAppsWithAppID(t *testing.T) {
	g := generator.New()
	scan := testScan()
	scan.Packages = append(scan.Packages,
		models.Package{Name: "Xcode", Version: "15.2", AppID: "497799835", Type: models.PackageTypeMAS},
		models.Package{Name: "Keynote", Version: "13.2", AppID: "", Type: models.PackageTypeMAS},
	)
	config, err := g.GenerateDarwinConfig(scan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(config, `"Xcode" = 497799835;`) {
		t.Errorf("expected Xcode with numeric App Store ID, got:\n%s", config)
	}
	if !strings.Contains(config, `"Keynote" = 0;`) {
		t.Errorf("expected Keynote fallback to 0, got:\n%s", config)
	}
}

func TestGenerateFlakeConfig(t *testing.T) {
	g := generator.New()
	config, err := g.GenerateFlakeConfig(testScan())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(config, "flake.nix") {
		t.Error("expected flake.nix header")
	}
	if !strings.Contains(config, "nix-darwin") {
		t.Error("expected nix-darwin input")
	}
	if !strings.Contains(config, "home-manager") {
		t.Error("expected home-manager input")
	}
	if !strings.Contains(config, "test-host") {
		t.Error("expected hostname in flake config")
	}
}
