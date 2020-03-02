package cmd

import (
	"fmt"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// regionCmd represents the region command
var regionCmd = &cobra.Command{
	Use:     "region",
	Aliases: []string{"regions"},
	Short:   "Details of Civo regions",
}

// regionListCmd represents the command to list available API keys
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
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			return
		}

		regions, err := client.ListRegions()
		if err != nil {
			fmt.Printf("Unable to list regions: %s\n", aurora.Red(err))
			return
		}

		ow := utility.NewOutputWriter()

		for _, region := range regions {
			ow.StartLine()
			ow.AppendData("Code", region.Code)
			ow.AppendData("Name", region.Name)

			defaultLabel := ""
			if OutputFormat == "json" || OutputFormat == "custom" {
				defaultLabel = utility.BoolToYesNo(region.Default)
			} else {
				if region.Default {
					defaultLabel = "<====="
				}
			}
			ow.AppendData("Default", defaultLabel)
		}

		switch OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(OutputFields)
		default:
			ow.WriteTable()
		}
	},
}

func init() {
	rootCmd.AddCommand(regionCmd)

	regionCmd.AddCommand(regionListCmd)
}
