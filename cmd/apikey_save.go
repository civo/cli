package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var apikeySaveCmd = &cobra.Command{
	Use:     "save",
	Aliases: []string{"add", "store", "create", "new"},
	Short:   "Save a new API key",
	Args:    cobra.MinimumNArgs(2),
	Example: "civo apikey save NAME KEY",
	Run: func(cmd *cobra.Command, args []string) {
		config.Current.APIKeys[args[0]] = args[1]

		if config.Current.Meta.DefaultRegion == "" {
			client, err := config.CivoAPIClient()
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
			config.Current.Meta.CurrentAPIKey = args[0]
			config.SaveConfig()
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"Name": args[0], "Key": args[1]})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Saved the API Key %s as %s\n", utility.Green(args[0]), utility.Green(args[1]))
		}
	},
}
