package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
	_ "strconv"
	"time"
)

var NumTargetNodes int
var waitKubernetes bool
var (
	KubernetesVersion string
	TargetNodesSize   string
)

var kubernetesCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new kubernetes cluster",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		isValidName := false

		_, err = client.FindKubernetesCluster(args[0])
		if err != nil {
			isValidName = true
		}

		if !isValidName {
			fmt.Printf("The %s is nos valida name for the cluster\n", utility.Red(args[0]))
			os.Exit(1)
		}

		configKubernetes := &civogo.KubernetesClusterConfig{
			Name:              args[0],
			NumTargetNodes:    NumTargetNodes,
			TargetNodesSize:   TargetNodesSize,
			KubernetesVersion: KubernetesVersion,
		}

		kubernetesCluster, err := client.NewKubernetesClusters(configKubernetes)
		if err != nil {
			utility.Error("Unable to create a kubernetes cluster %s", err)
			os.Exit(1)
		}

		if waitKubernetes == true {

			stillCreating := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = "Creating kubernetes cluster... "
			s.Start()

			for stillCreating {
				kubernetesCheck, _ := client.FindKubernetesCluster(kubernetesCluster.ID)
				if kubernetesCheck.Status == "ACTIVE" {
					stillCreating = false
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
			fmt.Printf("Created a new kubernetes cluster with Name %s with ID %s\n", utility.Green(kubernetesCluster.Name), utility.Green(kubernetesCluster.ID))
		}
	},
}
