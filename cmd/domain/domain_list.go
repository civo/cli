package domain

import (
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var domainListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List domains",
	Long: `List all current domains.
If you wish to use a custom format, the available fields are:

	* id
	* name

Example: civo domain ls -o custom -f "ID: Name"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		domains, err := client.ListDNSDomains()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, domain := range domains {
			ow.StartLine()

			ow.AppendDataWithLabel("id", domain.ID, "ID")
			ow.AppendDataWithLabel("name", domain.Name, "Name")
		}

		ow.FinishAndPrintOutput()
	},
}
