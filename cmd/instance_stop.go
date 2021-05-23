package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var waitStop bool
var instanceStopCmd = &cobra.Command{
	Use:     "stop",
	Example: "civo instance stop ID/HOSTNAME",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Stop an instance",
	Aliases: []string{"shutdown"},
	Long: `Pull the power from the specified instance by part of the ID or name.
If you wish to use a custom format, the available fields are:

	* id
	* hostname`,
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

		_, err = client.StopInstance(instance.ID)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if waitStop {
			stillStopping := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = "Stopping instance... "
			s.Start()

			for stillStopping {
				instanceCheck, err := client.FindInstance(instance.ID)
				if err != nil {
					utility.Error("Finding instance failed with %s\n", err)
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
			ow.AppendDataWithLabel("id", instance.ID, "ID")
			ow.AppendDataWithLabel("hostname", instance.Hostname, "Hostname")
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON(prettySet)
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
