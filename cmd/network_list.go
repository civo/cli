package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"strconv"
)

var networkListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List networks",
	Long: `List all available networks.
If you wish to use a custom format, the available fields are:

	* ID
	* Name
	* Label
	* Region
	* CIDR
	* Default

Example: civo network ls -o custom -f "ID: Name (CIDR)"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			return
		}

		networks, err := client.ListNetworks()
		if err != nil {
			fmt.Printf("Unable to list sizes: %s\n", aurora.Red(err))
			return
		}

		ow := utility.NewOutputWriter()

		for _, network := range networks {
			ow.StartLine()
			ow.AppendData("ID", network.ID)
			ow.AppendData("Name", network.Name)
			ow.AppendData("Label", network.Label)
			ow.AppendData("Region", network.Region)
			ow.AppendData("CIDR", network.CIDR)
			ow.AppendData("Default", strconv.FormatBool(network.Default))

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
