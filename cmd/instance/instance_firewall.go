package instance

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/cli/common"
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

	* id
	* hostname
	* firewall_id`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		if len(args) != 2 {
			fmt.Printf("You must specify %s parameters (you gave %s), the ID/name and the firewall ID\n", utility.Red("2"), utility.Red(strconv.Itoa(len(args))))
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		instanceArg := args[0]
		firewallArg := args[1]

		instance, err := client.FindInstance(instanceArg)
		if err != nil {
			utility.Error("Instance %s", err)
			os.Exit(1)
		}

		firewall, err := client.FindFirewall(firewallArg)
		if err != nil {
			utility.Error("Firewall %s", err)
			os.Exit(1)
		}

		instanceNetwork, err := client.FindNetwork(instance.NetworkID)
		if err != nil {
			utility.Error("Unable to find instance's network - %s", err)
			os.Exit(1)
		}

		firewallNetwork, err := client.FindNetwork(firewall.NetworkID)
		if err != nil {
			utility.Error("Unable to find firewall's network - %s", err)
			os.Exit(1)
		}

		if instanceNetwork.ID != firewallNetwork.ID {
			utility.Error("%q firewall is located in %q network but %q instance is located in %q network. Please try again with a firewall from %q network.", firewallArg, firewallNetwork.Label, instanceArg, instanceNetwork.Label, instanceNetwork.Label)
			os.Exit(1)
		}

		_, err = client.SetInstanceFirewall(instance.ID, firewall.ID)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if common.OutputFormat == "human" {
			fmt.Printf("Set the firewall for the instance %s (%s) to %s (%s)\n", utility.Green(instance.Hostname), instance.ID, utility.Green(firewall.Name), firewall.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendDataWithLabel("id", instance.ID, "ID")
			ow.AppendDataWithLabel("hostname", instance.Hostname, "Hostname")
			ow.AppendDataWithLabel("firewall_id", firewall.ID, "Firewall ID")
			if common.OutputFormat == "json" {
				ow.WriteSingleObjectJSON(common.PrettySet)
			} else {
				ow.WriteCustomOutput(common.OutputFields)
			}
		}
	},
}
