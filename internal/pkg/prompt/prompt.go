package prompt

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/burik666/gobin/internal/config"
)

var (
	stdin  io.Reader = os.Stdin
	stdout io.Writer = os.Stdout
)

func Continue() bool {
	if config.NoPrompt {
		return true
	}

	fmt.Fprint(stdout, "Do you want to continue? [Y/n] ")

	var resp string

	n, err := fmt.Fscanln(stdin, &resp)
	resp = strings.ToUpper(resp)

	if err != nil || n == 0 || (resp != "Y" && resp != "YES") {
		fmt.Fprintln(stdout, "Abort.")

		return false
	}

	fmt.Fprintln(stdout)

	return true
}
