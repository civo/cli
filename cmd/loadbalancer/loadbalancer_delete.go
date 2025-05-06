package loadbalancer

import (
	"errors"
	"fmt"
	"github.com/civo/cli/common"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var loadBalancerRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo loadbalancer rm ID/name",
	Short:   "Remove a load balancer",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		lb, err := client.FindLoadBalancer(args[0])
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry there is no %s load balancer in your account", utility.Red(args[0]))
				os.Exit(1)
			}
			if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one load balancer with that name in your account")
				os.Exit(1)
			}
		}

		if utility.UserConfirmedDeletion("load balancer", common.DefaultYes, lb.Name) {

			_, err = client.DeleteLoadBalancer(lb.ID)
			if err != nil {
				utility.Error("error deleting the loadbalancer: %s", err)
				os.Exit(1)
			}

			ow := utility.NewOutputWriterWithMap(map[string]string{"id": lb.ID, "hostname": lb.Name})

			switch common.OutputFormat {
			case "json":
				ow.WriteSingleObjectJSON(common.PrettySet)
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The load balancer %s with ID %s was deleted\n", utility.Green(lb.Name), utility.Green(lb.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
