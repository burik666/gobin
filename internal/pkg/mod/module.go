package mod

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/burik666/gobin/internal/config"
)

type ModuleInfo struct {
	Path       string       `json:",omitempty"` // module path
	Version    string       `json:",omitempty"` // module version
	Versions   []string     `json:",omitempty"` // available module versions
	Replace    *ModuleInfo  `json:",omitempty"` // replaced by this module
	Time       *time.Time   `json:",omitempty"` // time version was created
	Update     *ModuleInfo  `json:",omitempty"` // available update (with -u)
	Main       bool         `json:",omitempty"` // is this the main module?
	Indirect   bool         `json:",omitempty"` // module is only indirectly needed by main module
	Dir        string       `json:",omitempty"` // directory holding local copy of files, if any
	GoMod      string       `json:",omitempty"` // path to go.mod file describing module, if any
	GoVersion  string       `json:",omitempty"` // go version used in module
	Retracted  []string     `json:",omitempty"` // retraction information, if any (with -retracted or -u)
	Deprecated string       `json:",omitempty"` // deprecation message, if any (with -u)
	Error      *ModuleError `json:",omitempty"` // error loading module
}

type ModuleError struct {
	Err string // the error itself
}

var (
	modulesCache = make(map[string]*ModuleInfo)
	moduleErrors = make(map[string]error)
)

func getModuleInfo(path, version string) (rmi *ModuleInfo, rerr error) {
	name := fmt.Sprintf("%s@%s", path, version)

	if mi, ok := modulesCache[name]; ok {
		return mi, nil
	}

	if err, ok := moduleErrors[name]; ok {
		return nil, err
	}

	defer func() {
		if rerr != nil {
			moduleErrors[name] = rerr
		}
	}()

	cmd := exec.Command(config.Go, "list", "-json", "-versions", "-m", "-u", name)

	if config.Trace {
		fmt.Fprintf(os.Stderr, "> %s\n", cmd.String())
	}

	var errb bytes.Buffer
	cmd.Stderr = &errb

	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%s", errb.String())
	}

	fmt.Fprint(os.Stderr, errb.String())

	var mi ModuleInfo

	if err := json.Unmarshal(data, &mi); err != nil {
		return nil, err
	}

	if mi.Version == "" {
		return nil, fmt.Errorf("empty version")
	}

	modulesCache[name] = &mi

	if mi.Update != nil {
		mi.Update, err = getModuleInfo(mi.Path, mi.Update.Version)
		if err != nil {
			return nil, err
		}
	}

	return &mi, nil
}

func Install(path, version string) error {
	name := fmt.Sprintf("%s@%s", path, version)
	cmd := exec.Command(config.Go, "install", "-v", name)

	if config.Trace {
		fmt.Fprintf(os.Stderr, "> %s\n", cmd.String())
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
