package gocmd

import (
	"fmt"
	"os"

	"github.com/burik666/gobin/internal/config"
	"github.com/burik666/gobin/internal/pkg/mod"
)

func Uninstall(bi *mod.BuildInfo) error {
	if config.Trace {
		fmt.Fprintf(os.Stderr, "rm %s\n", bi.Filename)
	}

	return os.Remove(bi.Filename)
}
