.PHONY: build build-wails test lint clean install build-all release-dry-run

# Version — override with: make build VERSION=1.2.3
VERSION ?= dev

# Default target: build CLI
build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/nix-config-collector ./cmd/cli/

# Build Wails desktop app as a universal binary (requires Wails CLI and macOS)
build-wails:
	wails build -platform darwin/universal

# Run all tests
test:
	go test -v -race ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -rf bin/ dist/ frontend/dist/

# Install CLI binary
install: build
	cp bin/nix-config-collector /usr/local/bin/

# Cross-compile for all platforms
build-all:
	mkdir -p dist
	GOOS=darwin  GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o dist/nix-config-collector-darwin-arm64  ./cmd/cli/
	GOOS=darwin  GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o dist/nix-config-collector-darwin-amd64  ./cmd/cli/
	GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o dist/nix-config-collector-linux-amd64   ./cmd/cli/
	GOOS=linux   GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o dist/nix-config-collector-linux-arm64   ./cmd/cli/

# Dry-run the release: build all platforms locally and generate checksums
release-dry-run: build-all
	cd dist && sha256sum * > SHA256SUMS.txt
	@echo "Release artifacts in dist/:"
	@ls -lh dist/
