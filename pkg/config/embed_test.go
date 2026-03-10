package config_test

import (
	"strings"
	"testing"

	"github.com/el-j/nix-config-collector/pkg/config"
)

func TestEmbeddedTemplates(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		contains string
	}{
		{"darwin", config.DarwinConfigTemplate, "darwin-configuration.nix"},
		{"home", config.HomeConfigTemplate, "home.nix"},
		{"flake", config.FlakeConfigTemplate, "flake.nix"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := config.ReadTemplate(tt.path)
			if err != nil {
				t.Fatalf("ReadTemplate(%q) error: %v", tt.path, err)
			}
			if !strings.Contains(content, tt.contains) {
				t.Errorf("template %q does not contain %q", tt.path, tt.contains)
			}
			if len(content) < 100 {
				t.Errorf("template %q seems too short (%d bytes)", tt.path, len(content))
			}
		})
	}
}

func TestFS(t *testing.T) {
	entries, err := config.FS.ReadDir("templates")
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}
	if len(entries) < 3 {
		t.Errorf("expected at least 3 template files, got %d", len(entries))
	}
}
