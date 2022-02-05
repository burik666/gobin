package gocmd

import (
	"fmt"
)

func Install(path, version string, args []string) error {
	pathVer := fmt.Sprintf("%s@%s", path, version)

	cmdArgs := []string{"install", "-v"}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, pathVer)

	_, err := goExec(cmdArgs...)

	return err
}
