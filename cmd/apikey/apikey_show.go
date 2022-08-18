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
					ow.AppendDataWithLabel("name", name, "Name")
					ow.AppendDataWithLabel("key", apiKey, "Key")
				}
			} else {
				if config.Current.Meta.CurrentAPIKey == name {
					ow.StartLine()
					ow.AppendDataWithLabel("name", name, "Name")
					ow.AppendDataWithLabel("key", apiKey, "Key")
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
