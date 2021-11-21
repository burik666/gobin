package cmd

import (
	"fmt"
	"os"

	"github.com/burik666/gobin/internal/config"
	"github.com/burik666/gobin/internal/pkg/mod"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "gobin",
	Short:         "GOBIN modules manager",
	SilenceErrors: true,
	Example: "  gobin install golang.org/x/tools/...\n" +
		"  gobin install golang.org/x/tools/...@v0.1.7\n" +
		"  gobin list goimports\n" +
		"  gobin list golang.org/x/tools/cmd/...\n" +
		"  gobin upgrade golang.org/x/tools/cmd/...\n" +
		"  gobin uninstall golang.org/x/tools/cmd/...\n",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		goVersion, err := mod.GoVersion()
		if err != nil {
			fmt.Fprint(os.Stderr, err)

			os.Exit(1)
		}

		config.GoVersion = goVersion

		return nil
	},
}

func Execute() {
	cobra.EnableCommandSorting = false

	//	rootCmd.PersistentFlags().StringVarP(&config.GOBIN, "gobin", "b", config.GOBIN, "path to GOBIN directory")
	rootCmd.PersistentFlags().BoolVarP(&jsonFormat, "json", "j", false, "print result in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&config.NoPrompt, "yes", "y", false, "always yes to prompts")
	rootCmd.PersistentFlags().BoolVarP(&color.NoColor, "nocolor", "n", false, "disables color output")
	rootCmd.PersistentFlags().BoolVarP(&config.Trace, "trace", "x", false, "trace commands")
	rootCmd.PersistentFlags().StringVar(&config.Go, "go", "go", "go command")

	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
