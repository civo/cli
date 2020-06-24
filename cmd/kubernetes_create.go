package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var numTargetNodes int
var waitKubernetes bool
var kubernetesVersion, targetNodesSize, clusterName string

var kubernetesCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo kubernetes create CLUSTER_NAME [flags]",
	Short:   "Create a new Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if len(args) > 0 {
			clusterName = args[0]
		} else {
			clusterName = utility.RandomName()
		}

		configKubernetes := &civogo.KubernetesClusterConfig{
			Name:            clusterName,
			NumTargetNodes:  numTargetNodes,
			TargetNodesSize: targetNodesSize,
		}

		if kubernetesVersion != "latest" {
			configKubernetes.KubernetesVersion = kubernetesVersion
		}

		kubernetesCluster, err := client.NewKubernetesClusters(configKubernetes)
		if err != nil {
			utility.Error("Creating a Kubernetes cluster failed with %s", err)
			os.Exit(1)
		}

		var executionTime string

		if waitKubernetes == true {
			startTime := utility.StartTime()

			stillCreating := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = fmt.Sprintf("Creating a %s node k3s cluster of %s instances called %s... ", strconv.Itoa(kubernetesCluster.NumTargetNode), kubernetesCluster.TargetNodeSize, kubernetesCluster.Name)
			s.Start()

			for stillCreating {
				kubernetesCheck, err := client.FindKubernetesCluster(kubernetesCluster.ID)
				if err != nil {
					utility.Error("Finding the Kubernetes cluster failed with %s", err)
					os.Exit(1)
				}
				if kubernetesCheck.Status == "ACTIVE" {
					stillCreating = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}

			executionTime = utility.TrackTime(startTime)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": kubernetesCluster.ID, "Name": kubernetesCluster.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			if executionTime != "" {
				fmt.Printf("The cluster %s (%s) has been created in %s\n", utility.Green(kubernetesCluster.Name), kubernetesCluster.ID, executionTime)
			} else {
				fmt.Printf("The cluster %s (%s) has been created\n", utility.Green(kubernetesCluster.Name), kubernetesCluster.ID)
			}

		}
	},
}
