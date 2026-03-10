package models_test

import (
	"testing"
	"time"

	"github.com/el-j/nix-config-collector/internal/core/models"
)

func TestPackageTypes(t *testing.T) {
	tests := []struct {
		name     string
		pkgType  models.PackageType
		expected string
	}{
		{"brew formula", models.PackageTypeBrewFormula, "brew-formula"},
		{"brew cask", models.PackageTypeBrewCask, "brew-cask"},
		{"mas", models.PackageTypeMAS, "mas"},
		{"nix", models.PackageTypeNix, "nix"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.pkgType) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.pkgType)
			}
		})
	}
}

func TestSystemScan(t *testing.T) {
	scan := models.SystemScan{
		ScannedAt: time.Now(),
		Hostname:  "test-host",
		Packages: []models.Package{
			{Name: "git", Type: models.PackageTypeBrewFormula, NixName: "git"},
		},
	}
	if scan.Hostname != "test-host" {
		t.Errorf("expected hostname 'test-host', got %q", scan.Hostname)
	}
	if len(scan.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(scan.Packages))
	}
}

func TestNixConfig(t *testing.T) {
	config := models.NixConfig{
		DarwinConfig: "{ }",
		HomeConfig:   "{ }",
		FlakeConfig:  "{ }",
	}
	if config.DarwinConfig == "" {
		t.Error("expected non-empty DarwinConfig")
	}
}
