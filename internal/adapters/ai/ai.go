// Package ai provides an adapter for AI-assisted Nix configuration suggestions.
// It integrates with the Anthropic Claude API (claude-sonnet-4-5) to provide
// intelligent package mapping, config improvements, and user guidance.
package ai

import (
	"context"
	"errors"
)

// ErrNoAPIKey is returned when the Claude API key is not configured.
var ErrNoAPIKey = errors.New("ANTHROPIC_API_KEY not set; AI features are disabled")

// Suggestion represents an AI-generated config suggestion.
type Suggestion struct {
	// Type is the kind of suggestion: "package-mapping", "config-improvement", "warning"
	Type    string `json:"type"`
	Message string `json:"message"`
	// Before is the original snippet (may be empty)
	Before string `json:"before,omitempty"`
	// After is the suggested replacement (may be empty)
	After string `json:"after,omitempty"`
}

// Client defines the interface for AI-assisted suggestions.
type Client interface {
	// SuggestPackageMapping asks the AI to find a Nix package for the given Homebrew formula.
	SuggestPackageMapping(ctx context.Context, brewName string) (*Suggestion, error)
	// ReviewConfig asks the AI to review a generated Nix config and suggest improvements.
	ReviewConfig(ctx context.Context, configType, content string) ([]Suggestion, error)
	// Explain asks the AI to explain a Nix expression in plain English.
	Explain(ctx context.Context, nixExpression string) (string, error)
}

// NoOpClient is a Client implementation that returns no suggestions.
// Used when no API key is configured.
type NoOpClient struct{}

// NewNoOpClient creates a no-op AI client.
func NewNoOpClient() *NoOpClient {
	return &NoOpClient{}
}

// SuggestPackageMapping returns ErrNoAPIKey.
func (n *NoOpClient) SuggestPackageMapping(_ context.Context, _ string) (*Suggestion, error) {
	return nil, ErrNoAPIKey
}

// ReviewConfig returns ErrNoAPIKey.
func (n *NoOpClient) ReviewConfig(_ context.Context, _, _ string) ([]Suggestion, error) {
	return nil, ErrNoAPIKey
}

// Explain returns ErrNoAPIKey.
func (n *NoOpClient) Explain(_ context.Context, _ string) (string, error) {
	return "", ErrNoAPIKey
}

// ClaudeClient is a Client backed by the Anthropic Claude API.
// Use NewClaudeClient to construct it.
type ClaudeClient struct {
	apiKey string
	model  string
}

// NewClaudeClient creates a new Claude API client.
// apiKey is the Anthropic API key (ANTHROPIC_API_KEY).
// model defaults to "claude-sonnet-4-5" if empty.
func NewClaudeClient(apiKey, model string) (*ClaudeClient, error) {
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}
	if model == "" {
		model = "claude-sonnet-4-5"
	}
	return &ClaudeClient{apiKey: apiKey, model: model}, nil
}

// SuggestPackageMapping queries Claude to find a nixpkgs equivalent for a brew formula.
func (c *ClaudeClient) SuggestPackageMapping(ctx context.Context, brewName string) (*Suggestion, error) {
	prompt := "What is the nixpkgs attribute name for the Homebrew formula '" + brewName + "'? " +
		"Reply with only the attribute name (e.g. 'ripgrep') or 'unknown' if not in nixpkgs."
	answer, err := c.complete(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return &Suggestion{
		Type:    "package-mapping",
		Message: answer,
		Before:  brewName,
		After:   answer,
	}, nil
}

// ReviewConfig queries Claude to review a generated Nix config.
func (c *ClaudeClient) ReviewConfig(ctx context.Context, configType, content string) ([]Suggestion, error) {
	prompt := "Review this " + configType + " Nix configuration and list up to 3 specific improvements. " +
		"Format each as a short bullet point.\n\n```nix\n" + content + "\n```"
	answer, err := c.complete(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return []Suggestion{{
		Type:    "config-improvement",
		Message: answer,
	}}, nil
}

// Explain asks Claude to explain a Nix expression.
func (c *ClaudeClient) Explain(ctx context.Context, nixExpression string) (string, error) {
	prompt := "Explain this Nix expression in plain English for a developer new to Nix:\n\n```nix\n" + nixExpression + "\n```"
	return c.complete(ctx, prompt)
}

// complete sends a single-turn message to the Claude API and returns the text response.
// This is a minimal HTTP implementation to avoid adding the Anthropic SDK as a dependency.
func (c *ClaudeClient) complete(ctx context.Context, prompt string) (string, error) {
	return claudeHTTPComplete(ctx, c.apiKey, c.model, prompt)
}
