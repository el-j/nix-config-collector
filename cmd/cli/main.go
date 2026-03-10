package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/el-j/nix-config-collector/internal/adapters/cli"
	fsadapter "github.com/el-j/nix-config-collector/internal/adapters/fs"
	"github.com/el-j/nix-config-collector/internal/adapters/macos"
	"github.com/el-j/nix-config-collector/internal/core/app"
	"github.com/el-j/nix-config-collector/internal/core/generator"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "nix-config-collector",
		Short: "Collect macOS configuration and generate Nix configs",
		Long: `nix-config-collector scans your macOS system and generates
nix-darwin + home-manager configuration files.`,
	}

	rootCmd.AddCommand(scanCmd())
	rootCmd.AddCommand(generateCmd())
	rootCmd.AddCommand(applyCmd())
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func scanCmd() *cobra.Command {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan system and output JSON",
		Long:  "Scans the macOS system for packages, services, dotfiles, and preferences, outputs JSON.",
		RunE: func(cmd *cobra.Command, args []string) error {
			notifier := cli.New()
			scanner := macos.New()
			writer := fsadapter.New()
			gen := generator.New()

			svc := app.NewService(scanner, writer, notifier, gen)

			ctx := context.Background()
			scan, err := svc.Scan(ctx, app.ScanOptions{})
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			data, err := json.MarshalIndent(scan, "", "  ")
			if err != nil {
				return fmt.Errorf("marshaling JSON: %w", err)
			}

			if outputFile != "" {
				if err := os.WriteFile(outputFile, data, 0640); err != nil {
					return fmt.Errorf("writing output file: %w", err)
				}
				notifier.Success(fmt.Sprintf("Scan results written to %s", outputFile))
			} else {
				fmt.Println(string(data))
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: stdout)")
	return cmd
}

func generateCmd() *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Scan and generate Nix configs",
		Long:  "Scans the macOS system and generates nix-darwin + home-manager configuration files.",
		RunE: func(cmd *cobra.Command, args []string) error {
			notifier := cli.New()
			scanner := macos.New()
			writer := fsadapter.New()
			gen := generator.New()

			svc := app.NewService(scanner, writer, notifier, gen)

			ctx := context.Background()
			config, err := svc.Generate(ctx, app.GenerateOptions{
				OutputDir: outputDir,
			})
			if err != nil {
				return fmt.Errorf("generation failed: %w", err)
			}

			if outputDir == "" {
				fmt.Println("=== darwin-configuration.nix ===")
				fmt.Println(config.DarwinConfig)
				fmt.Println("=== home.nix ===")
				fmt.Println(config.HomeConfig)
				fmt.Println("=== flake.nix ===")
				fmt.Println(config.FlakeConfig)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for generated files")
	return cmd
}

func applyCmd() *cobra.Command {
	var outputDir string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Generate Nix configs and apply with darwin-rebuild",
		Long: `Scans the macOS system, generates nix-darwin + home-manager configuration,
and applies it using 'darwin-rebuild switch'.

Requires Nix and nix-darwin to be installed.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			notifier := cli.New()

			// Expand ~ in output dir
			if outputDir == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("getting home directory: %w", err)
				}
				outputDir = filepath.Join(home, ".config", "nixpkgs")
			} else if len(outputDir) >= 2 && outputDir[:2] == "~/" {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("getting home directory: %w", err)
				}
				outputDir = filepath.Join(home, outputDir[2:])
			}

			scanner := macos.New()
			writer := fsadapter.New()
			gen := generator.New()
			svc := app.NewService(scanner, writer, notifier, gen)

			ctx := context.Background()

			notifier.Info(fmt.Sprintf("Generating Nix configs into %s ...", outputDir))
			_, err := svc.Generate(ctx, app.GenerateOptions{OutputDir: outputDir})
			if err != nil {
				return fmt.Errorf("generation failed: %w", err)
			}
			notifier.Success("Nix configs generated successfully")

			if dryRun {
				notifier.Info("Dry-run mode: skipping darwin-rebuild")
				return nil
			}

			hostname, err := os.Hostname()
			if err != nil {
				hostname = "default"
			}

			flakeArg := fmt.Sprintf("%s#%s", outputDir, hostname)
			notifier.Info(fmt.Sprintf("Running: darwin-rebuild switch --flake %s", flakeArg))

			darwinCmd := exec.CommandContext(ctx, "darwin-rebuild", "switch", "--flake", flakeArg)
			darwinCmd.Stdout = os.Stdout
			darwinCmd.Stderr = os.Stderr
			if err := darwinCmd.Run(); err != nil {
				return fmt.Errorf("darwin-rebuild failed: %w", err)
			}
			notifier.Success("Configuration applied successfully!")
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory (default: ~/.config/nixpkgs)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Generate configs without applying")
	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("nix-config-collector version %s\n", version)
		},
	}
}
