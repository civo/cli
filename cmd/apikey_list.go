package cmd

import (
	"sort"

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
		keys := make([]string, 0, len(config.Current.APIKeys))
		for k := range config.Current.APIKeys {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		ow := utility.NewOutputWriter()

		for _, name := range keys {
			ow.StartLine()
			// apiKey := config.Current.APIKeys[name]
			defaultLabel := ""
			if config.Current.Meta.CurrentAPIKey == name {
				defaultLabel = "<====="
			}
			ow.AppendData("Name", name)
			// ow.AppendData("Key", apiKey)
			ow.AppendData("Default", defaultLabel)

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
