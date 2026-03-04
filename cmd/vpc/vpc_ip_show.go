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

var vpcIPShowCmd = &cobra.Command{
	Use:     "show [IP-NAME/IP-ID/IP-ADDRESS]",
	Short:   "Show details of a specific VPC reserved IP",
	Aliases: []string{"get", "describe", "inspect"},
	Args:    cobra.ExactArgs(1),
	Example: "civo vpc ip show IP_NAME",
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

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", ip.ID, "ID")
		ow.AppendDataWithLabel("name", ip.Name, "Name")
		ow.AppendDataWithLabel("address", ip.IP, "Address")

		switch common.OutputFormat {
		case "json":
			ow.ToJSON(ip, common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Println("VPC IP Details:")
			fmt.Printf("ID: %s\n", ip.ID)
			fmt.Printf("Name: %s\n", ip.Name)
			fmt.Printf("IP: %s\n", ip.IP)
			if ip.AssignedTo.ID != "" {
				fmt.Printf("Assigned To: %s (%s)\n", ip.AssignedTo.Name, ip.AssignedTo.Type)
			} else {
				fmt.Println("Assigned To: No resource")
			}
		}
	},
}
