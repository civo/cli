package cmd

import (
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var appListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List all appplications",
	Example: "civo app ls",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()
		client, err := config.CivoAPIClient()

		if regionSet != "" {
			client.Region = regionSet
		}

		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		applications, err := client.ListApplications()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, app := range applications.Items {
			ow.StartLine()

			ow.AppendDataWithLabel("process_type", app.ProcessInfo[0].ProcessType, "Process Type")
			ow.AppendDataWithLabel("process_count", string(rune(app.ProcessInfo[0].ProcessCount)), "Process Count")
		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
