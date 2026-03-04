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

var vpcIPUnassignCmd = &cobra.Command{
	Use:     "unassign",
	Aliases: []string{"detach"},
	Example: `civo vpc ip unassign IP_NAME`,
	Short:   "Unassign a VPC IP from a resource",
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

		if utility.UserConfirmedUnassign("VPC IP", common.DefaultYes, ip.Name) {
			_, err = client.UnassignVPCIP(ip.ID, client.Region)
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
				fmt.Printf("Unassigned VPC IP %s from Civo resource\n", utility.Green(ip.Name))
			}
		} else {
			fmt.Println("Aborted")
		}
	},
}
