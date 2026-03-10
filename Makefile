.PHONY: build build-wails test lint clean

# Default target: build CLI
build:
	go build -ldflags="-s -w" -o bin/nix-config-collector ./cmd/cli/

# Build Wails desktop app (requires Wails CLI and macOS)
build-wails:
	wails build -platform darwin/arm64

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
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/nix-config-collector-darwin-arm64 ./cmd/cli/
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/nix-config-collector-darwin-amd64 ./cmd/cli/
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/nix-config-collector-linux-amd64 ./cmd/cli/
