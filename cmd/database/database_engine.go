package database

import (
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var dbEngineCmd = &cobra.Command{
	Use:     "engine",
	Example: "civo db engine ls",
	Aliases: []string{"engines", "all", "software", "softwares"},
	Short:   "List Database engines",
	Long:    `List all currently available Database engines.`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		versions, err := client.ListDBVersions()
		if err != nil {
			utility.Error("%s", err)
			return
		}

		ow := utility.NewOutputWriter()
		for engine, version := range versions {
			ow.StartLine()
			ow.AppendDataWithLabel("engine", engine, "Engine")
			for _, v := range version {
				ow.AppendDataWithLabel("version", v.SoftwareVersion, "Version")
				if v.Default {
					ow.AppendDataWithLabel("default", v.SoftwareVersion, "Default Version")
				}
			}
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
