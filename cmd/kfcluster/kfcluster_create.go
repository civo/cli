package kfcluster

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var networkID, firewallID, size string

var kfcCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo kfcluster create <KFCLUSTER-NAME> --network <NETWORK_ID> --size <SIZE> --firewall <FIREWALL_ID>",
	Short:   "Create a new kubeflow cluster",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		var network *civogo.Network
		if networkID != "" {
			network, err = client.FindNetwork(networkID)
			if err != nil {
				utility.Error("the network %s doesn't exist", networkID)
				os.Exit(1)
			}
		}

		if firewallID != "" {
			_, err = client.FindFirewall(firewallID)
			if err != nil {
				utility.Error("the firewall %s doesn't exist", firewallID)
				os.Exit(1)
			}
		}

		sizes, err := client.ListInstanceSizes()
		if err != nil {
			utility.Error("Unable to list sizes %s", err)
			os.Exit(1)
		}

		sizeIsValid := false
		if size != "" {
			for _, s := range sizes {
				if s.Name == size && utility.SizeType(s.Name) == "KfCluster" {
					sizeIsValid = true
					break
				}
			}
			if !sizeIsValid {
				utility.Error("The provided size is not valid")
				os.Exit(1)
			}
		}

		kfCluster := civogo.CreateKfClusterReq{
			Name:       args[0],
			NetworkID:  network.ID,
			Size:       size,
			FirewallID: firewallID,
			Region:     client.Region,
		}

		kfc, err := client.CreateKfCluster(kfCluster)
		if err != nil {
			utility.Error("Error creating kfcluster %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": kfc.ID, "name": kfc.Name})
		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("KfCluster (%s) with ID %s has been created\n", utility.Green(kfc.Name), kfc.ID)
		}
	},
}
