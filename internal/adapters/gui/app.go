//go:build wails || bindings || desktop

package gui

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/el-j/nix-config-collector/internal/core/app"
	"github.com/el-j/nix-config-collector/internal/core/models"
)

// App holds the Wails application backend methods exposed to the frontend.
type App struct {
	ctx     context.Context
	service *app.Service
}

// NewApp creates a new App.
func NewApp(service *app.Service) *App {
	return &App{service: service}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Scan runs a full system scan and returns the results as JSON.
func (a *App) Scan() (string, error) {
	scan, err := a.service.Scan(a.ctx, app.ScanOptions{})
	if err != nil {
		return "", fmt.Errorf("scan failed: %w", err)
	}
	data, err := json.MarshalIndent(scan, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling scan: %w", err)
	}
	return string(data), nil
}

// Generate runs a full scan and generates Nix config files.
// outputDir is the directory where files will be written.
func (a *App) Generate(outputDir string) (*models.NixConfig, error) {
	config, err := a.service.Generate(a.ctx, app.GenerateOptions{
		OutputDir: outputDir,
	})
	if err != nil {
		return nil, fmt.Errorf("generate failed: %w", err)
	}
	return config, nil
}

// Preview generates configs without writing to disk and returns them as strings.
func (a *App) Preview() (*models.NixConfig, error) {
	config, err := a.service.Generate(a.ctx, app.GenerateOptions{})
	if err != nil {
		return nil, fmt.Errorf("preview failed: %w", err)
	}
	return config, nil
}
