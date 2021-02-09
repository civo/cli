package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var numTargetNodes int
var waitKubernetes, saveConfigKubernetes, switchConfigKubernetes bool
var kubernetesVersion, targetNodesSize, clusterName, applications, removeapplications, installApplications, networkID string

var kubernetesCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo kubernetes create CLUSTER_NAME [flags]",
	Short:   "Create a new Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {

		check, region, err := utility.CheckAvailability("kubernetes", regionSet)
		if err != nil {
			utility.Error("Error checking availability %s", err)
			os.Exit(1)
		}

		if check == false {
			utility.Error("Sorry you can't create a kubernetes cluster in the %s region", region)
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}

		// fix to the kubernetes because in other region
		// we use g3.medium
		if client.Region != "SVG1" {
			targetNodesSize = "g3.medium"
		}

		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if !waitKubernetes {
			if saveConfigKubernetes || switchConfigKubernetes {
				utility.Error("you can't use --save or --switch without --wait")
				os.Exit(1)
			}
		} else {
			if switchConfigKubernetes && !saveConfigKubernetes {
				utility.Error("you can't use --switch without --save")
				os.Exit(1)
			}
		}

		if len(args) > 0 {
			clusterName = args[0]
		} else {
			clusterName = utility.RandomName()
		}

		if networkID == "default" {
			network, err := client.GetDefaultNetwork()
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}

			networkID = network.ID

		} else {
			network, err := client.FindNetwork(networkID)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}

			networkID = network.ID
		}

		configKubernetes := &civogo.KubernetesClusterConfig{
			Name:            clusterName,
			NumTargetNodes:  numTargetNodes,
			TargetNodesSize: targetNodesSize,
			NetworkID:       networkID,
		}

		if kubernetesVersion != "latest" {
			configKubernetes.KubernetesVersion = kubernetesVersion
		}

		if applications != "" {
			installApplications = applications
		}

		if removeapplications != "" {
			var rmApp []string
			for _, v := range strings.Split(removeapplications, ",") {
				rmApp = append(rmApp, fmt.Sprintf("-%s", v))
			}
			if installApplications != "" {
				installApplications = fmt.Sprintf("%s,%s", installApplications, strings.Join(rmApp, ","))
			} else {
				installApplications = fmt.Sprintf("%s", strings.Join(rmApp, ","))
			}

		}

		if installApplications != "" {
			configKubernetes.Applications = installApplications
		}

		kubernetesCluster, err := client.NewKubernetesClusters(configKubernetes)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		var executionTime string

		if waitKubernetes {
			startTime := utility.StartTime()

			stillCreating := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = fmt.Sprintf("Creating a %s node k3s cluster of %s instances called %s... ", strconv.Itoa(kubernetesCluster.NumTargetNode), kubernetesCluster.TargetNodeSize, kubernetesCluster.Name)
			s.Start()

			for stillCreating {
				kubernetesCheck, err := client.FindKubernetesCluster(kubernetesCluster.ID)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				if kubernetesCheck.Ready {
					stillCreating = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}

			executionTime = utility.TrackTime(startTime)
		}

		if saveConfigKubernetes {
			kube, err := client.FindKubernetesCluster(kubernetesCluster.ID)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}

			err = utility.ObtainKubeConfig(localPathConfig, kube.KubeConfig, true, switchConfigKubernetes, kube.Name)
			if err != nil {
				utility.Error("Saving the cluster config failed with %s", err)
				os.Exit(1)
			}
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
