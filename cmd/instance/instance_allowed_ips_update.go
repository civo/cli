package instance

import (
	"fmt"
	"os"

	// "strings" // Uncomment if needed for string manipulations not covered by flags

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var allowedIPsUpdate []string

var instanceAllowedIPsUpdateCmd = &cobra.Command{
	Use:     "allowed-ips-update <INSTANCE_ID_OR_NAME>",
	Aliases: []string{"update-allowed-ips"},
	Short:   "Update the allowed IP addresses for an instance",
	Long: `Update the list of IP addresses that an instance is allowed to use for network traffic (IP/MAC spoofing protection).
Note: This replaces the existing list of allowed IPs. To clear all allowed IPs, pass an empty list with --ips "".`,
	Example: "civo instance allowed-ips-update my-instance --ips 192.168.0.10,10.0.0.5",
	Args:    cobra.ExactArgs(1),
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
			utility.Error("Finding instance %s", err)
			os.Exit(1)
		}

		// 'allowedIPsUpdate' is already a []string from the StringSliceVarP flag.
		// If --ips is not provided, allowedIPsUpdate will be an empty slice by default.
		// If --ips is provided as an empty string (e.g., --ips \"\"), it should also result in an empty slice,
		// which is useful for clearing the allowed IPs if the API supports it.

		resp, err := client.UpdateInstanceAllowedIPs(instance.ID, allowedIPsUpdate)
		if err != nil {
			utility.Error("Updating allowed IPs for instance %s: %s", instance.ID, err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", instance.ID, "ID")
		ow.AppendDataWithLabel("hostname", instance.Hostname, "Hostname")
		ow.AppendDataWithLabel("result", string(resp.Result), "Result")

		if common.OutputFormat == "human" {
			if resp.Result == "success" { // Assuming SimpleResponse has a "Result" field that indicates success
				fmt.Printf("Allowed IPs for instance %s (%s) updated successfully.\n", utility.Green(instance.Hostname), instance.ID)
				if len(allowedIPsUpdate) > 0 {
					fmt.Printf("New allowed IPs: %v\n", allowedIPsUpdate)
				} else {
					fmt.Println("All allowed IPs have been cleared.")
				}
			} else {
				fmt.Printf("Failed to update allowed IPs for instance %s (%s). Result: %s\n", utility.Red(instance.Hostname), instance.ID, resp.Result)
			}
		} else {
			ow.WriteSingleObjectJSON(common.PrettySet)
		}
	},
}
