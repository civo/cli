package ip

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var ipCreateCmd = &cobra.Command{
	Use:     "reserve",
	Aliases: []string{"new", "add", "allocate"},
	Example: `civo ip reserve 
civo ip reserve -n "server-1"`,
	Short: "Reserve a new ip",
	Long:  `You can name your ip with the --name flag.`,
	Args:  cobra.MinimumNArgs(0),
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

		createReq := civogo.CreateIPRequest{
			Region: client.Region,
		}
		if name != "" {
			createReq.Name = name
		}

		ip, err := client.NewIP(&createReq)
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
			if name != "" {
				fmt.Printf("Reserved IP called %s with ID %s\n", utility.Green(name), utility.Green(ip.ID))
			} else {
				fmt.Printf("Reserved IP with ID %s\n", utility.Green(ip.ID))
			}
		}
	},
}
