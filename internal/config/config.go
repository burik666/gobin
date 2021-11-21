package config

import (
	"os"
	"path/filepath"
)

var (
	GOBIN     = binDir()
	NoPrompt  bool
	Trace     bool
	Go        = "go"
	GoVersion = ""
)

func binDir() string {
	if gobin := os.Getenv("GOBIN"); gobin != "" {
		return gobin
	}

	if gopath := os.Getenv("GOPATH"); gopath != "" {
		return filepath.Join(gopath, "bin")
	}

	if home, _ := os.UserHomeDir(); home != "" {
		return filepath.Join(home, "go", "bin")
	}

	return filepath.Join(os.Getenv("GOROOT"), "bin")
}
