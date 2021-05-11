package cmd

import (
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var regionListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo region ls`,
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
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			return
		}

		regions, err := client.ListRegions()
		if err != nil {
			utility.Error("%s", err)
			return
		}

		ow := utility.NewOutputWriter()
		// 		Type
		// OutOfCapacity
		// Country
		// CountryName
		// Features
		// Iaas
		// Kubernetes
		for _, region := range regions {
			ow.StartLine()
			ow.AppendData("Code", region.Code)
			ow.AppendData("Name", region.Name)
			ow.AppendData("Country", region.CountryName)

			defaultLabel := ""
			if outputFormat == "json" || outputFormat == "custom" {
				defaultLabel = utility.BoolToYesNo(region.Default)
			} else {
				if config.Current.Meta.DefaultRegion != "" {
					if region.Code == config.Current.Meta.DefaultRegion {
						defaultLabel = "<====="
					}
				} else {
					if region.Default {
						defaultLabel = "<====="
					}
				}

			}
			ow.AppendData("Current", defaultLabel)
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
