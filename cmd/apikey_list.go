package cmd

import (
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var apikeyListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List stored API keys",
	Long: `List all API keys, making clear which is the current default.
If you wish to use a custom format, the available fields are:

* Name
* Key

Example: civo apikey ls -o custom -f "Name: Key"`,
	Run: func(cmd *cobra.Command, args []string) {
		ow := utility.NewOutputWriter()

		for _, name := range config.Current {
			ow.StartLine()
			for k, _ := range name.APIKeys {
				apiKey := name.APIKeys[k]
				defaultLabel := ""
				if name.Meta.CurrentAPIKey == k {
					defaultLabel = "<====="
				}
				ow.AppendData("Name", k)
				ow.AppendData("Key", apiKey)
				ow.AppendData("Default", defaultLabel)
			}
		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
