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

var dbCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo db create <DATABASE-NAME> <DATABASE-SIZE> <SOFTWARE>",
	Short:   "Create a new database",
	Args:    cobra.MinimumNArgs(3),
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

		network, err := client.GetDefaultNetwork()
		if err != nil {
			utility.Error("Unable to find the default network: %s", err)
			os.Exit(1)
		}

		firewall, err := client.FindFirewall("default")
		if err != nil {
			utility.Error("Unable to find the default firewall: %s", err)
			os.Exit(1)
		}

		configDB := civogo.CreateDatabaseRequest{
			Name:                 args[0],
			Size:                 args[1],
			Software:             args[2],
			SoftwareVersion:      softwareVersion,
			NetworkID:            network.ID,
			Replicas:             replicas,
			NumSnapshotsToRetain: numSnapshots,
			PublicIPRequired:     publicIP,
			FirewallID:           firewall.ID,
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
			fmt.Printf("Database %s (%s) has been created\n", utility.Green(db.Name), db.ID)
		}
	},
}
