package gocmd

import (
	"encoding/json"
	"fmt"

	"github.com/burik666/gobin/internal/pkg/mod"
)

var (
	moduleInfoCache  = make(map[string]mod.ModuleInfo)
	moduleInfoErrors = make(map[string]error)
)

func ModuleInfo(path, version string) (res *mod.ModuleInfo, rerr error) {
	mi, err := getModuleInfo(path, version)
	if err != nil {
		return nil, err
	}

	if mi.Version == "" {
		return nil, fmt.Errorf("no version")
	}

	if mi.Update != nil {
		mi.Update, err = getModuleInfo(path, version)
		if err != nil {
			return nil, err
		}
	}

	return mi, err
}

func getModuleInfo(path, version string) (res *mod.ModuleInfo, rerr error) {
	pathVer := fmt.Sprintf("%s@%s", path, version)
	if res, ok := moduleInfoCache[pathVer]; ok {
		return &res, nil
	}

	if err, ok := moduleInfoErrors[pathVer]; ok {
		return nil, err
	}

	defer func() {
		if rerr != nil {
			moduleInfoErrors[pathVer] = rerr
		} else {
			moduleInfoCache[pathVer] = *res
		}
	}()

	out, err := goExec("list", "-json", "-versions", "-m", "-u", pathVer)
	if err != nil {
		return nil, err
	}

	mi, err := parseModuleInfo(out)
	if err != nil {
		return nil, err
	}

	return mi, nil
}

func parseModuleInfo(data []byte) (*mod.ModuleInfo, error) {
	var mi mod.ModuleInfo

	if err := json.Unmarshal(data, &mi); err != nil {
		return nil, err
	}

	if mi.Version == "" {
		return nil, fmt.Errorf("empty version")
	}

	return &mi, nil
}
