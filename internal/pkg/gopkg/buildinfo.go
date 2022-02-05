package gopkg

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/burik666/gobin/internal/pkg/mod"
)

type BuildInfo struct {
	mod.BuildInfo
}

func (bi *BuildInfo) Matched(filters []string) bool {
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

	return false
}

func SplitVer(s string) (string, string) {
	p := strings.SplitN(s, "@", 2)

	name := p[0]
	ver := ""

	if len(p) > 1 {
		ver = p[1]
	}

	return name, ver
}
