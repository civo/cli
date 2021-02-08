package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instanceSetFirewallCmd = &cobra.Command{
	Use:     "firewall",
	Aliases: []string{"set-firewall", "change-firewall", "fw"},
	Args:    cobra.MinimumNArgs(2),
	Example: "civo instance firewall HOSTNAME/INSTANCE_ID FIREWALL_NAME/FIREWALL_ID",
	Short:   "Set firewall for instance",
	Long: `Change an instance's firewall by part of the instance's ID or name and the full firewall ID.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname
	* FirewallID`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Printf("You must specify %s parameters (you gave %s), the ID/name and the firewall ID\n", utility.Red("2"), utility.Red(strconv.Itoa(len(args))))
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		firewall, err := client.FindFirewall(args[1])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		_, err = client.SetInstanceFirewall(instance.ID, args[1])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if outputFormat == "human" {
			fmt.Printf("Set the firewall for the instance %s (%s) to %s (%s)\n", utility.Green(instance.Hostname), instance.ID, utility.Green(firewall.Name), firewall.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			ow.AppendDataWithLabel("FirewallID", firewall.ID, "Firewall ID")
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
