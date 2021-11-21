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
	rootCmd.AddCommand(uninstallCmd)

	uninstallExclude = uninstallCmd.PersistentFlags().StringSliceP("exclude", "e", nil, "comma-separated list of packages")
	uninstallCmd.PersistentFlags().BoolVarP(&dryRun, "dry", "d", false, "dry run")

	uninstallCmd.Flags().SortFlags = false
	uninstallCmd.PersistentFlags().SortFlags = false
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [package | filename]...",
	Short: "Delete installed packages",
	Example: "  gobin uninstall goimports\n" +
		"  gobin uninstall golang.org/x/tools/cmd/goimports\n" +
		"  gobin uninstall golang.org/x/tools/cmd/...\n",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		var result []*mod.Pkg

		packages, err := mod.ListInstalled(args, *uninstallExclude)
		if err != nil {
			return err
		}

		for _, pkg := range packages {
			result = append(result, &mod.Pkg{BuildInfo: pkg.BuildInfo})

			if !jsonFormat {
				pkg.Print()
			}
		}

		if jsonFormat {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			_ = enc.Encode(result)
		}

		if len(result) == 0 && !jsonFormat {
			fmt.Println("Nothing.")

			return nil
		}

		if dryRun || !prompt.Continue() {
			return nil
		}

		for i := range result {
			namever := ""
			if result[i].BuildInfo.Path == "" {
				namever = result[i].BuildInfo.Filename
			} else {
				namever = fmt.Sprintf("%s@%s", result[i].BuildInfo.Path, result[i].BuildInfo.Main.Version)
			}

			fmt.Printf("Uninstalling [%d/%d] %s\n", i+1, len(result), color.YellowString(namever))

			if err := mod.Uninstall(result[i].BuildInfo); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}

		return nil
	},
}
