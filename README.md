# nix-config-collector

> Collect your macOS setup and generate a full [nix-darwin](https://github.com/LnL7/nix-darwin) + [home-manager](https://github.com/nix-community/home-manager) configuration.

[![Tests](https://github.com/el-j/nix-config-collector/actions/workflows/test.yml/badge.svg)](https://github.com/el-j/nix-config-collector/actions/workflows/test.yml)
[![Build](https://github.com/el-j/nix-config-collector/actions/workflows/build.yml/badge.svg)](https://github.com/el-j/nix-config-collector/actions/workflows/build.yml)
[![Release](https://img.shields.io/github/v/release/el-j/nix-config-collector?include_prereleases)](https://github.com/el-j/nix-config-collector/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/el-j/nix-config-collector)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## What it does

`nix-config-collector` bridges the gap between your existing macOS environment and a reproducible, declarative [Nix](https://nixos.org) setup. It:

1. **Scans** your system — Homebrew packages, Mac App Store apps, system services, dotfiles, and preferences
2. **Maps** Homebrew formulae to their nixpkgs equivalents (60+ built-in mappings)
3. **Generates** three ready-to-use Nix config files:

| File | Purpose |
|---|---|
| `darwin-configuration.nix` | System packages, Homebrew config, system preferences |
| `home.nix` | User packages, dotfiles, shell, git config |
| `flake.nix` | Nix flake tying everything together |

---

## Installation

### Pre-built binaries

```bash
# Apple Silicon (M1/M2/M3)
curl -L https://github.com/el-j/nix-config-collector/releases/latest/download/nix-config-collector-darwin-arm64 \
  -o /usr/local/bin/nix-config-collector && chmod +x /usr/local/bin/nix-config-collector

# Intel Mac
curl -L https://github.com/el-j/nix-config-collector/releases/latest/download/nix-config-collector-darwin-amd64 \
  -o /usr/local/bin/nix-config-collector && chmod +x /usr/local/bin/nix-config-collector

# Linux AMD64
curl -L https://github.com/el-j/nix-config-collector/releases/latest/download/nix-config-collector-linux-amd64 \
  -o /usr/local/bin/nix-config-collector && chmod +x /usr/local/bin/nix-config-collector

# Linux ARM64
curl -L https://github.com/el-j/nix-config-collector/releases/latest/download/nix-config-collector-linux-arm64 \
  -o /usr/local/bin/nix-config-collector && chmod +x /usr/local/bin/nix-config-collector
```

### Build from source

```bash
git clone https://github.com/el-j/nix-config-collector
cd nix-config-collector
make build
# or: go build -o bin/nix-config-collector ./cmd/cli/
```

---

## Quick start

```bash
# Scan your system
nix-config-collector scan -o scan.json

# Generate Nix configs (preview to stdout)
nix-config-collector generate

# Generate and write to a directory
nix-config-collector generate -o ~/.config/nixpkgs/

# Apply immediately with darwin-rebuild
nix-config-collector apply -o ~/.config/nixpkgs/
```

---

## Commands

### `scan`
Scans the macOS system for packages, services, dotfiles, and preferences. Outputs a JSON inventory.

```
Usage:
  nix-config-collector scan [flags]

Flags:
  -o, --output string   Write JSON to file instead of stdout
```

### `generate`
Scans the system and generates `darwin-configuration.nix`, `home.nix`, and `flake.nix`.

```
Usage:
  nix-config-collector generate [flags]

Flags:
  -o, --output string   Output directory (default: print to stdout)
```

### `apply`
Generates configs and applies them with `darwin-rebuild switch`.

```
Usage:
  nix-config-collector apply [flags]

Flags:
  -o, --output string   Output directory (default: ~/.config/nixpkgs)
      --dry-run         Generate configs without applying
```

### `version`
Prints the version string.

---

## Architecture

The project follows **hexagonal architecture** (ports & adapters pattern):

```
cmd/
  cli/              CLI entry point (Cobra commands)
  wails/            Desktop app entry point (Wails 2, build with -tags wails)
internal/
  core/
    models/         Domain types: Package, Service, Dotfile, SystemPreference, …
    ports/          Interfaces: SystemScanner, FileWriter, UserNotifier, ConfigGenerator, ConfigMapper
    mapper/         Homebrew → nixpkgs mapping (60+ entries)
    generator/      Nix config template rendering (Go text/template)
    app/            Application service: orchestrates scan → map → generate → write
  adapters/
    macos/          macOS scanner: brew, mas, launchctl, defaults, fonts
    fs/             Filesystem writer with timestamped backups
    cli/            ANSI color CLI notifier
    git/            Git adapter: auto-commit generated configs
    ai/             Claude Sonnet AI adapter: package suggestions, config review
pkg/
  config/           Embedded Nix config templates (go:embed)
frontend/           Wails React frontend (scan dashboard, config preview)
```

### Data flow

```
macOS system
    │
    ▼
macos.Scanner (ScanPackages / ScanServices / ScanDotfiles / ScanPreferences)
    │
    ▼
mapper.Mapper  ──→  Package.NixName populated for known formulae
    │
    ▼
generator.Generator  ──→  darwin-configuration.nix / home.nix / flake.nix
    │
    ▼
fs.Writer  ──→  writes to output directory (with backup)
    │
    ▼
git.Adapter  ──→  optional: auto-commit to git
```

---

## What gets scanned

| Category | Source | Details |
|---|---|---|
| Homebrew formulae | `brew list --formula --versions` | Name, version, mapped to nixpkgs |
| Homebrew casks | `brew list --cask --versions` | Name, version |
| Mac App Store | `mas list` | Name, version |
| System services | `launchctl list` | Name, PID, enabled/disabled |
| Dotfiles | filesystem | ~/.zshrc, ~/.gitconfig, ~/.ssh/config, and more |
| Preferences | `defaults read` | Dock, Finder, keyboard, trackpad |
| Fonts | ~/Library/Fonts, /Library/Fonts | TTF, OTF, TTC files |
| Shell | $SHELL, rc files | Shell name and rc file paths |

---

## Package mapping

Homebrew formulae are automatically mapped to their nixpkgs equivalents. Packages without a known mapping are placed in `homebrew.brews` so they continue to be managed by Homebrew.

```
brew install git ripgrep neovim node
            ↓
environment.systemPackages = with pkgs; [ git ripgrep neovim nodejs ];
```

See [`internal/core/mapper/mapper.go`](internal/core/mapper/mapper.go) for the full table.

---

## Contributing

Contributions welcome! Please:

1. Fork the repository and create a feature branch
2. Write tests for new functionality
3. Run `go test ./...` and `go vet ./...`
4. Submit a pull request

### Development commands

```bash
make build      # Build CLI binary
make test       # Run all tests
make lint       # Run golangci-lint
make build-all  # Cross-compile for all platforms
make clean      # Remove build artifacts
```

---

## License

MIT — see [LICENSE](LICENSE)
