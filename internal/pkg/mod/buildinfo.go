package mod

import (
	"debug/buildinfo"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/burik666/gobin/internal/config"
	//"github.com/burik666/gobin/internal/pkg/mod/buildinfo"
)

type BuildInfo struct {
	Filename string
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

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	res := make([]BuildInfo, 0, len(files))

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		bi, err := buildinfo.ReadFile(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}

		res = append(res, BuildInfo{
			Filename:  filepath.Join(path, file.Name()),
			BuildInfo: *bi,
		})

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
