package volumetype

import (
	"fmt"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

var volumetypesListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List available volume types",
	Long:  `List the available volume types in Civo cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		// Call ListVolumeTypes from SDK
		volumeTypes, err := client.ListVolumeTypes()
		if err != nil {
			fmt.Printf("Error fetching volume types: %s\n", err)
			return
		}

		ow := utility.NewOutputWriter()

		// Print the volume types
		for _, volumeType := range volumeTypes {
			ow.StartLine()
			ow.AppendDataWithLabel("name", volumeType.Name, "Name")
			ow.AppendDataWithLabel("description", volumeType.Description, "Description")
			ow.AppendDataWithLabel("default", strconv.FormatBool(volumeType.Enabled), "Enabled")
			ow.AppendDataWithLabel("tags", strings.Join(volumeType.Labels, " "), "Labels")
		}
		ow.FinishAndPrintOutput()
	},
}
