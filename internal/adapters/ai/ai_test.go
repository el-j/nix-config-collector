package ai_test

import (
	"context"
	"errors"
	"testing"

	"github.com/el-j/nix-config-collector/internal/adapters/ai"
)

func TestNewNoOpClient(t *testing.T) {
	c := ai.NewNoOpClient()
	if c == nil {
		t.Fatal("NewNoOpClient() returned nil")
	}
}

func TestNoOpClient_SuggestPackageMapping(t *testing.T) {
	c := ai.NewNoOpClient()
	_, err := c.SuggestPackageMapping(context.Background(), "ripgrep")
	if !errors.Is(err, ai.ErrNoAPIKey) {
		t.Errorf("expected ErrNoAPIKey, got %v", err)
	}
}

func TestNoOpClient_ReviewConfig(t *testing.T) {
	c := ai.NewNoOpClient()
	_, err := c.ReviewConfig(context.Background(), "darwin", "# config")
	if !errors.Is(err, ai.ErrNoAPIKey) {
		t.Errorf("expected ErrNoAPIKey, got %v", err)
	}
}

func TestNoOpClient_Explain(t *testing.T) {
	c := ai.NewNoOpClient()
	_, err := c.Explain(context.Background(), "pkgs.ripgrep")
	if !errors.Is(err, ai.ErrNoAPIKey) {
		t.Errorf("expected ErrNoAPIKey, got %v", err)
	}
}

func TestNewClaudeClient_NoAPIKey(t *testing.T) {
	_, err := ai.NewClaudeClient("", "")
	if !errors.Is(err, ai.ErrNoAPIKey) {
		t.Errorf("expected ErrNoAPIKey for empty key, got %v", err)
	}
}

func TestNewClaudeClient_WithKey(t *testing.T) {
	c, err := ai.NewClaudeClient("sk-test-key", "")
	if err != nil {
		t.Fatalf("NewClaudeClient() error: %v", err)
	}
	if c == nil {
		t.Fatal("NewClaudeClient() returned nil")
	}
}

func TestSuggestion_Struct(t *testing.T) {
	s := ai.Suggestion{
		Type:    "package-mapping",
		Message: "use ripgrep",
		Before:  "rg",
		After:   "ripgrep",
	}
	if s.Type != "package-mapping" {
		t.Errorf("Type = %q, want package-mapping", s.Type)
	}
}
