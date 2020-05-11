package cmd

import (
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var regionListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List regions",
	Long: `List all available regions, including which is the default.
If you wish to use a custom format, the available fields are:

* Code
* Name
* Default

Example: civo region ls -o custom -f "Code: Name (Region)"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			return
		}

		regions, err := client.ListRegions()
		if err != nil {
			utility.Error("Unable to list regions %s", err)
			return
		}

		ow := utility.NewOutputWriter()

		for _, region := range regions {
			ow.StartLine()
			ow.AppendData("Code", region.Code)
			ow.AppendData("Name", region.Name)

			defaultLabel := ""
			if outputFormat == "json" || outputFormat == "custom" {
				defaultLabel = utility.BoolToYesNo(region.Default)
			} else {
				if region.Default {
					defaultLabel = "<====="
				}
			}
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
