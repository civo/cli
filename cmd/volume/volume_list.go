package volume

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/cli/common"
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
	* network_id
	* cluster_id
	* instance_id
	* size_gigabytes
	* mount_point
	* status

Example: civo volume ls -o custom -f "ID: Name (SizeGigabytes)`,
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

			if volume.NetworkID != "" {
				network, err := client.FindNetwork(volume.NetworkID)
				if err != nil {
					utility.Error("Finding the network failed with %s", err)
					os.Exit(1)
				}
				ow.AppendDataWithLabel("network_id", network.Label, "Network")
			} else {
				ow.AppendDataWithLabel("network_id", "", "Network")
			}

			if volume.ClusterID != "" {
				cluster, err := client.FindKubernetesCluster(volume.ClusterID)
				if err != nil {
					utility.Error("Finding the cluster failed with %s", err)
					os.Exit(1)
				}
				ow.AppendDataWithLabel("cluster_id", cluster.Name, "Cluster")
			} else {
				ow.AppendDataWithLabel("cluster_id", "", "Cluster")
			}

			if volume.InstanceID != "" {
				instance, err := client.FindInstance(volume.InstanceID)
				if err != nil {
					utility.Error("Finding the instance failed with %s", err)
					os.Exit(1)
				}
				ow.AppendDataWithLabel("instance_id", instance.Hostname, "Instance")
			} else {
				ow.AppendDataWithLabel("instance_id", "", "Instance")
			}

			ow.AppendDataWithLabel("size_gigabytes", fmt.Sprintf("%s GB", strconv.Itoa(volume.SizeGigabytes)), "Size")
			ow.AppendDataWithLabel("mount_point", volume.MountPoint, "Mount Point")
			ow.AppendDataWithLabel("status", volume.Status, "Status")
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
