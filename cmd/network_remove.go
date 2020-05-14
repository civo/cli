package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
)

var networkRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo network rm NAME",
	Short:   "Remove a network",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		if utility.AskForConfirmDelete("network") == nil {
			network, err := client.FindNetwork(args[0])
			if err != nil {
				utility.Error("Unable to find network for your search %s", err)
				os.Exit(1)
			}

			_, err = client.DeleteNetwork(network.ID)

			ow := utility.NewOutputWriterWithMap(map[string]string{"ID": network.ID, "Name": network.Name, "Label": network.Label})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The network called %s with ID %s was delete\n", utility.Green(network.Label), utility.Green(network.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}

	},
}
