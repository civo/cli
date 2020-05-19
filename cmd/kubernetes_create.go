package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

var numTargetNodes int
var waitKubernetes bool
var (
	kubernetesVersion string
	targetNodesSize   string
)

var clusterName string

var kubernetesCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo kubernetes create CLUSTER_NAME [flags]",
	Short:   "Create a new kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		if len(args) > 0 {
			clusterName = args[0]
		} else {
			clusterName = utility.RandomName()
		}

		configKubernetes := &civogo.KubernetesClusterConfig{
			Name:              clusterName,
			NumTargetNodes:    numTargetNodes,
			TargetNodesSize:   targetNodesSize,
			KubernetesVersion: kubernetesVersion,
		}

		kubernetesCluster, err := client.NewKubernetesClusters(configKubernetes)
		if err != nil {
			utility.Error("Unable to create a kubernetes cluster %s", err)
			os.Exit(1)
		}

		if waitKubernetes == true {

			stillCreating := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = fmt.Sprintf("Creating a %s node k3s cluster of %s instances called %s... ", strconv.Itoa(kubernetesCluster.NumTargetNode), kubernetesCluster.TargetNodeSize, kubernetesCluster.Name)
			s.Start()

			for stillCreating {
				kubernetesCheck, err := client.FindKubernetesCluster(kubernetesCluster.ID)
				if err != nil {
					utility.Error("Unable to find the kubernetes cluster %s", err)
					os.Exit(1)
				}
				if kubernetesCheck.Status == "ACTIVE" {
					stillCreating = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": kubernetesCluster.ID, "Name": kubernetesCluster.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("The cluster %s (%s) has been created and took %s\n", utility.Green(kubernetesCluster.Name), kubernetesCluster.ID, "2m3s")
		}
	},
}
