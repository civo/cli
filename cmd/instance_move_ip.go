package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var instanceMoveIPCmd = &cobra.Command{
	Use:     "move-ip",
	Aliases: []string{"switch-ip", "moveip", "switchip"},
	Short:   "Move a public IP",
	Long: `Move a public IP address to a target instance by part of the ID or name.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname
	* PublicIP

Example: civo instance move-ip ID/NAME 1.2.3.4`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Printf("You must specify %d parameters (you gave %d), the ID/name and the public IP\n", aurora.Red(2), aurora.Red(len(args)))
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

		instances, err := client.ListAllInstances()
		var moving bool
		for _, i := range instances {
			if i.PublicIP == args[1] && i.ID != instance.ID {
				moving = true
				_, err = client.MovePublicIPToInstance(instance.ID, args[1])
				if err != nil {
					fmt.Printf("Moving IP: %s\n", aurora.Red(err))
					os.Exit(1)
				}
			}
		}

		if !moving {
			fmt.Printf("Unable to find public IP: %s\n", aurora.Red(args[1]))
			os.Exit(1)
		}

		if outputFormat == "human" {
			fmt.Printf("Moving the IP %s to the instance %s (%s)\n", aurora.Green(args[1]), aurora.Green(instance.Hostname), instance.ID)
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
