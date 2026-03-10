package ports

import (
	"context"

	"github.com/el-j/nix-config-collector/internal/core/models"
)

// SystemScanner defines the interface for scanning a macOS system
type SystemScanner interface {
	ScanPackages(ctx context.Context) ([]models.Package, error)
	ScanServices(ctx context.Context) ([]models.Service, error)
	ScanDotfiles(ctx context.Context, paths []string) ([]models.Dotfile, error)
	ScanPreferences(ctx context.Context) ([]models.SystemPreference, error)
	ScanFonts(ctx context.Context) ([]models.Font, error)
	ScanShellConfig(ctx context.Context) (models.ShellConfig, error)
}

// FileWriter defines the interface for writing files
type FileWriter interface {
	WriteFile(path string, content string) error
	EnsureDir(path string) error
	BackupFile(path string) error
	FileExists(path string) bool
}

// UserNotifier defines the interface for user interaction
type UserNotifier interface {
	Info(msg string)
	Success(msg string)
	Warning(msg string)
	Error(msg string)
	Progress(current, total int, msg string)
	Confirm(msg string) bool
}

// ConfigMapper defines the interface for mapping scanned data to Nix config
type ConfigMapper interface {
	MapPackages(packages []models.Package) ([]models.Package, error)
	MapServices(services []models.Service) ([]models.Service, error)
}

// ConfigGenerator defines the interface for generating Nix config files
type ConfigGenerator interface {
	GenerateDarwinConfig(scan models.SystemScan) (string, error)
	GenerateHomeConfig(scan models.SystemScan) (string, error)
	GenerateFlakeConfig(scan models.SystemScan) (string, error)
}
