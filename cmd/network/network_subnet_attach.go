package network

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instanceID, subnetID string
var waitSubnetAttach bool

var networkSubnetAttachCmd = &cobra.Command{
	Use:     "attach",
	Short:   "Attach a subnet to a resource",
	Example: "civo network subnet attach --instance INSTANCE-ID --subnet SUBNET-ID",
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

		instance, err := client.FindInstance(instanceID)
		if err != nil {
			utility.Error("Instance %s", err)
			os.Exit(1)
		}

		subnet, err := client.FindSubnet(subnetID, instance.NetworkID)
		if err != nil {
			utility.Error("Subnet %s", err)
			os.Exit(1)
		}

		_, err = client.AttachSubnetToInstance(instance.NetworkID, subnet.ID, &civogo.CreateRoute{
			ResourceID:   instance.ID,
			ResourceType: "instance",
		})
		if err != nil {
			utility.Error("error attaching the subnet: %s", err)
			os.Exit(1)
		}

		if waitSubnetAttach {
			stillAttaching := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = "Attaching subnet to the instance... "
			s.Start()

			for stillAttaching {
				subnetCheck, err := client.FindSubnet(subnet.ID, instance.NetworkID)
				if err != nil {
					utility.Error("Finding the subnet failed with %s", err)
					os.Exit(1)
				}
				if subnetCheck.Status == "attached" {
					stillAttaching = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": subnet.ID, "name": subnet.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The subnet (%s) was attached to the instance with ID (%s)\n", utility.Green(subnet.Name), utility.Green(instance.ID))
		}
	},
}
