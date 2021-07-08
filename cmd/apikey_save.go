package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var apikeySaveCmd = &cobra.Command{
	Use:     "save",
	Aliases: []string{"add", "store", "create", "new"},
	Short:   "Save a new API key",
	Args:    cobra.MinimumNArgs(0),
	Example: "civo apikey save NAME KEY",
	Run: func(cmd *cobra.Command, args []string) {

		var name, apiKey string
		var err error
		if len(args) == 0 {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Enter a nice name for this account/API Key: ")

			name, err = reader.ReadString('\n')
			if err != nil {
				utility.Error("Error reading name", err)
				os.Exit(1)
			}
			fmt.Printf("Enter the API key: ")
			apikeyBytes, err := terminal.ReadPassword(0)
			if err != nil {
				utility.Error("Error reading api key", err)
				os.Exit(1)
			}
			apiKey = string(apikeyBytes)
		} else if len(args) == 2 {
			name = args[0]
			apiKey = args[1]
		}

		config.Current.APIKeys[name] = apiKey
		if config.Current.Meta.DefaultRegion == "" {
			client, err := civogo.NewClientWithURL(apiKey, config.Current.Meta.URL, "")
			if err != nil {
				utility.Error("Unable to create a Civo API client, please report this at https://github.com/civo/cli")
				os.Exit(1)
			}
			region, err := client.GetDefaultRegion()
			if err != nil {
				utility.Error("Unable to get the default regions from the Civo API")
				os.Exit(1)
			}
			config.Current.Meta.DefaultRegion = region.Code
		}

		config.SaveConfig()

		if len(config.Current.APIKeys) == 1 {
			config.Current.Meta.CurrentAPIKey = name
			config.SaveConfig()
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"name": name, "key": apiKey})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("\nSaved the API Key %s ", utility.Green(name))
		}

	},
}
