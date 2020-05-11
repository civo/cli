package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
)

var apikeyRemoveCmd = &cobra.Command{
	Use:     "remove NAME",
	Aliases: []string{"delete", "rm"},
	Short:   "Remove a saved API key",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, _, err := apiKeyFind(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		if utility.AskForConfirmDelete("api key") == nil {
			numKeys := len(config.Current.APIKeys)
			delete(config.Current.APIKeys, index)
			config.SaveConfig()

			if numKeys > len(config.Current.APIKeys) {
				fmt.Printf("Removed the API Key %s\n", utility.Green(index))
			} else {
				utility.Error("The API Key couldn't be found", args[0])
				os.Exit(1)
			}
		} else {
			fmt.Println("Operation aborted.")
		}

	},
}
