package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"os"
	"time"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var waitStop bool
var instanceStopCmd = &cobra.Command{
	Use:     "stop",
	Example: "civo instance stop ID/HOSTNAME",
	Short:   "Stop an instance",
	Aliases: []string{"shutdown"},
	Long: `Pull the power from the specified instance by part of the ID or name.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("Finding instance: %s\n", err)
			os.Exit(1)
		}

		_, err = client.StopInstance(instance.ID)
		if err != nil {
			utility.Error("Stopping instance %s", err)
			os.Exit(1)
		}

		if waitStop == true {

			stillStopping := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = "Stopping instance... "
			s.Start()

			for stillStopping {
				instanceCheck, err := client.FindInstance(instance.ID)
				if err != nil {
					utility.Error("Finding instance: %s\n", err)
					os.Exit(1)
				}
				if instanceCheck.Status == "SHUTOFF" {
					stillStopping = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}
		}

		if outputFormat == "human" {
			fmt.Printf("The instance %s (%s) is being stopped\n", utility.Green(instance.Hostname), instance.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
