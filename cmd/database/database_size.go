package database

import (
	"os"
	"strconv"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var dbSizeCmd = &cobra.Command{
	Use:     "size",
	Example: "civo db size ls",
	Aliases: []string{"sizes", "all"},
	Short:   "List Database node sizes",
	Long:    `List all currently available Database node sizes.`,
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

		filter := []civogo.InstanceSize{}
		sizes, err := client.ListInstanceSizes()
		if err != nil {
			utility.Error("%s", err)
			return
		}

		for _, size := range sizes {
			if strings.Contains(size.Name, ".db.") {
				filter = append(filter, size)
			}
		}

		ow := utility.NewOutputWriter()
		for _, size := range filter {
			ow.StartLine()
			ow.AppendDataWithLabel("name", size.Name, "Name")
			ow.AppendDataWithLabel("description", size.Description, "Description")
			ow.AppendDataWithLabel("type", "Database", "Type")
			ow.AppendDataWithLabel("cpu_cores", strconv.Itoa(size.CPUCores), "CPU Cores")
			ow.AppendDataWithLabel("ram_mb", strconv.Itoa(size.RAMMegabytes), "RAM MB")
			ow.AppendDataWithLabel("disk_gb", strconv.Itoa(size.DiskGigabytes), "SSD GB")
			ow.AppendDataWithLabel("selectable", utility.BoolToYesNo(size.Selectable), "Selectable")
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
