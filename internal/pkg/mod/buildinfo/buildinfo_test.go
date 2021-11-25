package buildinfo

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const goodModuleLine = "example.org/foo/bar\tv0.1.0\tsum"

var goodModule = Module{
	Path:    "example.org/foo/bar",
	Version: "v0.1.0",
	Sum:     "sum",
}

func TestReadModuleLine(t *testing.T) {
	tests := []struct {
		line   string
		module *Module
		err    string
	}{
		{
			"",
			nil,
			"expected 2 or 3 columns; got %d",
		},
		{
			"\t\t\t",
			nil,
			"expected 2 or 3 columns; got %d",
		},
		{
			"\t\t\t",
			nil,
			"expected 2 or 3 columns; got %d",
		},
		{
			"example.org/foo/bar\tversion",
			&Module{
				Path:    "example.org/foo/bar",
				Version: "version",
			},
			"",
		},
		{
			line:   goodModuleLine,
			module: &goodModule,
			err:    "",
		},
	}

	tab := []byte("\t")

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			elems := bytes.Split([]byte(tt.line), tab)

			m, err := readModuleLine(elems)
			if err != nil && !assert.EqualErrorf(t, err, fmt.Sprintf(tt.err, len(elems)), "") {
				return
			}

			if tt.err == "" && !assert.Empty(t, err) {
				return
			}

			if tt.module != nil {
				assert.Equal(t, *tt.module, m)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	errMsg := func(lineNum int, msg string, args ...interface{}) string {
		return fmt.Sprintf("could not parse Go build info: line %d: %s", lineNum, fmt.Sprintf(msg, args...))
	}

	tests := []struct {
		data string
		bi   BuildInfo
		err  string
	}{
		{"", BuildInfo{}, ""},
		{"foo", BuildInfo{}, ""},
		{"foo\n", BuildInfo{}, errMsg(1, `excepted 2 columns separated by colon; got: "foo"`)},

		// check first line
		{
			"foo: version\n",
			BuildInfo{
				Filename:  "foo",
				GoVersion: "version",
			},
			"",
		},
		{
			"foo: go1.17.2\n",
			BuildInfo{
				Filename:  "foo",
				GoVersion: "go1.17.2",
			},
			"",
		},
		{
			"foo: devel\n",
			BuildInfo{
				Filename:  "foo",
				GoVersion: "devel",
			},
			"",
		},
		{
			"foo: devel go1.2.3 a b c\n",
			BuildInfo{
				Filename:  "foo",
				GoVersion: "devel go1.2.3 a b c",
			},
			"",
		},

		{
			"foo: go1.2.3\n" +
				"path\tfoo/bar\n",
			BuildInfo{
				Filename:  "foo",
				GoVersion: "go1.2.3",
				Path:      "foo/bar",
			},
			"",
		},

		{
			"foo: go1.2.3\n" +
				fmt.Sprintf("mod\t%s\n", goodModuleLine),
			BuildInfo{
				Filename:  "foo",
				GoVersion: "go1.2.3",
				Main:      goodModule,
			},
			"",
		},
		{
			"foo: go1.2.3\n" +
				fmt.Sprintf("dep\t%s\n", goodModuleLine) +
				fmt.Sprintf("dep\t%s\n", goodModuleLine),
			BuildInfo{
				Filename:  "foo",
				GoVersion: "go1.2.3",
				Deps: []*Module{
					&goodModule,
					&goodModule,
				},
			},
			"",
		},
		{
			"foo: go1.2.3\n" +
				"=>\t./bar\n",
			BuildInfo{
				Filename:  "foo",
				GoVersion: "go1.2.3",
			},
			errMsg(2, "expected 3 columns for replacement; got 1"),
		},
		{
			"foo: go1.2.3\n" +
				"=>\t./bar\tver\tsum\n",
			BuildInfo{
				Filename:  "foo",
				GoVersion: "go1.2.3",
			},
			errMsg(2, "replacement with no module on previous line"),
		},
		{
			"foo: go1.2.3\n" +
				"dep\texample.org/foo/bar\tver\tsum\n" +
				fmt.Sprintf("=>\t%s\n", goodModuleLine),
			BuildInfo{
				Filename:  "foo",
				GoVersion: "go1.2.3",
				Deps: []*Module{
					{
						Path:    "example.org/foo/bar",
						Version: "ver",
						Sum:     "sum",
						Replace: &goodModule,
					},
				},
			},
			"",
		},

		{
			"foo: go1.2.3\n" +
				fmt.Sprintf("mod\t%s\n", goodModuleLine) +
				"build\t\n",
			BuildInfo{
				Filename:  "foo",
				GoVersion: "go1.2.3",
				Main:      goodModule,
			},
			errMsg(3, "empty key"),
		},
		{
			"foo: go1.2.3\n" +
				fmt.Sprintf("mod\t%s\n", goodModuleLine) +
				"build\tkey\n" +
				"build\tkey\tvalue\n",
			BuildInfo{
				Filename:  "foo",
				GoVersion: "go1.2.3",
				Main:      goodModule,
				Settings: []BuildSetting{
					{Key: "key"},
					{Key: "key", Value: "value"},
				},
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			var bi BuildInfo

			err := bi.UnmarshalText([]byte(tt.data))
			if tt.err != "" {
				if !assert.EqualError(t, err, tt.err) {
					return
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.bi, bi)
		})
	}
}
