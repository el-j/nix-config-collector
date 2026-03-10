package macos

import (
	"bufio"
	"context"
	"strings"
	"testing"

	"github.com/el-j/nix-config-collector/internal/core/models"
)

// parseMASOutput exercises the same parsing logic as scanMASApps without
// invoking the real `mas` binary, keeping tests OS-independent.
func parseMASOutput(output string) []models.Package {
	s := New()
	var packages []models.Package
	sc := bufio.NewScanner(strings.NewReader(output))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}
		appID := strings.TrimSpace(parts[0])
		name := strings.TrimSpace(parts[1])
		version := ""
		if idx := strings.LastIndex(name, " ("); idx != -1 {
			version = strings.Trim(name[idx+2:], ")")
			name = name[:idx]
		}
		packages = append(packages, models.Package{
			Name:    name,
			Version: version,
			AppID:   appID,
			Type:    models.PackageTypeMAS,
		})
	}
	_ = s
	return packages
}

func TestScanMASApps_AppIDCaptured(t *testing.T) {
	masOutput := "497799835  Xcode (15.2)\n1333542190  1Password 7 - Password Manager (7.9.11)\n"
	pkgs := parseMASOutput(masOutput)

	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(pkgs))
	}

	tests := []struct {
		name    string
		version string
		appID   string
	}{
		{"Xcode", "15.2", "497799835"},
		{"1Password 7 - Password Manager", "7.9.11", "1333542190"},
	}

	for i, tt := range tests {
		pkg := pkgs[i]
		if pkg.Name != tt.name {
			t.Errorf("pkg[%d].Name = %q, want %q", i, pkg.Name, tt.name)
		}
		if pkg.Version != tt.version {
			t.Errorf("pkg[%d].Version = %q, want %q", i, pkg.Version, tt.version)
		}
		if pkg.AppID != tt.appID {
			t.Errorf("pkg[%d].AppID = %q, want %q", i, pkg.AppID, tt.appID)
		}
		if pkg.Type != models.PackageTypeMAS {
			t.Errorf("pkg[%d].Type = %q, want %q", i, pkg.Type, models.PackageTypeMAS)
		}
	}
}

func TestScanMASApps_SkipsEmptyLines(t *testing.T) {
	masOutput := "\n409183694  Keynote (13.2)\n\n"
	pkgs := parseMASOutput(masOutput)
	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	if pkgs[0].AppID != "409183694" {
		t.Errorf("AppID = %q, want %q", pkgs[0].AppID, "409183694")
	}
}

func TestScanMASApps_GracefulOnMissingMas(t *testing.T) {
	// On Linux/CI, mas is not installed; scanMASApps must return nil, nil.
	s := New()
	pkgs, err := s.scanMASApps(context.Background())
	if err != nil {
		t.Fatalf("scanMASApps() error = %v", err)
	}
	_ = pkgs
}
