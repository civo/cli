package cmd

import (
	"os"
	"strconv"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instanceSizeCmd = &cobra.Command{
	Use:     "size",
	Example: `civo instance size"`,
	Aliases: []string{"sizes", "all"},
	Short:   "List instances size",
	Long:    `List all current instances size.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		filter := []civogo.InstanceSize{}
		sizes, err := client.ListInstanceSizes()
		if err != nil {
			utility.Error("%s", err)
			return
		}

		for _, size := range sizes {
			if !strings.Contains(size.Name, "db") && !strings.Contains(size.Name, "k3s") {
				filter = append(filter, size)
			}
		}

		ow := utility.NewOutputWriter()
		for _, size := range filter {
			ow.StartLine()
			ow.AppendData("Name", size.Name)
			ow.AppendData("Description", size.Description)
			ow.AppendData("Type", "Instance")
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
