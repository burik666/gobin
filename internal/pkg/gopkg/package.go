package gopkg

import (
	"strings"

	"github.com/burik666/gobin/internal/config"
	"github.com/burik666/gobin/internal/pkg/gocmd"
	"github.com/burik666/gobin/internal/pkg/mod"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

type Pkg struct {
	Name            string
	BuildInfo       *BuildInfo
	ModuleInfo      *mod.ModuleInfo
	SelectedVersion *mod.ModuleInfo
}

func (pkg *Pkg) FetchModuleInfo() error {
	mi, err := gocmd.ModuleInfo(pkg.BuildInfo.Main.Path, pkg.BuildInfo.Main.Version)
	if err != nil {
		return err
	}

	pkg.ModuleInfo = mi

	if module.IsPseudoVersion(pkg.BuildInfo.Main.Version) {
		up, err := gocmd.ModuleInfo(mi.Path, "latest")
		if err != nil {
			return err
		}

		if semver.Compare(up.Version, mi.Version) > 0 {
			mi.Update = up
		}
	}

	return nil
}

func (pkg *Pkg) HasUpdate() bool {
	if pkg.ModuleInfo == nil {
		return false
	}

	if (pkg.ModuleInfo.Update != nil) || pkg.hasGoUpdate() {
		return true
	}

	return false
}

func (pkg *Pkg) CheckPath() bool {
	return module.CheckPath(pkg.BuildInfo.Path) == nil
}

func (pkg *Pkg) hasGoUpdate() bool {
	pv := trimGoVersion(pkg.BuildInfo.GoVersion)
	cv := trimGoVersion(config.GoVersion)

	return semver.Compare(pv, cv) < 0
}

func trimGoVersion(v string) string {
	if strings.HasPrefix(v, "devel") {
		p := strings.SplitN(v, " ", 3)
		if len(p) >= 2 {
			v = p[1]
		}
	}

	return "v" + strings.TrimLeft(v, "go")
}

func ListInstalled(filter []string, exclude []string) ([]Pkg, error) {
	builds, err := gocmd.BuildInfo(config.GOBIN)
	if err != nil {
		return nil, err
	}

	res := make([]Pkg, 0, len(builds))

	for i := range builds {
		bi := BuildInfo{builds[i]}
		if len(filter) > 0 && !bi.Matched(filter) ||
			len(exclude) > 0 && bi.Matched(exclude) {
			continue
		}

		res = append(res, Pkg{
			Name:      bi.BuildInfo.Path,
			BuildInfo: &bi,
		})
	}

	return res, nil
}

func FindPackage(name, path, ver string) (*Pkg, error) {
	p := strings.Split(path, "/")

	var ferr error

	for {
		mi, err := gocmd.ModuleInfo(strings.Join(p, "/"), ver)
		if err != nil {
			if ferr == nil {
				ferr = err
			}

			if len(p) <= 2 {
				return nil, ferr
			}

			p = p[:len(p)-1]

			continue
		}

		return &Pkg{
			Name:       name,
			ModuleInfo: mi,
		}, err
	}
}
