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

var ipRenameCmd = &cobra.Command{
	Use:     "rename",
	Aliases: []string{"update", "change"},
	Example: `civo ip rename 127.0.0.1 server-2
civo ip rename server-1 server-2
civo ip rename <ip id> server-2 `,
	Short: "Rename reserved ip",
	Args:  cobra.MinimumNArgs(2),
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
			}
		}

		rename := args[1]
		RenameReq := civogo.UpdateIPRequest{
			Name:   rename,
			Region: client.Region,
		}

		ip, err = client.UpdateIP(ip.ID, &RenameReq)
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
			fmt.Printf("Renamed IP to %s with ID %s\n", utility.Green(rename), utility.Green(ip.ID))
		}
	},
}
