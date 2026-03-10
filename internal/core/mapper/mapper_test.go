package mapper_test

import (
	"testing"

	"github.com/el-j/nix-config-collector/internal/core/mapper"
	"github.com/el-j/nix-config-collector/internal/core/models"
)

func TestMapPackages(t *testing.T) {
	m := mapper.New()
	packages := []models.Package{
		{Name: "git", Type: models.PackageTypeBrewFormula},
		{Name: "unknown-pkg", Type: models.PackageTypeBrewFormula},
		{Name: "SomeApp", Type: models.PackageTypeBrewCask},
	}

	mapped, err := m.MapPackages(packages)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mapped) != 3 {
		t.Fatalf("expected 3 packages, got %d", len(mapped))
	}
	if mapped[0].NixName != "git" {
		t.Errorf("expected NixName 'git', got %q", mapped[0].NixName)
	}
	if mapped[1].NixName != "" {
		t.Errorf("expected empty NixName for unknown package, got %q", mapped[1].NixName)
	}
	if mapped[2].NixName != "" {
		t.Errorf("expected empty NixName for cask, got %q", mapped[2].NixName)
	}
}

func TestLookupNixName(t *testing.T) {
	tests := []struct {
		brew     string
		expected string
		found    bool
	}{
		{"git", "git", true},
		{"neovim", "neovim", true},
		{"ripgrep", "ripgrep", true},
		{"nonexistent", "", false},
		{"GIT", "git", true}, // case-insensitive
	}
	for _, tt := range tests {
		t.Run(tt.brew, func(t *testing.T) {
			name, ok := mapper.LookupNixName(tt.brew)
			if ok != tt.found {
				t.Errorf("LookupNixName(%q): found=%v, expected found=%v", tt.brew, ok, tt.found)
			}
			if name != tt.expected {
				t.Errorf("LookupNixName(%q): name=%q, expected %q", tt.brew, name, tt.expected)
			}
		})
	}
}

func TestGetUnmappedPackages(t *testing.T) {
	packages := []models.Package{
		{Name: "git", Type: models.PackageTypeBrewFormula, NixName: "git"},
		{Name: "unknown", Type: models.PackageTypeBrewFormula},
		{Name: "SomeApp", Type: models.PackageTypeBrewCask},
	}
	unmapped := mapper.GetUnmappedPackages(packages)
	if len(unmapped) != 1 {
		t.Fatalf("expected 1 unmapped, got %d", len(unmapped))
	}
	if unmapped[0].Name != "unknown" {
		t.Errorf("expected 'unknown', got %q", unmapped[0].Name)
	}
}

func TestSplitByDestination(t *testing.T) {
	packages := []models.Package{
		{Name: "git", Type: models.PackageTypeBrewFormula, NixName: "git"},
		{Name: "unknown", Type: models.PackageTypeBrewFormula},
		{Name: "app", Type: models.PackageTypeBrewCask},
		{Name: "MAS App", Type: models.PackageTypeMAS},
	}
	nixPkgs, brewFormulae, brewCasks, masApps := mapper.SplitByDestination(packages)
	if len(nixPkgs) != 1 {
		t.Errorf("expected 1 nix package, got %d", len(nixPkgs))
	}
	if len(brewFormulae) != 1 {
		t.Errorf("expected 1 brew formula, got %d", len(brewFormulae))
	}
	if len(brewCasks) != 1 {
		t.Errorf("expected 1 brew cask, got %d", len(brewCasks))
	}
	if len(masApps) != 1 {
		t.Errorf("expected 1 MAS app, got %d", len(masApps))
	}
}
