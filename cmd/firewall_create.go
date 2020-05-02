package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
)

var firewallCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new firewall",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		firewall, err := client.NewFirewall(args[0])
		if err != nil {
			fmt.Printf("Unable to create firewall: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": firewall.ID, "Name": firewall.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created a firewall called %s with ID %s\n", aurora.Green(firewall.Name), aurora.Green(firewall.ID))
		}
	},
}
