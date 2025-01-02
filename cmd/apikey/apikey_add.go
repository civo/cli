package apikey

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var apikeyAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create", "new"},
	Short:   "Add a new API key",
	Args:    cobra.ExactArgs(2), // Expecting exactly two arguments: key name and key value
	Example: "civo apikey add NAME API_KEY",
	Run: func(cmd *cobra.Command, args []string) {
		keyName := args[0]
		keyValue := args[1]

		if _, exists := config.Current.APIKeys[keyName]; exists {
			utility.Error("API key with name %q already exists. Choose a different name or update the existing key.", keyName)
			os.Exit(1)
		}

		config.Current.APIKeys[keyName] = keyValue
		config.SaveConfig()
		fmt.Printf("Added new API Key %s\n", utility.Green(keyName))
	},
}
