package mod

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/burik666/gobin/internal/config"
)

func GoVersion() (string, error) {
	cmd := exec.Command(config.Go, "version")

	if config.Trace {
		fmt.Fprintf(os.Stderr, "> %s\n", cmd.String())
	}

	var errb bytes.Buffer
	cmd.Stderr = &errb

	data, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s", errb.String())
	}

	v := string(data)
	v = strings.TrimPrefix(v, "go version ")

	p := strings.SplitN(v, " ", 3)
	if strings.HasPrefix(v, "devel") {
		v = p[1]
	} else {
		v = p[0]
	}

	v = strings.TrimPrefix(v, "go")

	return v, nil
}

func SplitVer(s string) (string, string) {
	p := strings.SplitN(s, "@", 2)

	name := p[0]
	ver := ""

	if len(p) > 1 {
		ver = p[1]
	}

	return name, ver
}
