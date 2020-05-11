package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var apikeyCurrentCmd = &cobra.Command{
	Use:     "current [NAME]",
	Aliases: []string{"use", "default", "set"},
	Short:   "Show the current API key",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, name, err := apiKeyFind(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		if name != "" {
			config.Current.Meta.CurrentAPIKey = index
			config.SaveConfig()

			fmt.Printf("Set the default API Key to be %s\n", aurora.Green(name))
		}
	},
}
