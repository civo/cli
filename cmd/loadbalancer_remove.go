package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
)

var loadBalancerRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Short:   "Remove a Load Balancer",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		if utility.AskForConfirmDelete("load balancer") == nil {
			lb, err := client.FindLoadBalancer(args[0])
			if err != nil {
				utility.Error("Unable to find load balancer for your search %s", err)
				os.Exit(1)
			}

			_, err = client.DeleteLoadBalancer(lb.ID)

			ow := utility.NewOutputWriterWithMap(map[string]string{"ID": lb.ID, "Hostname": lb.Hostname})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The load blancer %s with ID %s was delete\n", utility.Green(lb.Hostname), utility.Green(lb.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
