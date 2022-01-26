package cmd

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var dbCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo db create --size g3.medium  hello",
	Short:   "Create a new database",
	Run: func(cmd *cobra.Command, args []string) {
		reg, err := utility.GetCurrentRegion()
		if err != nil {
			utility.Error("Failed to get region: %s", err)
			os.Exit(1)
		}
		// TODO: Should we have this check or will all regions support DBaaS?
		// check, region, err := utility.CheckAvailability("db", reg)
		// if err != nil {
		// 	utility.Error("Error checking availability %s", err)
		// 	os.Exit(1)
		// }
		// if !check {
		// 	utility.Error("Sorry you can't create a DBaaS instance in the %s region", region)
		// 	os.Exit(1)
		// }

		client, err := config.CivoAPIClient()
		if reg != "" {
			client.Region = reg
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		// TODO: Enable this check once size configmap gets updated with DB sizes
		// if !strings.Contains(size, "db") {
		// 	dbSize, err := utility.GetDBSizes()
		// 	if err != nil {
		// 		utility.Error("Error %s", err)
		// 		os.Exit(1)
		// 	}

		// 	utility.Error("You can create a cluster with this %s size, Possible values: %s", size, dbSize)
		// 	os.Exit(1)
		// }

		if len(args) > 0 {
			if utility.ValidNameLength(args[0]) {
				utility.Warning("the cluster name cannot be longer than 63 characters")
				os.Exit(1)
			}
			dbName = args[0]
		} else {
			dbName = utility.RandomName()
		}

		var network = &civogo.Network{}
		if networkID == "default" {
			network, err = client.GetDefaultNetwork()
			if err != nil {
				utility.Error("Network %s", err)
				os.Exit(1)
			}
		} else {
			network, err = client.FindNetwork(dbNetworkID)
			if err != nil {
				utility.Error("Network %s", err)
				os.Exit(1)
			}
		}

		configDB := civogo.CreateDatabaseRequest{
			Name:                 dbName,
			Size:                 dbSize,
			Software:             software,
			SoftwareVersion:      softwareVersion,
			NetworkID:            network.ID,
			Replicas:             replicas,
			NumSnapshotsToRetain: numSnapshots,
			PublicIPRequired:     publicIP,
			// TODO: Fetch firewall ID
			// FirewallID: ,
		}

		db, err := client.NewDatabase(&configDB)
		if err != nil {
			utility.Error("Error creating database %s", err)
			os.Exit(1)
		}

		// TODO: Check how to use these variables. Common package between root and sub command packages?
		// ow := utility.NewOutputWriterWithMap(map[string]string{"id": db.ID, "name": db.Name})
		// switch outputFormat {
		// case "json":
		// 	ow.WriteSingleObjectJSON(prettySet)
		// case "custom":
		// 	ow.WriteCustomOutput(outputFields)
		// default:
		// 	fmt.Printf("Database %s (%s) has been created\n", utility.Green(db.Name), db.ID)
		// }
		fmt.Printf("Database %s (%s) has been created\n", utility.Green(db.Name), db.ID)

	},
}
