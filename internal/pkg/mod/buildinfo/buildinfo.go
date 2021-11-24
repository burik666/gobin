// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildinfo

import (
	"bytes"
	"fmt"
	"strings"
)

// BuildInfo represents the build information read from a Go binary.
type BuildInfo struct {
	Filename  string
	GoVersion string         // Version of Go that produced this binary.
	Path      string         // The main package path
	Main      Module         // The module containing the main package
	Deps      []*Module      // Module dependencies
	Settings  []BuildSetting // Other information about the build.
}

// Module represents a module.
type Module struct {
	Path    string  // module path
	Version string  // module version
	Sum     string  // checksum
	Replace *Module // replaced by this module
}

// BuildSetting describes a setting that may be used to understand how the
// binary was built. For example, VCS commit and dirty status is stored here.
type BuildSetting struct {
	// Key and Value describe the build setting. They must not contain tabs
	// or newlines.
	Key, Value string
}

func (bi *BuildInfo) UnmarshalText(data []byte) (err error) {
	*bi = BuildInfo{}
	lineNum := 1

	defer func() {
		if err != nil {
			err = fmt.Errorf("could not parse Go build info: line %d: %w", lineNum, err)
		}
	}()

	var (
		pathLine  = []byte("path\t")
		modLine   = []byte("mod\t")
		depLine   = []byte("dep\t")
		repLine   = []byte("=>\t")
		buildLine = []byte("build\t")
		newline   = []byte("\n")
		tab       = []byte("\t")
		colon     = []byte(":")
	)

	readModuleLine := func(elem [][]byte) (Module, error) {
		if len(elem) != 2 && len(elem) != 3 {
			return Module{}, fmt.Errorf("expected 2 or 3 columns; got %d", len(elem))
		}

		sum := ""
		if len(elem) == 3 {
			sum = string(elem[2])
		}

		return Module{
			Path:    string(elem[0]),
			Version: strings.TrimPrefix(string(elem[1]), "go"),
			Sum:     sum,
		}, nil
	}

	var (
		last *Module
		line []byte
		ok   bool
	)

	// Reverse of BuildInfo.String(), except for go version.
	for len(data) > 0 {
		line, data, ok = bytesCut(data, newline)
		if !ok {
			break
		}

		line = bytes.TrimLeft(line, "\t")

		switch {
		case lineNum == 1:
			elem := bytes.SplitN(line, colon, 2)
			bi.Filename = string(elem[0])
			bi.GoVersion = string(elem[1][1:])

			if strings.HasPrefix(bi.GoVersion, "devel") {
				p := strings.SplitN(bi.GoVersion, " ", 3)
				bi.GoVersion = p[1]
			}

			bi.GoVersion = strings.TrimPrefix(bi.GoVersion, "go")
		case bytes.HasPrefix(line, pathLine):
			elem := line[len(pathLine):]
			bi.Path = string(elem)
		case bytes.HasPrefix(line, modLine):
			elem := bytes.Split(line[len(modLine):], tab)
			last = &bi.Main
			*last, err = readModuleLine(elem)

			if err != nil {
				return err
			}
		case bytes.HasPrefix(line, depLine):
			elem := bytes.Split(line[len(depLine):], tab)
			last = new(Module)
			bi.Deps = append(bi.Deps, last)
			*last, err = readModuleLine(elem)

			if err != nil {
				return err
			}
		case bytes.HasPrefix(line, repLine):
			elem := bytes.Split(line[len(repLine):], tab)
			if len(elem) != 3 {
				return fmt.Errorf("expected 3 columns for replacement; got %d", len(elem))
			}

			if last == nil {
				return fmt.Errorf("replacement with no module on previous line")
			}

			last.Replace = &Module{
				Path:    string(elem[0]),
				Version: string(elem[1]),
				Sum:     string(elem[2]),
			}

			last = nil
		case bytes.HasPrefix(line, buildLine):
			elem := bytes.Split(line[len(buildLine):], tab)
			value := ""

			if len(elem) >= 2 {
				value = string(elem[1])
			}

			if len(elem[0]) == 0 {
				return fmt.Errorf("empty key")
			}

			bi.Settings = append(bi.Settings, BuildSetting{Key: string(elem[0]), Value: value})
		case len(line) == 0:

		default:
			return nil
		}

		lineNum++
	}

	return nil
}

func bytesCut(s, sep []byte) (before, after []byte, found bool) {
	if i := bytes.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}

	return s, nil, false
}
