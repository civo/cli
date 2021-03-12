package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesNewNodes int
var waitKubernetesNodes bool

var kubernetesScaleCmd = &cobra.Command{
	Use:     "scale",
	Short:   "Scale a Kubernetes cluster",
	Example: "civo kubernetes scale CLUSTER_NAME --nodes 4 [flags]",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		kubernetesFindCluster, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		configKubernetes := &civogo.KubernetesClusterConfig{
			NumTargetNodes: kubernetesNewNodes,
		}

		kubernetesCluster, err := client.UpdateKubernetesCluster(kubernetesFindCluster.ID, configKubernetes)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if waitKubernetesNodes {

			stillScaling := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = "Scaling the Kubernetes cluster... "
			s.Start()

			for stillScaling {
				kubernetesCheck, err := client.FindKubernetesCluster(kubernetesCluster.ID)
				if err != nil {
					utility.Error("Finding the kubernetes cluster failed with %s", err)
					os.Exit(1)
				}
				if kubernetesCheck.Status == "ACTIVE" {
					stillScaling = false
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
			fmt.Printf("The Kubernetes cluster %s (with ID %s) was rescaled from %v to %v nodes\n", utility.Green(kubernetesCluster.Name), utility.Green(kubernetesCluster.ID), kubernetesFindCluster.NumTargetNode, kubernetesCluster.NumTargetNode)
		}
	},
}
