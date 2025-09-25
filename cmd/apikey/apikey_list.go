package apikey

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

* name
* key

Example: civo apikey ls -o custom -f "Name: Key"`,
	Run: func(cmd *cobra.Command, args []string) {
		keys := make([]string, 0, len(config.Current.APIKeys))
		for _, k := range config.Current.APIKeys {
			keys = append(keys, k.Name)
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
			ow.AppendDataWithLabel("name", name, "Name")
			// ow.AppendData("Key", apiKey)
			ow.AppendDataWithLabel("default", defaultLabel, "Default")

		}

		ow.FinishAndPrintOutput()
	},
}
