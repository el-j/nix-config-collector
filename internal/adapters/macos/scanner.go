package macos

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/el-j/nix-config-collector/internal/core/models"
)

// Scanner implements ports.SystemScanner for macOS
type Scanner struct {
	homeDir string
	timeout time.Duration
}

// New creates a new macOS Scanner
func New() *Scanner {
	homeDir, _ := os.UserHomeDir()
	return &Scanner{
		homeDir: homeDir,
		timeout: 30 * time.Second,
	}
}

// WithTimeout sets the command execution timeout
func (s *Scanner) WithTimeout(d time.Duration) *Scanner {
	s.timeout = d
	return s
}

// runCommand runs a shell command with context timeout
func (s *Scanner) runCommand(ctx context.Context, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("command timed out: %s %v", name, args)
		}
		// Return empty string for non-critical command failures (e.g., command not found)
		return "", nil
	}
	return stdout.String(), nil
}

// ScanPackages scans all installed packages
func (s *Scanner) ScanPackages(ctx context.Context) ([]models.Package, error) {
	var packages []models.Package

	// Homebrew formulae
	brewFormulae, err := s.scanBrewFormulae(ctx)
	if err != nil {
		return nil, fmt.Errorf("scanning brew formulae: %w", err)
	}
	packages = append(packages, brewFormulae...)

	// Homebrew casks
	brewCasks, err := s.scanBrewCasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("scanning brew casks: %w", err)
	}
	packages = append(packages, brewCasks...)

	// Mac App Store
	masApps, err := s.scanMASApps(ctx)
	if err != nil {
		return nil, fmt.Errorf("scanning MAS apps: %w", err)
	}
	packages = append(packages, masApps...)

	return packages, nil
}

func (s *Scanner) scanBrewFormulae(ctx context.Context) ([]models.Package, error) {
	output, err := s.runCommand(ctx, "brew", "list", "--formula", "--versions")
	if err != nil || output == "" {
		return nil, err
	}

	var packages []models.Package
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		pkg := models.Package{
			Name: parts[0],
			Type: models.PackageTypeBrewFormula,
		}
		if len(parts) > 1 {
			pkg.Version = parts[1]
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

func (s *Scanner) scanBrewCasks(ctx context.Context) ([]models.Package, error) {
	output, err := s.runCommand(ctx, "brew", "list", "--cask", "--versions")
	if err != nil || output == "" {
		return nil, err
	}

	var packages []models.Package
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		pkg := models.Package{
			Name: parts[0],
			Type: models.PackageTypeBrewCask,
		}
		if len(parts) > 1 {
			pkg.Version = parts[1]
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

func (s *Scanner) scanMASApps(ctx context.Context) ([]models.Package, error) {
	output, err := s.runCommand(ctx, "mas", "list")
	if err != nil || output == "" {
		return nil, err
	}

	var packages []models.Package
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// MAS format: "id  Name (version)"
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}
		name := strings.TrimSpace(parts[1])
		version := ""
		if idx := strings.LastIndex(name, " ("); idx != -1 {
			version = strings.Trim(name[idx+2:], ")")
			name = name[:idx]
		}
		packages = append(packages, models.Package{
			Name:    name,
			Version: version,
			Type:    models.PackageTypeMAS,
		})
	}
	return packages, nil
}

// ScanServices scans all system services
func (s *Scanner) ScanServices(ctx context.Context) ([]models.Service, error) {
	var services []models.Service

	launchServices, err := s.scanLaunchctlServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("scanning launchctl services: %w", err)
	}
	services = append(services, launchServices...)

	return services, nil
}

func (s *Scanner) scanLaunchctlServices(ctx context.Context) ([]models.Service, error) {
	output, err := s.runCommand(ctx, "launchctl", "list")
	if err != nil || output == "" {
		return nil, err
	}

	var services []models.Service
	scanner := bufio.NewScanner(strings.NewReader(output))
	first := true
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if first { // skip header
			first = false
			continue
		}
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}
		pid, _ := strconv.Atoi(parts[0])
		name := parts[2]
		if name == "-" || strings.HasPrefix(name, "0x") {
			continue
		}
		services = append(services, models.Service{
			Name:    name,
			Type:    models.ServiceTypeLaunchAgent,
			Enabled: pid > 0,
			PID:     pid,
		})
	}
	return services, nil
}

// ScanDotfiles scans configuration files
func (s *Scanner) ScanDotfiles(ctx context.Context, paths []string) ([]models.Dotfile, error) {
	if len(paths) == 0 {
		paths = s.defaultDotfilePaths()
	}

	var dotfiles []models.Dotfile
	for _, p := range paths {
		expanded := s.expandPath(p)
		info, err := os.Stat(expanded)
		if err != nil {
			continue // file doesn't exist, skip
		}

		if info.IsDir() {
			dirFiles, err := s.scanDotfileDir(expanded)
			if err != nil {
				continue
			}
			dotfiles = append(dotfiles, dirFiles...)
		} else {
			content, err := os.ReadFile(expanded)
			if err != nil {
				continue
			}
			dotfiles = append(dotfiles, models.Dotfile{
				Path:    p,
				Content: string(content),
			})
		}
	}
	return dotfiles, nil
}

func (s *Scanner) scanDotfileDir(dir string) ([]models.Dotfile, error) {
	var dotfiles []models.Dotfile
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		// Skip binary files and large files
		info, err := d.Info()
		if err != nil || info.Size() > 1024*1024 { // skip files > 1MB
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		// Skip binary content
		if !isText(content) {
			return nil
		}
		rel, _ := filepath.Rel(s.homeDir, path)
		dotfiles = append(dotfiles, models.Dotfile{
			Path:    "~/" + rel,
			Content: string(content),
		})
		return nil
	})
	return dotfiles, err
}

func isText(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return false
		}
	}
	return true
}

func (s *Scanner) defaultDotfilePaths() []string {
	return []string{
		"~/.zshrc",
		"~/.bashrc",
		"~/.bash_profile",
		"~/.profile",
		"~/.zprofile",
		"~/.gitconfig",
		"~/.ssh/config",
		"~/.config/starship.toml",
		"~/.config/fish/config.fish",
		"~/.config/nvim/init.lua",
		"~/.config/nvim/init.vim",
	}
}

func (s *Scanner) expandPath(p string) string {
	if strings.HasPrefix(p, "~/") {
		return filepath.Join(s.homeDir, p[2:])
	}
	return p
}

// ScanPreferences scans macOS system preferences
func (s *Scanner) ScanPreferences(ctx context.Context) ([]models.SystemPreference, error) {
	domains := []struct {
		domain string
		keys   []string
	}{
		{"com.apple.dock", []string{"autohide", "tilesize", "orientation", "show-recents"}},
		{"com.apple.finder", []string{"ShowExternalHardDrivesOnDesktop", "ShowHardDrivesOnDesktop", "AppleShowAllFiles"}},
		{"NSGlobalDomain", []string{"AppleShowAllExtensions", "AppleInterfaceStyle", "KeyRepeat", "InitialKeyRepeat"}},
		{"com.apple.trackpad", []string{"Clicking", "TrackpadRightClick"}},
	}

	var prefs []models.SystemPreference
	for _, d := range domains {
		for _, key := range d.keys {
			output, err := s.runCommand(ctx, "defaults", "read", d.domain, key)
			if err != nil || output == "" {
				continue
			}
			prefs = append(prefs, models.SystemPreference{
				Domain: d.domain,
				Key:    key,
				Value:  strings.TrimSpace(output),
			})
		}
	}
	return prefs, nil
}

// ScanFonts scans installed fonts (macOS specific locations)
func (s *Scanner) ScanFonts(ctx context.Context) ([]models.Font, error) {
	fontDirs := []string{
		"~/Library/Fonts",
		"/Library/Fonts",
		"/System/Library/Fonts",
	}

	var fonts []models.Font
	seen := make(map[string]bool)

	for _, dir := range fontDirs {
		expanded := s.expandPath(dir)
		entries, err := os.ReadDir(expanded)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			ext := strings.ToLower(filepath.Ext(name))
			if ext != ".ttf" && ext != ".otf" && ext != ".ttc" {
				continue
			}
			family := strings.TrimSuffix(name, filepath.Ext(name))
			family = strings.ReplaceAll(family, "-", " ")
			if !seen[family] {
				seen[family] = true
				fonts = append(fonts, models.Font{
					Name:   name,
					Family: family,
				})
			}
		}
	}
	return fonts, nil
}

// ScanShellConfig scans the user's shell configuration
func (s *Scanner) ScanShellConfig(ctx context.Context) (models.ShellConfig, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/zsh"
	}

	var rcFiles []string
	candidates := []string{"~/.zshrc", "~/.zprofile", "~/.bashrc", "~/.bash_profile", "~/.profile"}
	for _, f := range candidates {
		expanded := s.expandPath(f)
		if _, err := os.Stat(expanded); err == nil {
			rcFiles = append(rcFiles, f)
		}
	}

	return models.ShellConfig{
		Shell:   shell,
		RCFiles: rcFiles,
	}, nil
}
