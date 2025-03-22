package apikey

import (
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var forceFlag bool

var apikeyRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "rm"},
	Short:   "Remove a saved API key",
	Args:    cobra.MinimumNArgs(1),
	Example: "civo apikey remove NAME",
	Run: func(cmd *cobra.Command, args []string) {
		index, err := apiKeyFind(args[0])
		if err != nil {
			utility.Error("Unable to find the API key %s", err.Error())
			os.Exit(1)
		}

		// Check if the requested API key is the current one
		if index == config.Current.Meta.CurrentAPIKey {
			if forceFlag {
				utility.Warning("The API key %q is the current one. You are using the --force flag, so it will be deleted.", args[0])
			} else {
				utility.Warning("The API key %q is the current one. If you remove it, you will need to set another API key as the current one to continue using the CLI.", args[0])
			}
		}

		// Confirm deletion of the API key
		if forceFlag || utility.UserConfirmedDeletion("API key", common.DefaultYes, args[0]) {
			numKeys := len(config.Current.APIKeys)
			delete(config.Current.APIKeys, index)
			config.SaveConfig()

			if numKeys > len(config.Current.APIKeys) {
				utility.Printf("Removed the API Key %s\n", utility.Green(index))
			} else {
				utility.Error("The API Key %q couldn't be found", args[0])
				os.Exit(1)
			}
		} else {
			utility.Println("Operation aborted.")
		}

	},
}
