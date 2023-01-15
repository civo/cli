package database

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var rulesFirewall string
var dbCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo db create <DATABASE-NAME> --size <SIZE>",
	Short:   "Create a new database",
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
		} else {
			network, err = client.GetDefaultNetwork()
			if err != nil {
				utility.Error("Unable to find the default network: %s", err)
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
				if s.Name == size && utility.SizeType(s.Name) == "Database" {
					sizeIsValid = true
					break
				}
			}
			if !sizeIsValid {
				utility.Error("The provided size is not valid")
				os.Exit(1)
			}
		}

		configDB := civogo.CreateDatabaseRequest{
			Name:          args[0],
			Size:          size,
			NetworkID:     network.ID,
			Nodes:         nodes,
			FirewallID:    firewallID,
			FirewallRules: rulesFirewall,
			Region:        client.Region,
		}

		db, err := client.NewDatabase(&configDB)
		if err != nil {
			utility.Error("Error creating database %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": db.ID, "name": db.Name})
		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Database (%s) with ID %s has been created\n", utility.Green(db.Name), db.ID)
		}
	},
}
