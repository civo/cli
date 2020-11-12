package cmd

import (
	"sort"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var apikeyShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"show"},
	Example: "civo apikey show NAME",
	Short:   "Show the default (current) API key",
	Long: `Show default API key.
If you wish to use a custom format, the available fields are:

* Name
* Key
`,
	Run: func(cmd *cobra.Command, args []string) {
		keys := make([]string, 0, len(config.Current.APIKeys))
		for k := range config.Current.APIKeys {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		ow := utility.NewOutputWriter()

		for _, name := range keys {
			apiKey := config.Current.APIKeys[name]
			if len(args) > 0 {
				if strings.Contains(name, args[0]) {
					ow.StartLine()
					ow.AppendData("Name", name)
					ow.AppendData("Key", apiKey)
				}
			} else {
				if config.Current.Meta.CurrentAPIKey == name {
					ow.StartLine()
					ow.AppendData("Name", name)
					ow.AppendData("Key", apiKey)
					// ow.AppendData("Default", defaultLabel)
				}
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
