---
layout: default
title: Changelog
---

# Changelog

All notable changes to this project will be documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added
- `pkg/config/` embedded Nix config templates with `go:embed`
- `internal/adapters/git/` adapter for auto-committing generated configs
- `internal/adapters/ai/` Claude Sonnet AI adapter for package suggestions and config review
- `internal/adapters/macos/` unit tests covering all scanner methods
- `internal/core/app/` integration tests with mock adapters
- Wails 2 desktop app scaffold (`cmd/wails/`, `frontend/`)
- `apply` CLI command with `--dry-run` flag
- `.gitignore` file
- Homebrew tap formula
- Enriched GitHub Pages documentation site

## [0.1.0] - 2026-03-10

### Added
- Initial release
- Core domain models (Package, Service, Dotfile, SystemPreference, Font, SystemScan, NixConfig)
- Ports (interfaces) for all external interactions
- Homebrew → nixpkgs mapping engine (60+ mappings)
- Nix config generator (darwin-configuration.nix, home.nix, flake.nix)
- macOS system scanner: `brew list`, `mas list`, `launchctl list`, `defaults read`
- CLI commands: `scan`, `generate`, `version`
- GitHub Actions: test, build (multi-platform), release, docs
- GitHub Pages documentation site
- JSON-based Claude orchestration system (`.claude/orchestrator.json`)
