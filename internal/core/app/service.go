package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/el-j/nix-config-collector/internal/core/mapper"
	"github.com/el-j/nix-config-collector/internal/core/models"
	"github.com/el-j/nix-config-collector/internal/core/ports"
)

// Service orchestrates the scan → map → generate → write workflow
type Service struct {
	scanner   ports.SystemScanner
	writer    ports.FileWriter
	notifier  ports.UserNotifier
	generator ports.ConfigGenerator
	mapper    ports.ConfigMapper
}

// NewService creates a new application service
func NewService(
	scanner ports.SystemScanner,
	writer ports.FileWriter,
	notifier ports.UserNotifier,
	generator ports.ConfigGenerator,
) *Service {
	return &Service{
		scanner:   scanner,
		writer:    writer,
		notifier:  notifier,
		generator: generator,
		mapper:    mapper.New(),
	}
}

// ScanOptions configures the scan behavior
type ScanOptions struct {
	DotfilePaths []string
	Timeout      time.Duration
}

// GenerateOptions configures the generation behavior
type GenerateOptions struct {
	OutputDir   string
	ScanOptions ScanOptions
}

// Scan performs a full system scan and returns the result
func (s *Service) Scan(ctx context.Context, opts ScanOptions) (*models.SystemScan, error) {
	scan := &models.SystemScan{
		ScannedAt: time.Now(),
	}

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	scan.Hostname = hostname

	s.notifier.Progress(1, 6, "Scanning packages...")
	packages, err := s.scanner.ScanPackages(ctx)
	if err != nil {
		return nil, fmt.Errorf("scanning packages: %w", err)
	}

	// Map packages to Nix equivalents
	mappedPackages, err := s.mapper.MapPackages(packages)
	if err != nil {
		return nil, fmt.Errorf("mapping packages: %w", err)
	}
	scan.Packages = mappedPackages

	s.notifier.Progress(2, 6, "Scanning services...")
	services, err := s.scanner.ScanServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("scanning services: %w", err)
	}
	scan.Services = services

	s.notifier.Progress(3, 6, "Scanning dotfiles...")
	dotfiles, err := s.scanner.ScanDotfiles(ctx, opts.DotfilePaths)
	if err != nil {
		return nil, fmt.Errorf("scanning dotfiles: %w", err)
	}
	scan.Dotfiles = dotfiles

	s.notifier.Progress(4, 6, "Scanning preferences...")
	prefs, err := s.scanner.ScanPreferences(ctx)
	if err != nil {
		return nil, fmt.Errorf("scanning preferences: %w", err)
	}
	scan.Preferences = prefs

	s.notifier.Progress(5, 6, "Scanning fonts...")
	fonts, err := s.scanner.ScanFonts(ctx)
	if err != nil {
		return nil, fmt.Errorf("scanning fonts: %w", err)
	}
	scan.Fonts = fonts

	s.notifier.Progress(6, 6, "Scanning shell config...")
	shellConfig, err := s.scanner.ScanShellConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("scanning shell config: %w", err)
	}
	scan.ShellConfig = shellConfig

	return scan, nil
}

// Generate performs a full scan and generates Nix configuration files
func (s *Service) Generate(ctx context.Context, opts GenerateOptions) (*models.NixConfig, error) {
	scan, err := s.Scan(ctx, opts.ScanOptions)
	if err != nil {
		return nil, fmt.Errorf("scanning system: %w", err)
	}

	s.notifier.Info("Generating Nix configurations...")

	darwinConfig, err := s.generator.GenerateDarwinConfig(*scan)
	if err != nil {
		return nil, fmt.Errorf("generating darwin config: %w", err)
	}

	homeConfig, err := s.generator.GenerateHomeConfig(*scan)
	if err != nil {
		return nil, fmt.Errorf("generating home config: %w", err)
	}

	flakeConfig, err := s.generator.GenerateFlakeConfig(*scan)
	if err != nil {
		return nil, fmt.Errorf("generating flake config: %w", err)
	}

	config := &models.NixConfig{
		DarwinConfig: darwinConfig,
		HomeConfig:   homeConfig,
		FlakeConfig:  flakeConfig,
	}

	if opts.OutputDir != "" {
		if err := s.writer.EnsureDir(opts.OutputDir); err != nil {
			return nil, fmt.Errorf("creating output directory: %w", err)
		}

		files := map[string]string{
			"darwin-configuration.nix": darwinConfig,
			"home.nix":                 homeConfig,
			"flake.nix":                flakeConfig,
		}

		for filename, content := range files {
			path := opts.OutputDir + "/" + filename
			if s.writer.FileExists(path) {
				if err := s.writer.BackupFile(path); err != nil {
					s.notifier.Warning(fmt.Sprintf("Could not backup %s: %v", filename, err))
				}
			}
			if err := s.writer.WriteFile(path, content); err != nil {
				return nil, fmt.Errorf("writing %s: %w", filename, err)
			}
			s.notifier.Success(fmt.Sprintf("Written: %s", path))
		}
	}

	return config, nil
}
