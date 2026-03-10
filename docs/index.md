# nix-config-collector

A tool to collect macOS system configuration and generate Nix (nix-darwin + home-manager) configuration files.

## Overview

`nix-config-collector` scans your macOS system and generates ready-to-use Nix configuration files:

- **`darwin-configuration.nix`** – nix-darwin system configuration
- **`home.nix`** – home-manager user configuration  
- **`flake.nix`** – Nix flake wrapping everything together

## Installation

### From Source

```bash
git clone https://github.com/el-j/nix-config-collector
cd nix-config-collector
go build -o nix-config-collector ./cmd/cli/
```

### Pre-built Binaries

Download from the [releases page](https://github.com/el-j/nix-config-collector/releases).

## Usage

### Scan your system

```bash
# Output scan results as JSON to stdout
nix-config-collector scan

# Save scan results to a file
nix-config-collector scan -o scan.json
```

### Generate Nix configs

```bash
# Print configs to stdout
nix-config-collector generate

# Write configs to a directory
nix-config-collector generate -o ./nix-config/
```

## What Gets Scanned

| Category | Details |
|----------|---------|
| Packages | Homebrew formulae, casks, Mac App Store apps |
| Services | launchctl agents and daemons |
| Dotfiles | Shell configs, git config, SSH config, and more |
| Preferences | Dock, Finder, keyboard, trackpad settings |
| Fonts | User and system fonts |
| Shell | Current shell, rc files |

## Package Mapping

Homebrew formulae are automatically mapped to their nixpkgs equivalents when possible. For example:

| Homebrew | nixpkgs |
|----------|---------|
| `git` | `git` |
| `neovim` | `neovim` |
| `ripgrep` | `ripgrep` |
| `node` | `nodejs` |
| `kubectl` | `kubectl` |

Packages without a known Nix equivalent are kept in the `homebrew.brews` list.

## Architecture

The project follows hexagonal architecture:

```
cmd/cli/          – CLI entry point (Cobra)
internal/
  core/
    models/       – Domain types
    ports/        – Interface definitions
    mapper/       – Brew→Nix package mapping
    generator/    – Nix config template rendering
    app/          – Application service (orchestration)
  adapters/
    macos/        – macOS system scanner
    fs/           – Filesystem writer
    cli/          – CLI notifier
```

## License

MIT
