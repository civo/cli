package region

import (
	"github.com/civo/cli/common"
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

	* code
	* name
	* country
	* current

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

		for _, region := range regions {
			ow.StartLine()
			ow.AppendDataWithLabel("code", region.Code, "Code")
			ow.AppendDataWithLabel("name", region.Name, "Name")
			ow.AppendDataWithLabel("country", region.CountryName, "Country")

			defaultLabel := ""
			if common.OutputFormat == "json" || common.OutputFormat == "custom" {
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
			ow.AppendDataWithLabel("current", defaultLabel, "Current")
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
