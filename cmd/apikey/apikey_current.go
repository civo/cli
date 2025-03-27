package apikey

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var apikeyCurrentCmd = &cobra.Command{
	Use:     "set [NAME]",
	Aliases: []string{"use", "default", "current"},
	Short:   "Set the current API key",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, err := apiKeyFind(args[0])
		if err != nil {
			utility.Error("Unable to find the API key %s", err.Error())
			os.Exit(1)
		}

		if index != "" {
			config.Current.Meta.CurrentAPIKey = index
			config.SaveConfig()
			fmt.Printf("Set the default API Key to be %s\n", utility.Green(index))
		}

	},
}
