package vpc

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var vpcIPName string

var vpcIPCreateCmd = &cobra.Command{
	Use:     "reserve",
	Aliases: []string{"new", "add", "allocate", "create"},
	Example: `civo vpc ip reserve
civo vpc ip reserve --name "my-ip"`,
	Short: "Reserve a new VPC IP",
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

		createReq := &civogo.CreateIPRequest{
			Region: client.Region,
		}
		if vpcIPName != "" {
			createReq.Name = vpcIPName
		}

		ip, err := client.NewVPCIP(createReq)
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
			if vpcIPName != "" {
				fmt.Printf("Reserved VPC IP called %s with ID %s\n", utility.Green(vpcIPName), utility.Green(ip.ID))
			} else {
				fmt.Printf("Reserved VPC IP with ID %s\n", utility.Green(ip.ID))
			}
		}
	},
}
