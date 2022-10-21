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

var ipDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"unallocate", "free", "remove", "rm"},
	Example: `civo ip delete 127.0.0.1 
civo ip delete server-1
civo ip delete <ip id>
`,
	Short: "Delete a reserved ip",
	Long: `Once a reserved IP is deleted, it cannot be obtained again.
Please make sure to delete your domains aren't pointed to it before deleting it.`,
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

		if utility.UserConfirmedDeletion("IP", common.DefaultYes, ip.Name) {
			_, err = client.DeleteIP(ip.ID)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}

			ow := utility.NewOutputWriterWithMap(map[string]string{"id": ip.ID, "name": ip.Name})

			switch common.OutputFormat {
			case "json":
				ow.WriteSingleObjectJSON(common.PrettySet)
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("IP called %s with ID %s was deleted\n", utility.Green(ip.Name), utility.Green(ip.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
