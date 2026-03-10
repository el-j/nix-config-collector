package config

import "embed"

// FS holds the embedded template files.
//
//go:embed templates/*
var FS embed.FS

// TemplatePath constants for the embedded templates.
const (
	DarwinConfigTemplate = "templates/darwin-configuration.nix.tmpl"
	HomeConfigTemplate   = "templates/home.nix.tmpl"
	FlakeConfigTemplate  = "templates/flake.nix.tmpl"
)

// ReadTemplate reads an embedded template by its path constant.
func ReadTemplate(path string) (string, error) {
	data, err := FS.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
