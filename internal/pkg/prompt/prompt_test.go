package prompt

import (
	"bytes"
	"strings"
	"testing"

	"github.com/burik666/gobin/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestContinue(t *testing.T) {
	const (
		msg      = "Do you want to continue? [Y/n] "
		okMsg    = msg + "\n"
		abortMsg = msg + "Abort.\n"
	)

	tests := []struct {
		noPrompt bool
		input    string
		result   bool
	}{
		{true, "", true},

		{false, "Y", true},
		{false, "y", true},
		{false, "Yes", true},
		{false, "YES", true},
		{false, "YeS", true},

		{false, "Hello", false},
		{false, "", false},
		{false, "N", false},
	}

	buf := bytes.NewBuffer(make([]byte, 1024))
	stdout = buf

	for _, tt := range tests {
		buf.Reset()

		t.Run("", func(t *testing.T) {
			config.NoPrompt = tt.noPrompt

			stdin = strings.NewReader(tt.input)

			res := Continue()

			assert.Equal(t, res, tt.result)

			if config.NoPrompt {
				if buf.Len() > 0 {
					t.Errorf("unexcepted output with NoPrompt")
				}

				assert.Empty(t, buf.String())

				return
			}

			if tt.result {
				assert.Equal(t, buf.String(), okMsg)
			} else {
				assert.Equal(t, buf.String(), abortMsg)
			}
		})
	}
}
