package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var apikeyRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "rm"},
	Short:   "Remove a saved API key",
	Args:    cobra.MinimumNArgs(1),
	Example: "civo apikey remove NAME",
	Run: func(cmd *cobra.Command, args []string) {
		index, err := apiKeyFind(args[0])
		if err != nil {
			utility.Error("Unable find the API key %s", err.Error())
			os.Exit(1)
		}

		if utility.UserConfirmedDeletion("api key", common.DefaultYes, args[0]) {
			numKeys := len(config.Current.APIKeys)
			delete(config.Current.APIKeys, index)
			config.SaveConfig()

			if numKeys > len(config.Current.APIKeys) {
				fmt.Printf("Removed the API Key %s\n", utility.Green(index))
			} else {
				utility.Error("The API Key %q couldn't be found", args[0])
				os.Exit(1)
			}
		} else {
			fmt.Println("Operation aborted.")
		}

	},
}
