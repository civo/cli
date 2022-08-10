package volume

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/civogo"
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

		networks, err := client.ListNetworks()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		instances, err := client.ListAllInstances()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		clusters, err := client.ListKubernetesClusters()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		for _, volume := range volumes {
			ow.StartLine()
			ow.AppendDataWithLabel("id", volume.ID, "ID")
			ow.AppendDataWithLabel("name", volume.Name, "Name")

			var network civogo.Network

			if volume.NetworkID != "" {
				for _, network = range networks {
					if network.ID == volume.NetworkID {
						break
					}
				}
				ow.AppendDataWithLabel("network_id", network.Label, "Network")
			} else {
				ow.AppendDataWithLabel("network_id", "", "Network")
			}

			var cluster *civogo.KubernetesCluster
			if volume.ClusterID != "" {
				for _, c := range clusters.Items {
					if c.ID == volume.ClusterID {
						cluster = &c
						break
					}
				}
				ow.AppendDataWithLabel("cluster_id", cluster.Name, "Cluster")
			} else {
				ow.AppendDataWithLabel("cluster_id", "", "Cluster")
			}

			if volume.InstanceID != "" {
				if cluster != nil {
					for _, instance := range cluster.Instances {
						if instance.ID == volume.InstanceID {
							ow.AppendDataWithLabel("instance_id", instance.Hostname, "Instance")
							break
						}
					}
				} else {
					var instance civogo.Instance
					for _, instance = range instances {
						if instance.ID == volume.InstanceID {
							break
						}
					}

					ow.AppendDataWithLabel("instance_id", instance.Hostname, "Instance")
				}
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
