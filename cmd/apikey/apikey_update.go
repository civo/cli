package apikey

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var apikeyUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"set", "edit"},
	Short:   "Update an existing API key",
	Args:    cobra.ExactArgs(2), // Expecting exactly two arguments: key name and new key value
	Example: "civo apikey update NAME NEW_API_KEY",
	Run: func(cmd *cobra.Command, args []string) {
		keyName := args[0]
		newKeyValue := args[1]

		if _, exists := config.Current.APIKeys[keyName]; !exists {
			utility.Error("No API key with name %q found. Use add command to insert new key.", keyName)
			os.Exit(1)
		}

		config.Current.APIKeys[keyName] = newKeyValue
		config.SaveConfig()
		fmt.Printf("Updated API Key %s\n", utility.Green(keyName))
	},
}
