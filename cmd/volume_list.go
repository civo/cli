package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var volumeListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo volume ls`,
	Short:   "List volumes",
	Long: `List all available volumes.
If you wish to use a custom format, the available fields are:

	* id
	* name
	* instance_id
	* size_gigabytes
	* mount_point
	* bootable
	* created_at

Example: civo volume ls -o custom -f "ID: Name (SizeGigabytes)`,
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

		volumes, err := client.ListVolumes()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		for _, volume := range volumes {
			ow.StartLine()
			ow.AppendDataWithLabel("id", volume.ID, "ID")
			ow.AppendDataWithLabel("name", volume.Name, "Name")

			if volume.InstanceID != "" {
				instance, err := client.FindInstance(volume.InstanceID)
				if err != nil {
					utility.Error("Finding the instance failed with %s", err)
					os.Exit(1)
				}
				ow.AppendDataWithLabel("instance_id", instance.Hostname, "Instance")
			} else {
				ow.AppendDataWithLabel("instance_id", "-", "Instance")
			}

			ow.AppendDataWithLabel("size_gigabytes", fmt.Sprintf("%s GB", strconv.Itoa(volume.SizeGigabytes)), "Size")
			ow.AppendDataWithLabel("mount_point", volume.MountPoint, "Mount Point")
			ow.AppendDataWithLabel("bootable", strconv.FormatBool(volume.Bootable), "Bootable")
			ow.AppendDataWithLabel("created_at", volume.CreatedAt.Format(time.RFC1123), "Created At")
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
