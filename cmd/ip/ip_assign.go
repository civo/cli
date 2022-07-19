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

var ipAssignCmd = &cobra.Command{
	Use:     "assign",
	Aliases: []string{"attach"},
	Example: `civo ip assign  127.0.0.1 --instance <NAME>
civo ip assign server-1 --instance <instance ID>
civo ip assign <ip id> --instance <instance ID>`,

	Short: "Assign IP address to an instance",
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
			}
		}

		instance, err := client.FindInstance(instance)
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry there is no %s instance in your account", utility.Red(args[0]))
				os.Exit(1)
			} else if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one instance with that value in your account")
				os.Exit(1)
			} else {
				utility.Error("%s", err)
			}
		}

		_, err = client.AssignIP(ip.ID, instance.ID, "instance", client.Region)
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
			fmt.Printf("Assigned IP %s to instance %s\n", utility.Green(ip.Name), utility.Green(instance.Hostname))
		}
	},
}
