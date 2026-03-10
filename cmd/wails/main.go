//go:build wails

package main

import (
	"embed"

	"github.com/el-j/nix-config-collector/internal/adapters/cli"
	fsadapter "github.com/el-j/nix-config-collector/internal/adapters/fs"
	"github.com/el-j/nix-config-collector/internal/adapters/gui"
	"github.com/el-j/nix-config-collector/internal/adapters/macos"
	"github.com/el-j/nix-config-collector/internal/core/app"
	"github.com/el-j/nix-config-collector/internal/core/generator"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	scanner := macos.New()
	writer := fsadapter.New()
	notifier := cli.New()
	gen := generator.New()

	svc := app.NewService(scanner, writer, notifier, gen)
	guiApp := gui.NewApp(svc)

	if err := wails.Run(&options.App{
		Title:  "nix-config-collector",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:        guiApp.Startup,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind: []interface{}{
			guiApp,
		},
	}); err != nil {
		panic(err)
	}
}
