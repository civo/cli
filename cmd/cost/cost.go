package main

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var currentcost = &cobra.Command{
	Use:     "cost",
	Aliases: []string{"show"},
	Short:   "Show the current cost",
	Example: "civo cost",

	Run: func(cmd *cobra.Command, args []string) {

		// Ensure the user has set a region
		utility.EnsureCurrentRegion()

		// Create a new client to the Civo API using the user's API key
		client, err := config.CivoAPIClient()
		// If the user has set a region, use that
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		// If there is an error, display it and exit
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// client is already defined, so we can just call the GetCost() method
		cost, err := client.GetCost()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Display the current cost
		fmt.Printf("cost: %s\n", cost)
	},
}

func init() {
	costCmd.AddCommand(currentcost)
	// Add the --region flag to the currentcost command
	currentcost.Flags().StringVarP(&common.RegionSet, "region", "r", "", "The region to use")

}
