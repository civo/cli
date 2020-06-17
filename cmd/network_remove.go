package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
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
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if utility.UserConfirmedDeletion("network", defaultYes) == true {
			network, err := client.FindNetwork(args[0])
			if err != nil {
				utility.Error("Finding the network for your search failed with %s", err)
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
				fmt.Printf("The network called %s with ID %s was deleted\n", utility.Green(network.Label), utility.Green(network.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}

	},
}
