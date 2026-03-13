# Contributing to nix-config-collector

**This project is in active development and we welcome contributors of all experience levels!** 🎉

> ⚠️ **Heads up:** The codebase is evolving quickly. Things may change, break, or get refactored. Use at your own risk and with no warranty.

## Ways to contribute

- **Bug reports** — open an [issue](https://github.com/el-j/nix-config-collector/issues) with steps to reproduce
- **Feature requests** — open an [issue](https://github.com/el-j/nix-config-collector/issues) with your idea
- **Pull requests** — fork the repo, make your changes, and submit a PR
- **Documentation** — improve docs, fix typos, add examples
- **Testing** — try it on your macOS setup and report what breaks

## Development setup

```bash
git clone https://github.com/el-j/nix-config-collector
cd nix-config-collector
go mod download
go test ./...
go build ./cmd/cli/
```

## Running tests

```bash
go test -v -race ./...
```

## Code style

- Standard Go formatting (`gofmt` / `goimports`)
- Follow existing patterns in the codebase
- Add tests for new functionality

## Pull request checklist

- [ ] Tests pass (`go test ./...`)
- [ ] Code is formatted (`gofmt`)
- [ ] Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/)

## License

By contributing you agree that your contributions will be licensed under the [MIT License](LICENSE).
