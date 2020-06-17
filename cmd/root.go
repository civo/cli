package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/spf13/cobra"
)

var outputFields string
var outputFormat string
var defaultYes bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "civo",
	Short: "CLI to manage cloud resources at Civo.com",
	Long: `civo is a CLI library for managing cloud resources such
as instances and Kubernetes clusters at Civo.com.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.ReadConfig)

	rootCmd.PersistentFlags().StringVarP(&config.Filename, "config", "", "", "config file (default is $HOME/.civo.json)")
	rootCmd.PersistentFlags().StringVarP(&outputFields, "fields", "f", "", "output fields for custom format output (use -h to determine fields)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "human", "output format (json/human/custom)")
	rootCmd.PersistentFlags().BoolVarP(&defaultYes, "yes", "y", false, "Automatic yes to prompts; assume \"yes\" as answer to all prompts and run non-interactively")

}
