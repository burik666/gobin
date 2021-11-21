package prompt

import (
	"fmt"
	"strings"

	"github.com/burik666/gobin/internal/config"
)

func Continue() bool {
	if config.NoPrompt {
		return true
	}

	fmt.Printf("Do you want to continue? [Y/n] ")

	var resp string

	n, err := fmt.Scanln(&resp)
	resp = strings.ToUpper(resp)

	if err != nil || n == 0 || (resp != "Y" && resp != "YES") {
		fmt.Println("Abort.")

		return false
	}

	fmt.Println()

	return true
}
