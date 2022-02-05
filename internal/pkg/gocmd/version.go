package gocmd

import (
	"strings"
)

func GoVersion() (string, error) {
	out, err := goExec("version")
	if err != nil {
		return "", err
	}

	return parseGoVersion(string(out))
}

func parseGoVersion(v string) (string, error) {
	v = strings.TrimPrefix(v, "go version ")

	i := strings.LastIndex(v, " ")
	if i > 0 {
		v = v[:i]
	}

	return v, nil
}
