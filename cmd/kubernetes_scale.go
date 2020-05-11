package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	_ "github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
	"time"
)

var KubernetesNewNodes int
var waitKubernetesNodes bool

var kubernetesScaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale a kubernetes cluster",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		kubernetesFindCluster, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("Unable to find a kubernetes cluster %s", err)
			os.Exit(1)
		}

		configKubernetes := &civogo.KubernetesClusterConfig{
			NumTargetNodes: KubernetesNewNodes,
		}

		kubernetesCluster, err := client.UpdateKubernetesCluster(kubernetesFindCluster.ID, configKubernetes)
		if err != nil {
			utility.Error("Unable to rename a kubernetes cluster %s", err)
			os.Exit(1)
		}

		if waitKubernetesNodes == true {

			stillScaling := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = "Scaling the kubernetes cluster... "
			s.Start()

			for stillScaling {
				kubernetesCheck, _ := client.FindKubernetesCluster(kubernetesCluster.ID)
				if kubernetesCheck.Status == "ACTIVE" {
					stillScaling = false
					s.Stop()
				}
				time.Sleep(5 * time.Second)
			}
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": kubernetesCluster.ID, "Name": kubernetesCluster.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("The kubernetes cluster %s was rescale from (%v) to (%v) nodes with ID %s\n", utility.Green(kubernetesCluster.Name), kubernetesFindCluster.NumTargetNode, kubernetesCluster.NumTargetNode, utility.Green(kubernetesCluster.ID))
		}
	},
}
