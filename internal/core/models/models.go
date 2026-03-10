package models

import "time"

// PackageType represents the source of a package
type PackageType string

const (
	PackageTypeBrewFormula PackageType = "brew-formula"
	PackageTypeBrewCask    PackageType = "brew-cask"
	PackageTypeMAS         PackageType = "mas"
	PackageTypeNix         PackageType = "nix"
)

// Package represents an installed software package
type Package struct {
	Name        string      `json:"name"`
	Version     string      `json:"version,omitempty"`
	Type        PackageType `json:"type"`
	AppID       string      `json:"app_id,omitempty"` // Mac App Store numeric ID
	NixName     string      `json:"nix_name,omitempty"`
	Description string      `json:"description,omitempty"`
}

// ServiceType represents the type of system service
type ServiceType string

const (
	ServiceTypeLaunchAgent  ServiceType = "launchagent"
	ServiceTypeLaunchDaemon ServiceType = "launchdaemon"
	ServiceTypeBrewService  ServiceType = "brew-service"
)

// Service represents a system service
type Service struct {
	Name    string      `json:"name"`
	Type    ServiceType `json:"type"`
	Enabled bool        `json:"enabled"`
	PID     int         `json:"pid,omitempty"`
}

// Dotfile represents a configuration file
type Dotfile struct {
	Path        string `json:"path"`
	Content     string `json:"content,omitempty"`
	SymlinkDest string `json:"symlink_dest,omitempty"`
}

// SystemPreference represents a macOS system preference
type SystemPreference struct {
	Domain string      `json:"domain"`
	Key    string      `json:"key"`
	Value  interface{} `json:"value"`
}

// Font represents an installed font
type Font struct {
	Name   string `json:"name"`
	Family string `json:"family,omitempty"`
}

// SystemScan is the complete result of scanning a macOS system
type SystemScan struct {
	ScannedAt   time.Time          `json:"scanned_at"`
	Hostname    string             `json:"hostname"`
	Packages    []Package          `json:"packages"`
	Services    []Service          `json:"services"`
	Dotfiles    []Dotfile          `json:"dotfiles"`
	Preferences []SystemPreference `json:"preferences"`
	Fonts       []Font             `json:"fonts,omitempty"`
	ShellConfig ShellConfig        `json:"shell_config"`
}

// ShellConfig represents the user's shell configuration
type ShellConfig struct {
	Shell   string   `json:"shell"`
	RCFiles []string `json:"rc_files"`
	Path    []string `json:"path,omitempty"`
}

// NixConfig represents the generated Nix configuration
type NixConfig struct {
	DarwinConfig string `json:"darwin_config"`
	HomeConfig   string `json:"home_config"`
	FlakeConfig  string `json:"flake_config,omitempty"`
}
