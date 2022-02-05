package gocmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/burik666/gobin/internal/pkg/mod"

	"github.com/stretchr/testify/assert"
)

const goodModuleLine = "example.org/foo/bar\tv0.1.0\tsum"

var goodModule = mod.Module{
	Path:    "example.org/foo/bar",
	Version: "v0.1.0",
	Sum:     "sum",
}

var goodBuildInfoData = "/path/to/file: go1.2.3\n" +
	"\tpath\texample.org/foo/bar\n" +
	fmt.Sprintf("\tmod\t%s\n", goodModuleLine)

var goodBuildInfo = mod.BuildInfo{
	Filename:  "/path/to/file",
	GoVersion: "go1.2.3",
	Path:      "example.org/foo/bar",
	Main:      goodModule,
}

func TestParseBuildInfo(t *testing.T) {
	tests := []struct {
		output    string
		buildinfo []mod.BuildInfo
		err       string
	}{
		{
			"",
			nil,
			"",
		},
		{
			strings.TrimRight(goodBuildInfoData, "\n"),
			nil,
			"no newline at end of data",
		},
		{
			goodBuildInfoData,
			[]mod.BuildInfo{
				goodBuildInfo,
			},
			"",
		},
		{
			goodBuildInfoData +
				goodBuildInfoData,
			[]mod.BuildInfo{
				goodBuildInfo,
				goodBuildInfo,
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, err := parseBuildInfo([]byte(tt.output))
			if tt.err != "" {
				if !assert.EqualError(t, err, tt.err) {
					return
				}
			} else {
				assert.NoError(t, err)
			}

			if len(res) == 0 && len(tt.buildinfo) == 0 {
				return
			}

			assert.Equal(t, tt.buildinfo, res)
		})
	}
}
