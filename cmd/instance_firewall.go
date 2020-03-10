package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var instanceSetFirewallCmd = &cobra.Command{
	Use:     "firewall",
	Aliases: []string{"set-firewall", "change-firewall"},
	Short:   "Use different firewall",
	Long: `Change an instance's firewall by part of the instance's ID or name and the full firewall ID.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname
	* FirewallID

Example: civo instance firewall ID/NAME 12345`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Printf("You must specify %d parameters (you gave %d), the ID/name and the firewall ID\n", aurora.Red(2), aurora.Red(len(args)))
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			fmt.Printf("Finding instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		firewall, err := client.FindFirewall(args[1])
		if err != nil {
			fmt.Printf("Finding firewall: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		_, err = client.SetInstanceFirewall(instance.ID, args[1])
		if err != nil {
			fmt.Printf("Setting firewall: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		if OutputFormat == "human" {
			fmt.Printf("Setting the firewall for the instance %s (%s) to %s (%s)\n", aurora.Green(instance.Hostname), instance.ID, aurora.Green(firewall.Name), firewall.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			ow.AppendDataWithLabel("FirewallID", firewall.ID, "Firewall ID")
			if OutputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(OutputFields)
			}
		}
	},
}
