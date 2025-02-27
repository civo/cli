package instance

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instanceRecoveryStatusCmd = &cobra.Command{
	Use:     "recovery-status",
	Aliases: []string{"recovery"},
	Example: "civo instance recovery-status INSTANCE_ID/HOSTNAME",
	Short:   "Get the recovery status of an instance",
	Long: `Show the current recovery mode status for a specified instance.
If you wish to use a custom format, the available fields are:

	* id
	* hostname
	* status

Example:
  * Check recovery status:
    civo instance recovery-status my-instance`,
	Args: cobra.MinimumNArgs(1),
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

		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("Instance %s", err)
			os.Exit(1)
		}

		status, err := client.GetRecoveryStatus(instance.ID)
		if err != nil {
			utility.Error("Failed to get recovery status: %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{
			"id":       instance.ID,
			"hostname": instance.Hostname,
			"status":   string(status.Result),
		})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Recovery Status for instance %s (%s):\n", utility.Green(instance.Hostname), instance.ID)
			ow.WriteKeyValues()
		}
	},
}
