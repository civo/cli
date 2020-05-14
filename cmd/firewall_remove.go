package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
)

var firewallRemoveCmd = &cobra.Command{
	Use:     "remove [NAME]",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo firewall remove NAME",
	Short:   "Remove a firewall",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		if utility.AskForConfirmDelete("firewall") == nil {
			firewall, err := client.FindFirewall(args[0])
			if err != nil {
				utility.Error("Unable to find firewall for your search %s", err)
				os.Exit(1)
			}

			_, err = client.DeleteFirewall(firewall.ID)

			ow := utility.NewOutputWriterWithMap(map[string]string{"ID": firewall.ID, "Name": firewall.Name})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The firewall %s with ID %s was delete\n", utility.Green(firewall.Name), utility.Green(firewall.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
