package cmd

import (
	"strconv"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var sizeListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo size ls`,
	Short:   "List sizes",
	Long: `List all available sizes for instances or Kubernetes nodes.
If you wish to use a custom format, the available fields are:

	* Name
	* NiceName
	* CPUCores
	* RAMMegabytes
	* DiskGigabytes
	* Description
	* Selectable

Example: civo size ls -o custom -f "Code: Name (size)"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			return
		}

		sizes, err := client.ListInstanceSizes()
		if err != nil {
			utility.Error("%s", err)
			return
		}

		ow := utility.NewOutputWriter()

		for _, size := range sizes {
			ow.StartLine()
			ow.AppendData("Name", size.Name)
			ow.AppendData("Description", size.Description)
			ow.AppendData("CPU", strconv.Itoa(size.CPUCores))

			if outputFormat == "json" || outputFormat == "custom" {
				ow.AppendData("RAM_MB", strconv.Itoa(size.RAMMegabytes))
				ow.AppendData("Disk_GB", strconv.Itoa(size.DiskGigabytes))
			} else {
				ow.AppendData("RAM (MB)", strconv.Itoa(size.RAMMegabytes))
				ow.AppendData("Disk (GB)", strconv.Itoa(size.DiskGigabytes))
			}
			ow.AppendData("Selectable", utility.BoolToYesNo(size.Selectable))
		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
