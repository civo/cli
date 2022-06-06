package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var appShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Example: `civo app show APP-NAME"`,
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show Application information",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		application, err := client.FindApplication(args[0])
		if err != nil {
			utility.Error("App %s", err)
			os.Exit(1)
		}

		networks, err := client.ListNetworks()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		ow.AppendData("ID", application.ID)
		ow.AppendDataWithLabel("name", application.Name, "Name")

		for _, network := range networks {
			if application.NetworkID == network.ID {
				ow.AppendDataWithLabel("network_name", network.Label, "Network Name")
			}
		}

		ow.AppendDataWithLabel("region", client.Region, "Region")
		ow.AppendDataWithLabel("description", application.Description, "Description")
		ow.AppendDataWithLabel("image", application.Image, "Image")
		ow.AppendDataWithLabel("size", application.Size, "Size")
		//ow.AppendDataWithLabel("domains", strings.Join(application.Domains, ", "), "Domains")
		ow.AppendDataWithLabel("status", application.Status, "Status")

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteKeyValues()
		}

		fmt.Println()
		//ow.WriteHeader("Domains ")
		owDomain := utility.NewOutputWriter()
		for _, domain := range application.Domains {
			owDomain.StartLine()
			owDomain.AppendData("Domains :", domain)
		}
		fmt.Println()
		owDomain.WriteTable()
		fmt.Println()

		fmt.Println()
		ow.WriteHeader("Application Config ")
		owConfig := utility.NewOutputWriter()
		for _, config := range application.Config {
			owConfig.StartLine()
			owConfig.AppendData("Name", config.Name)
			owConfig.AppendData("Value", config.Value)
		}
		fmt.Println()
		owConfig.WriteTable()
		fmt.Println()

		if application.ProcessInfo != nil {
			ow.WriteHeader("Application Processes ")
			owProcess := utility.NewOutputWriter()
			for _, process := range application.ProcessInfo {
				owProcess.StartLine()
				owProcess.AppendData("Type", process.ProcessType)
				owProcess.AppendData("Count", strconv.Itoa(process.ProcessCount))
			}
			fmt.Println()
			owProcess.WriteTable()
			fmt.Println()
		}

	},
}
