package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var instanceMoveIPCmd = &cobra.Command{
	Use:     "move-ip",
	Example: "civo instance move-ip ID/HOSTNAME 1.2.3.4",
	Args:    cobra.MinimumNArgs(2),
	Aliases: []string{"switch-ip", "moveip", "switchip"},
	Short:   "Move a public IP",
	Long: `Move a public IP address to a target instance by part of its ID or name.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname
	* PublicIP`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Printf("You must specify %s parameters (you gave %s), the ID/name and the public IP\n", utility.Red("2"), utility.Red(strconv.Itoa(len(args))))
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

		instances, err := client.ListAllInstances()
		if err != nil {
			utility.Error("error listing all instances: %s", err)
			os.Exit(1)
		}

		var moving bool
		for _, i := range instances {
			if i.PublicIP == args[1] && i.ID != instance.ID {
				moving = true
				_, err = client.MovePublicIPToInstance(instance.ID, args[1])
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
			}
		}

		if !moving {
			utility.Error("Unable to find that public IP connected to one of your instances", args[1])
			os.Exit(1)
		}

		if outputFormat == "human" {
			fmt.Printf("Moved the IP %s to the instance %s (%s)\n", utility.Green(args[1]), utility.Green(instance.Hostname), instance.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			ow.AppendDataWithLabel("PublicIP", args[1], "Public IP")
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
