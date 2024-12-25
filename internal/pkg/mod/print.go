package mod

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/burik666/gobin/internal/config"
	"github.com/fatih/color"
)

var (
	colorName       = color.New(color.FgHiWhite).SprintFunc()
	colorKey        = color.New(color.FgGreen)
	colorNewVersion = color.New(color.FgHiGreen).SprintFunc()
	colorCurVersion = color.New(color.FgGreen).SprintFunc()
	colorTime       = color.New(color.FgHiBlack).SprintFunc()
	colorWarn       = color.New(color.FgRed).SprintFunc()
)

func (pkg *Pkg) Print() {
	name := pkg.Name
	if name == "" && pkg.BuildInfo != nil {
		name = path.Base(pkg.BuildInfo.Filename)
	}

	fmt.Printf("* %s\n", colorName(name))

	pkg.printFilename()
	pkg.printPath()
	pkg.printInstalledVersion()
	pkg.printSelectedVersion()
	pkg.printLatestVersion()
	pkg.printVersions()

	fmt.Print("\n")
}

func (pkg *Pkg) printFilename() {
	if pkg.BuildInfo != nil && pkg.BuildInfo.Filename != "" {
		printKeyValue("File", strings.TrimPrefix(pkg.BuildInfo.Filename, config.GOBIN+string(os.PathSeparator)))
	}
}

func (pkg *Pkg) printPath() {
	path := colorWarn("unknown")

	if pkg.BuildInfo != nil {
		if pkg.BuildInfo.Main.Path != "" {
			path = pkg.BuildInfo.Main.Path
		}
	} else if pkg.ModuleInfo != nil {
		path = pkg.ModuleInfo.Path
	}

	printKeyValue("Package", path)
}

func (pkg *Pkg) printSelectedVersion() {
	if pkg.SelectedVersion != nil {
		printKeyValue("Selected version",
			colorNewVersion(pkg.SelectedVersion.Version)+" "+
				colorTime(pkg.SelectedVersion.Time),
		)
		printKeyValue("  go.mod version", pkg.SelectedVersion.GoVersion)
	}
}

func (pkg *Pkg) printInstalledVersion() {
	if pkg.BuildInfo == nil {
		return
	}

	ver := colorWarn("unknown")

	if pkg.BuildInfo.Main.Version != "" {
		ver = colorCurVersion(pkg.BuildInfo.Main.Version)

		if pkg.ModuleInfo != nil {
			ver += " " + colorTime(pkg.ModuleInfo.Time)
		}
	}

	goVer := pkg.BuildInfo.GoVersion
	if pkg.hasGoUpdate() {
		goVer = colorWarn(goVer)
	}

	printKeyValue("Installed version", ver)
	printKeyValue("  GoVersion", goVer)
}

func (pkg *Pkg) printLatestVersion() {
	var mi *ModuleInfo

	if pkg.ModuleInfo != nil && pkg.ModuleInfo.Update != nil {
		mi = pkg.ModuleInfo.Update
	}

	if pkg.SelectedVersion != nil {
		mi = pkg.SelectedVersion.Update
	}

	if mi != nil && (pkg.SelectedVersion == nil || pkg.SelectedVersion.Version != mi.Version) {
		printKeyValue("Latest version",
			colorCurVersion(mi.Version)+" "+
				colorTime(mi.Time),
		)
		printKeyValue("  go.mod version", mi.GoVersion)
	}
}

func printKeyValue(key, value string) {
	fmt.Println(colorKey.Sprintf("  %-20s ", key+":"), value)
}

func (pkg *Pkg) printVersions() {
	bi := pkg.BuildInfo
	mi := pkg.ModuleInfo
	si := pkg.SelectedVersion

	if mi == nil {
		return
	}

	versions := mi.Versions

	output := make([]string, len(versions))
	for i := range versions {
		output[i] = versions[i]

		if si != nil && versions[i] == si.Version {
			output[i] = colorNewVersion(versions[i])

			continue
		}

		if mi != nil && versions[i] == mi.Version {
			output[i] = colorCurVersion(versions[i])

			continue
		}

		if bi != nil && versions[i] == bi.Main.Version {
			output[i] = colorCurVersion(versions[i])

			continue
		}
	}

	if len(versions) > 0 {
		printKeyValue("Available versions", strings.Join(output, " "))
	}

	if len(mi.Retracted) > 0 {
		printKeyValue("Retracted", colorWarn(strings.Join(mi.Retracted, " ")))
	}
}
