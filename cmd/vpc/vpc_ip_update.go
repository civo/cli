package vpc

import (
	"errors"
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var vpcIPUpdateName string

var vpcIPUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"rename", "change"},
	Example: `civo vpc ip update IP_NAME/IP_ID --name NEW_NAME`,
	Short:   "Update a VPC reserved IP",
	Args:    cobra.MinimumNArgs(1),
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

		ip, err := client.FindVPCIP(args[0])
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry there is no %s VPC IP in your account", utility.Red(args[0]))
				os.Exit(1)
			}
			if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one VPC IP with that value in your account")
				os.Exit(1)
			}
			utility.Error("%s", err)
			os.Exit(1)
		}

		ip, err = client.UpdateVPCIP(ip.ID, &civogo.UpdateIPRequest{
			Name:   vpcIPUpdateName,
			Region: client.Region,
		})
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Renamed VPC IP to %s with ID %s\n", utility.Green(vpcIPUpdateName), utility.Green(ip.ID))
		}
	},
}
