package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/burik666/gobin/internal/pkg/mod"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolVarP(&checkUpdates, "update", "u", false, "check for updates")
	listCmd.PersistentFlags().BoolVarP(&onlyUpdates, "upgrade", "U", false, "list of packages for which updates are available")
	listExclude = listCmd.PersistentFlags().StringSliceP("exclude", "e", nil, "comma-separated list of packages")

	listCmd.Flags().SortFlags = false
	listCmd.PersistentFlags().SortFlags = false
}

var listCmd = &cobra.Command{
	Use:   "list [package | filename]...",
	Short: "list installed packages",
	Example: "  gobin list\n" +
		"  gobin list goimports\n" +
		"  gobin list golang.org/x/tools/cmd/...\n" +
		"  gobin list goimports@v0.1.7",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		packages, err := mod.ListInstalled(args, *listExclude)
		if err != nil {
			return err
		}

		result := make([]mod.Pkg, 0)

		for _, pkg := range packages {
			if (checkUpdates || onlyUpdates) && pkg.CheckPath() {
				if err := pkg.FetchModuleInfo(); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}

			if onlyUpdates && (!pkg.HasUpdate() || !pkg.CheckPath()) {
				continue
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

		return nil
	},
}
