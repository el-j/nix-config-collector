package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("nix-config-collector version %s\n", version)
		},
	}
}
