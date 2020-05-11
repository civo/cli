package cmd

import (
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
)

var domainListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List domains",
	Long: `List all current domains.
If you wish to use a custom format, the available fields are:

	* ID
	* Name

Example: civo domain ls -o custom -f "ID: Name"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		domains, err := client.ListDNSDomains()
		if err != nil {
			utility.Error("Unable to list domains %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, domain := range domains {
			ow.StartLine()

			ow.AppendData("ID", domain.ID)
			ow.AppendData("Name", domain.Name)
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
