package database

import (
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
)

var dbVersionListCmd = &cobra.Command{
	Use:     "versions",
	Aliases: []string{"version"},
	Example: `civo db versions`,
	Short:   "List all the available database versions",
	Run: func(cmd *cobra.Command, args []string) {

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		dbVersions, err := client.ListDBVersions()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		for dbName, versionDetails := range dbVersions {
			ow.StartLine()
			ow.AppendDataWithLabel("name", dbName, "Name")
			ow.AppendDataWithLabel("version", versionDetails[0].SoftwareVersion, "version")
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
