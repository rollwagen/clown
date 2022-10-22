package cmd

import (
	"os"

	"github.com/rollwagen/clown/cmd/gitlab"
	"github.com/spf13/cobra"
)

var Verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clown",
	Short: "Recursively clone git repositories.",
	Long:  `Recursively clone git repositories.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable debug output")

	addAllCommands()
}

func addAllCommands() {
	rootCmd.AddCommand(gitlab.NewCmd())
}
