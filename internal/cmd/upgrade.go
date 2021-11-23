package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/burik666/gobin/internal/pkg/mod"
	"github.com/burik666/gobin/internal/pkg/prompt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)

	upgradeExclude = upgradeCmd.PersistentFlags().StringSliceP("exclude", "e", nil, "comma-separated list of packages")
	upgradeCmd.PersistentFlags().BoolVarP(&dryRun, "dry", "d", false, "dry run")

	upgradeCmd.Flags().SortFlags = false
	upgradeCmd.PersistentFlags().SortFlags = false
}

var upgradeCmd = &cobra.Command{
	DisableFlagsInUseLine: true,

	Use:   "upgrade [flags] [package | filename]... -- [build flags]",
	Short: "Upgrade installed packages to the latest version.",
	Example: "  gobin upgrade\n" +
		"  gobin upgrade goimports\n" +
		"  gobin upgrade golang.org/x/tools/cmd/...",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		var installArgs []string

		if cmd.ArgsLenAtDash() > 0 {
			installArgs = args[cmd.ArgsLenAtDash():]
			args = args[:cmd.ArgsLenAtDash()]
		}

		result := make([]mod.Pkg, 0)

		packages, err := mod.ListInstalled(args, *upgradeExclude)
		if err != nil {
			return err
		}

		for _, pkg := range packages {
			if !pkg.CheckPath() {
				continue
			}

			if err := pkg.FetchModuleInfo(); err != nil {
				fmt.Fprintln(os.Stderr, err)

				continue
			}

			if !pkg.HasUpdate() {
				continue
			}

			if pkg.ModuleInfo.Update != nil {
				pkg.SelectedVersion = pkg.ModuleInfo.Update
			} else {
				pkg.SelectedVersion = pkg.ModuleInfo
			}

			result = append(result, pkg)
			if !jsonFormat {
				pkg.Print()
			}
		}

		if jsonFormat {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			_ = enc.Encode(result)
		}

		if len(result) == 0 {
			if !jsonFormat {
				fmt.Println("Nothing.")
			}

			return nil
		}

		if dryRun || !prompt.Continue() {
			return nil
		}

		for i := range result {
			name := result[i].Name
			ver := result[i].SelectedVersion.Version
			namever := fmt.Sprintf("%s@%s", name, ver)

			fmt.Printf("Installing [%d/%d] %s\n", i+1, len(result), color.YellowString(namever))

			if err := mod.Install(name, ver, installArgs); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}

		return nil
	},
}
