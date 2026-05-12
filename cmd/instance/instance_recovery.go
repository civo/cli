package instance

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instanceRecoveryCmd = &cobra.Command{
	Use:     "recovery",
	Example: "civo instance recovery enable/disable ID/HOSTNAME",
	Args:    cobra.MinimumNArgs(2),
	Short:   "Enable/disable recovery mode for an instance",
	Long: `Enable or disable recovery mode for a specified instance.
Recovery mode allows you to boot your instance into a recovery environment for troubleshooting.

Example:
  * Enable recovery mode:
    civo instance recovery enable my-instance
  * Disable recovery mode:
    civo instance recovery disable my-instance`,
	Aliases: []string{"rescue"},
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
		if args[0] == "enable" {
			_, err := client.EnableRecoveryMode(instance.ID)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
		} else {
			_, err := client.DisableRecoveryMode(instance.ID)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", instance.ID, "ID")
		ow.AppendDataWithLabel("hostname", instance.Hostname, "Hostname")
		ow.AppendDataWithLabel("recovery_mode", args[0], "Recovery Mode")

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if args[0] == "enable" {
				fmt.Printf("Recovery mode has been enabled for instance %s\n", utility.Green(instance.Hostname))
			} else {
				fmt.Printf("Recovery mode has been disabled for instance %s\n", utility.Green(instance.Hostname))
			}
		}
	},
}
