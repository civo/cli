package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var instancePublicIPCmd = &cobra.Command{
	Use:     "public-ip",
	Example: "civo instance public-ip ID/HOSTNAME",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show instance's public IP [deprectated]",
	Aliases: []string{"ip", "publicip"},
	Long: `Show the specified instance's public IP by part of the instance's ID or name.
If you wish to use a custom format, the available fields are:

	* id
	* hostname
	* public_ip

This command is deprecated and instead you should use:

civo instance show ID/HOSTNAME -o custom -f "PublicIP"`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

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

		if outputFormat == "human" {
			fmt.Printf("The instance %s (%s) has the public IP %s\n", utility.Green(instance.Hostname), instance.ID, utility.Green(instance.PublicIP))
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendDataWithLabel("id", instance.ID, "ID")
			ow.AppendDataWithLabel("hostname", instance.Hostname, "Hostname")
			ow.AppendDataWithLabel("public_ip", instance.PublicIP, "Public ID")
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON(prettySet)
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
