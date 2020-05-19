package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
)

var apikeyCurrentCmd = &cobra.Command{
	Use:     "current [NAME]",
	Aliases: []string{"use", "default", "set"},
	Short:   "Set the API key",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := apiKeyFind(args[0])
		if err != nil {
			utility.Error("Unable find the api key %s", err.Error())
			os.Exit(1)
		}

		if index != "" {
			config.Current.Meta.CurrentAPIKey = index
			config.SaveConfig()
			fmt.Printf("Set the default API Key to be %s\n", utility.Green(index))
		}

	},
}
