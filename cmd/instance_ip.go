package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var instancePublicIPCmd = &cobra.Command{
	Use:     "public-ip",
	Example: "civo instance public-ip enable/disable ID/HOSTNAME",
	Args:    cobra.MinimumNArgs(2),
	Short:   "Enable/disable controls if instance should have a public IP",
	Aliases: []string{"ip", "publicip"},
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

		if args[0] != "enable" && args[0] != "disable" {
			utility.Error("Please specify either enable or disable")
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[1])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if args[0] == "disable" {
			instance.PublicIP = "false"
		} else {
			instance.PublicIP = "true"
		}

		_, err = client.UpdateInstance(instance)
		if err != nil {
			utility.Error("%s", err)
		}

		if args[0] == "disable" {
			fmt.Printf("Instance %s has been updated to NOT have a Public IP\n", utility.Green(instance.Hostname))
		} else {
			fmt.Printf("Instance %s has been updated to have a Public IP. IP addressed will be assigned shortly.\n", utility.Green(instance.Hostname))
		}
	},
}
