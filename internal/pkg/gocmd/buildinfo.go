package gocmd

import (
	"bytes"
	"fmt"

	"github.com/burik666/gobin/internal/pkg/mod"
)

var (
	buildInfoCache  = make(map[string][]mod.BuildInfo)
	buildInfoErrors = make(map[string]error)
)

func BuildInfo(path string) (res []mod.BuildInfo, rerr error) {
	if res, ok := buildInfoCache[path]; ok {
		return res, nil
	}

	if err, ok := buildInfoErrors[path]; ok {
		return nil, err
	}

	defer func() {
		if rerr != nil {
			buildInfoErrors[path] = rerr
		} else {
			buildInfoCache[path] = res
		}
	}()

	out, err := goExec("version", "-m", path)
	if err != nil {
		return nil, err
	}

	return parseBuildInfo(out)
}

func parseBuildInfo(data []byte) ([]mod.BuildInfo, error) {
	res := make([]mod.BuildInfo, 0)

	for len(data) > 0 {
		var bi mod.BuildInfo

		if err := bi.UnmarshalText(data); err != nil {
			return nil, err
		}

		res = append(res, bi)

		for len(data) > 0 {
			n := bytes.IndexByte(data, '\n')
			if n < 0 {
				return nil, fmt.Errorf("no newline at end of data")
			}

			data = data[n+1:]
			if len(data) > 0 && data[0] != '\t' {
				break
			}
		}
	}

	return res, nil
}
