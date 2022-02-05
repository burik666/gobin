package mod

import (
	"time"
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
