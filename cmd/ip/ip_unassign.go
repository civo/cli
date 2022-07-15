package ip

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

var ipUnassignCmd = &cobra.Command{
	Use:     "unassign",
	Aliases: []string{"detach"},
	Example: `civo ip unassign  127.0.0.1
civo ip Unassign server-1
civo ip Unassign <ip id>`,

	Short: "Unassign IP address to an instance",
	Args:  cobra.MinimumNArgs(1),
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

		ip, err := client.FindIP(args[0])
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry there is no %s IP in your account", utility.Red(args[0]))
				os.Exit(1)
			} else if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one IP with that value in your account")
				os.Exit(1)
			} else {
				utility.Error("%s", err)
				os.Exit(1)
			}
		}
		if utility.UserConfirmedUnassign("ip", common.DefaultYes, ip.Name) {
			_, err = client.UnassignIP(ip.ID, client.Region)
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
				fmt.Printf("Unassigned IP %s from Civo resource\n", utility.Green(ip.Name))
			}
		} else {
			fmt.Println("Aborted")
		}

	},
}
