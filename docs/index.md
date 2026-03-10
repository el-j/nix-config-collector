---
layout: default
title: nix-config-collector
---

# nix-config-collector

**Collect your macOS setup and generate a full [nix-darwin](https://github.com/LnL7/nix-darwin) + [home-manager](https://github.com/nix-community/home-manager) configuration.**

[![Tests](https://github.com/el-j/nix-config-collector/actions/workflows/test.yml/badge.svg)](https://github.com/el-j/nix-config-collector/actions/workflows/test.yml)
[![Build](https://github.com/el-j/nix-config-collector/actions/workflows/build.yml/badge.svg)](https://github.com/el-j/nix-config-collector/actions/workflows/build.yml)
[![Release](https://img.shields.io/github/v/release/el-j/nix-config-collector?include_prereleases)](https://github.com/el-j/nix-config-collector/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/el-j/nix-config-collector/blob/main/LICENSE)

---

## What it does

`nix-config-collector` scans your macOS system and generates three ready-to-use Nix configuration files:

| Output file | Purpose |
|---|---|
| `darwin-configuration.nix` | nix-darwin system configuration: packages, homebrew, system preferences |
| `home.nix` | home-manager user config: dotfiles, shell, git, user packages |
| `flake.nix` | Nix flake wiring everything together |

---

## Installation

### Pre-built binaries (recommended)

Download the latest release for your architecture:

| Platform | Download |
|---|---|
| macOS (Apple Silicon) | [nix-config-collector-darwin-arm64](https://github.com/el-j/nix-config-collector/releases/latest/download/nix-config-collector-darwin-arm64) |
| macOS (Intel) | [nix-config-collector-darwin-amd64](https://github.com/el-j/nix-config-collector/releases/latest/download/nix-config-collector-darwin-amd64) |
| Linux (amd64) | [nix-config-collector-linux-amd64](https://github.com/el-j/nix-config-collector/releases/latest/download/nix-config-collector-linux-amd64) |

```bash
# Apple Silicon example
curl -L https://github.com/el-j/nix-config-collector/releases/latest/download/nix-config-collector-darwin-arm64 \
  -o /usr/local/bin/nix-config-collector
chmod +x /usr/local/bin/nix-config-collector
```

### Build from source

```bash
git clone https://github.com/el-j/nix-config-collector
cd nix-config-collector
go build -o nix-config-collector ./cmd/cli/
```

---

## Quick start

```bash
# 1. Scan your system and save results
nix-config-collector scan -o scan.json

# 2. Generate Nix config files into ./nix-config/
nix-config-collector generate -o ./nix-config/

# 3. Review and edit the generated files
ls ./nix-config/
# darwin-configuration.nix  home.nix  flake.nix

# 4. Apply with nix-darwin (requires Nix installed)
darwin-rebuild switch --flake ./nix-config#$(hostname)
```

---

## Commands

### `scan`

Scans the macOS system and outputs a JSON inventory.

```
nix-config-collector scan [flags]

Flags:
  -o, --output string   Output file path (default: stdout)
```

**Example output (abbreviated):**
```json
{
  "scanned_at": "2024-01-15T10:30:00Z",
  "hostname": "mymac",
  "packages": [
    {"name": "git", "type": "brew-formula", "nix_name": "git"},
    {"name": "ripgrep", "type": "brew-formula", "nix_name": "ripgrep"},
    {"name": "1Password", "type": "mas"}
  ]
}
```

### `generate`

Scans the system and generates all three Nix config files.

```
nix-config-collector generate [flags]

Flags:
  -o, --output string   Output directory for generated files (default: stdout)
```

### `apply`

Generates configs and applies them with `darwin-rebuild switch`.

```
nix-config-collector apply [flags]

Flags:
  -o, --output string   Output directory (default: ~/.config/nixpkgs)
      --dry-run         Generate configs without applying
```

---

## What gets scanned

| Category | Tool used | Details |
|---|---|---|
| Homebrew formulae | `brew list --formula` | Name, version |
| Homebrew casks | `brew list --cask` | Name, version |
| Mac App Store apps | `mas list` | Name, version |
| System services | `launchctl list` | Name, PID, enabled |
| Dotfiles | filesystem | ~/.zshrc, ~/.gitconfig, ~/.ssh/config, etc. |
| Preferences | `defaults read` | Dock, Finder, keyboard, trackpad |
| Fonts | ~/Library/Fonts, /Library/Fonts | TTF, OTF, TTC |
| Shell | $SHELL | rc files |

---

## Package mapping

Homebrew formulae are automatically mapped to their nixpkgs equivalents. Unmapped packages stay in `homebrew.brews`.

| Homebrew | nixpkgs | Homebrew | nixpkgs |
|---|---|---|---|
| `git` | `git` | `neovim` | `neovim` |
| `ripgrep` | `ripgrep` | `fzf` | `fzf` |
| `node` | `nodejs` | `kubectl` | `kubectl` |
| `go` | `go` | `terraform` | `terraform` |
| `bat` | `bat` | `lazygit` | `lazygit` |
| `starship` | `starship` | `zoxide` | `zoxide` |

[Full mapping table →](https://github.com/el-j/nix-config-collector/blob/main/internal/core/mapper/mapper.go)

---

## Architecture

The project follows **hexagonal architecture** (ports & adapters):

```
cmd/cli/              CLI entry point (Cobra)
cmd/wails/            Desktop app entry point (Wails 2)
internal/
  core/
    models/           Domain types (Package, Service, Dotfile, …)
    ports/            Interface definitions
    mapper/           Brew → Nix package mapping (60+ mappings)
    generator/        Nix config template rendering
    app/              Application service (scan → map → generate → write)
  adapters/
    macos/            macOS system scanner
    fs/               Filesystem writer (with backup)
    cli/              ANSI color CLI notifier
    git/              Git operations adapter
    ai/               Claude AI adapter (package suggestions, config review)
pkg/
  config/             Embedded Nix config templates
```

---

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make your changes with tests
4. Run tests: `go test ./...`
5. Submit a pull request

### Development setup

```bash
git clone https://github.com/el-j/nix-config-collector
cd nix-config-collector
go mod download
go test ./...
go build ./cmd/cli/
```

---

## License

MIT — see [LICENSE](https://github.com/el-j/nix-config-collector/blob/main/LICENSE)
