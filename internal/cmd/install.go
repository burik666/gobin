package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/burik666/gobin/internal/pkg/mod"
	"github.com/burik666/gobin/internal/pkg/prompt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.PersistentFlags().BoolVarP(&dryRun, "dry", "d", false, "dry run")

	installCmd.Flags().SortFlags = false
	installCmd.PersistentFlags().SortFlags = false
}

var installCmd = &cobra.Command{
	Use:   "install [packages]...",
	Short: "Install packages",
	Example: "  gobin install golang.org/x/tools/cmd/goimports\n" +
		"  gobin install golang.org/x/tools/cmd/...\n" +
		"  gobin install golang.org/x/tools/cmd/goimports@v0.1.7\n" +
		"\n  If the package is already installed you can use the filename.\n" +
		"  gobin install goimports@v0.1.7",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		var result struct {
			Installed []mod.Pkg
			New       []mod.Pkg
		}

		for _, m := range args {
			name, ver := mod.SplitVer(m)
			if ver == "" {
				ver = "latest"
			}

			newPkg, ferr := mod.FindPackage(name, name, ver)
			if ferr != nil {
				s, err := mod.ListInstalled([]string{name}, nil)
				if err != nil || len(s) == 0 || !s[0].CheckPath() {
					return ferr
				}

				if len(s) == 0 {
					return err
				}

				ok := false
				for i := range s {
					if s[i].CheckPath() {
						newPkg, err := mod.FindPackage(s[i].BuildInfo.Path, s[i].BuildInfo.Main.Path, ver)
						if err != nil {
							return err
						}

						if err := s[i].FetchModuleInfo(); err != nil {
							fmt.Fprintln(os.Stderr, err)

							continue
						}

						s[i].SelectedVersion = newPkg.ModuleInfo
						result.New = append(result.New, s[i])

						ok = true
					}
				}

				if !ok {
					return ferr
				}
			} else {
				installedPackages, err := mod.ListInstalled([]string{name}, nil)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				if len(installedPackages) == 1 && !strings.HasSuffix(name, "...") {
					pkg := installedPackages[0]

					if err := pkg.FetchModuleInfo(); err != nil {
						fmt.Fprintln(os.Stderr, err)
					}

					pkg.SelectedVersion = newPkg.ModuleInfo

					result.New = append(result.New, pkg)
				} else {
					newPkg.SelectedVersion = newPkg.ModuleInfo
					result.New = append(result.New, *newPkg)

					for _, pkg := range installedPackages {
						if err := pkg.FetchModuleInfo(); err != nil {
							fmt.Fprintln(os.Stderr, err)
						}

						pkg.SelectedVersion = newPkg.ModuleInfo

						result.Installed = append(result.Installed, pkg)
					}
				}
			}
		}

		if jsonFormat {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			_ = enc.Encode(result)
		} else {
			if len(result.Installed) != 0 {
				fmt.Println("Installed packages:")
				for _, pkg := range result.Installed {
					pkg.Print()
				}
			}

			if len(result.New) != 0 {
				fmt.Println("Install packages:")
				for _, pkg := range result.New {
					pkg.Print()
				}
			}
		}

		if len(result.New) == 0 && len(result.Installed) == 0 {
			if !jsonFormat {
				fmt.Println("Nothing.")
			}

			return nil
		}

		if dryRun || !prompt.Continue() {
			return nil
		}

		i := 0
		for _, pkg := range result.New {
			name := pkg.Name
			ver := pkg.SelectedVersion.Version
			namever := fmt.Sprintf("%s@%s", name, ver)

			fmt.Printf("Installing [%d/%d] %s\n", i+1, len(result.New), color.YellowString(namever))

			if err := mod.Install(name, ver); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			i++
		}

		return nil
	},
}
