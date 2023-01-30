package app

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var name, size, firewallID, processType string
var processCount int

var appUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"edit", "modify", "change", "scale", "resize"},
	Short:   "Update an App",
	Example: "civo app update APP_NAME --flags",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		findApp, err := client.FindApplication(args[0])
		if err != nil {
			utility.Error("App %s", err)
			os.Exit(1)
		}

		processInfo := civogo.ProcInfo{
			ProcessType:  processType,
			ProcessCount: processCount,
		}

		app, err := client.UpdateApplication(findApp.ID, &civogo.UpdateApplicationRequest{
			Name:        name,
			Size:        size,
			FirewallID:  firewallID,
			ProcessInfo: append(findApp.ProcessInfo, processInfo),
		})

		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": app.ID, "name": findApp.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The App with ID %s was updated. \n", utility.Green(app.ID))
			os.Exit(0)
		}
	},
}
