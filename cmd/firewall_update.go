package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
)

var firewallUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"rename", "change"},
	Short:   "Update a firewall",
	Example: "civo firewall update OLD_NAME NEW_NAME",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		firewall, err := client.FindFirewall(args[0])
		if err != nil {
			fmt.Printf("Unable to find firewall for your search: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		_, err = client.RenameFirewall(firewall.ID, args[1])
		if err != nil {
			fmt.Printf("Unable to rename firewall: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": firewall.ID, "Name": firewall.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("The firewall called %s with ID %s was rename to %s\n", aurora.Green(firewall.Name), aurora.Green(firewall.ID), aurora.Green(args[1]))
		}
	},
}
