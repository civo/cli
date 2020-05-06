package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
)

var networkUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"change", "modify"},
	Short:   "Update a new network",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		oldNetwork, err := client.FindNetwork(args[0])
		if err != nil {
			fmt.Printf("Unable to find the network: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		network, err := client.RenameNetwork(args[1], oldNetwork.ID)
		if err != nil {
			fmt.Printf("Unable to update the network: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": network.ID, "Label": network.Label})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Update a network called %s with ID %s to %s\n", aurora.Green(oldNetwork.Label), aurora.Green(network.ID), aurora.Green(network.Label))
		}
	},
}
