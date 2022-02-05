package gocmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		output  string
		version string
		err     string
	}{
		{"go version go1.17.2 linux/amd64", "go1.17.2", ""},
		{"go version go1.17.2", "go1.17.2", ""},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			version, err := parseGoVersion(tt.output)
			if tt.err != "" {
				if !assert.EqualError(t, err, tt.err) {
					return
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.version, version)
		})
	}
}
