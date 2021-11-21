package mod

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/burik666/gobin/internal/config"
	"github.com/burik666/gobin/internal/pkg/mod/buildinfo"
)

type BuildInfo struct {
	buildinfo.BuildInfo
}

func (bi *BuildInfo) Matched(filters []string) bool {
	if len(filters) > 0 {
		for i := range filters {
			name, ver := SplitVer(filters[i])
			if ((bi.Path == "" && name == path.Base(bi.Filename)) ||
				name == filepath.Base(bi.Path) ||
				name == bi.Path ||
				name == bi.Main.Path) &&
				(ver == "" || ver == "latest" || ver == bi.Main.Version) {
				return true
			}

			if strings.HasSuffix(name, "...") && strings.HasPrefix(bi.Path, strings.TrimSuffix(name, "...")) {
				return true
			}
		}
	}

	return false
}

var buildInfoCache []BuildInfo

func getBuildInfo(path string) ([]BuildInfo, error) {
	if buildInfoCache != nil {
		return buildInfoCache, nil
	}

	cmd := exec.Command(config.Go, "version", "-m", path)

	if config.Trace {
		fmt.Fprintf(os.Stderr, "> %s\n", cmd.String())
	}

	var errb bytes.Buffer
	cmd.Stderr = &errb

	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%s", errb.String())
	}

	res := make([]BuildInfo, 0)

	for len(data) > 0 {
		var bi BuildInfo

		if err := bi.UnmarshalText(data); err != nil {
			return nil, err
		}

		res = append(res, bi)

		for len(data) > 0 {
			n := bytes.IndexByte(data, '\n')
			if n < 0 {
				break
			}

			data = data[n+1:]
			if len(data) > 0 && data[0] != '\t' {
				break
			}
		}
	}

	buildInfoCache = res

	return res, nil
}

func Uninstall(bi *BuildInfo) error {
	if config.Trace {
		fmt.Fprintf(os.Stderr, "rm %s\n", bi.Filename)
	}

	return os.Remove(bi.Filename)
}
