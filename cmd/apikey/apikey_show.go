package apikey

import (
	"sort"
	"strings"

	"github.com/civo/cli/common"
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

* name
* key
`,
	Run: func(cmd *cobra.Command, args []string) {
		keys := make([]string, 0, len(config.Current.APIKeys))
		for _, apiKey := range config.Current.APIKeys {
			keys = append(keys, apiKey.Name)
		}
		sort.Strings(keys)

		ow := utility.NewOutputWriter()
		for _, name := range keys {
			var idx int
			for index, apiKey := range config.Current.APIKeys {
				if name == apiKey.Name {
					idx = index
					break
				}
			}
			apiKey := config.Current.APIKeys[idx]
			if len(args) > 0 {
				if strings.Contains(name, args[0]) {
					ow.StartLine()
					ow.AppendDataWithLabel("name", apiKey.Name, "Name")
					ow.AppendDataWithLabel("key", apiKey.Value, "Key")
					ow.AppendDataWithLabel("url", apiKey.APIURL, "APIURL")
				}
			} else {
				if config.Current.Meta.CurrentAPIKey == name {
					ow.StartLine()
					ow.AppendDataWithLabel("name", apiKey.Name, "Name")
					ow.AppendDataWithLabel("key", apiKey.Value, "Key")
					ow.AppendDataWithLabel("url", apiKey.APIURL, "APIURL")
				}
			}
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
