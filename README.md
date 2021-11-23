# gobin
gobin is a package manager for /go/bin

## Features
- List installed packages.
- Check for updates.
- Install packages (like `go install`).
- Uninstall packages (just remove binaries).
- Update outdated packages.
- Reinstall packages which was built with old version of go.
- JSON format output.

## Installation

    go install github.com/burik666/gobin@latest

## Usage
Use `gobin --help` for more information.

    gobin install golang.org/x/tools/cmd/...
    gobin install golang.org/x/tools/cmd/...@v0.1.7 -- -ldflags=-s
    gobin list goimports
    gobin list golang.org/x/tools/cmd/...
    gobin upgrade golang.org/x/tools/cmd/...
    gobin uninstall golang.org/x/tools/cmd/...

## License

GPLv3

