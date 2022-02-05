package gocmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/burik666/gobin/internal/config"
)

func goExec(args ...string) ([]byte, error) {
	cmd := exec.Command(config.Go, args...)

	if config.Trace {
		fmt.Fprintf(os.Stderr, "> %s\n", cmd.String())
	}

	var stderr bytes.Buffer
	cmd.Stderr = io.MultiWriter(&stderr, os.Stderr)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return stdout.Bytes(), fmt.Errorf("%w: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}
