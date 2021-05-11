package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesNode string

var kubernetesRecycleCmd = &cobra.Command{
	Use:     "recycle",
	Short:   "recycle a Kubernetes node",
	Example: "civo kubernetes recycle CLUSTER_NAME --node NODE_NAME [flags]",
	Args:    cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return getAllKubernetesClusterName(), cobra.ShellCompDirectiveNoFileComp
		}
		return getKubernetesClusterName(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
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

		kubernetesFindCluster, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		_, err = client.RecycleKubernetesCluster(kubernetesFindCluster.ID, kubernetesNode)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": kubernetesFindCluster.ID, "Name": kubernetesFindCluster.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("The node (%s) was recycled\n", utility.Green(kubernetesNode))
		}
	},
}

func getKubernetesClusterName(value string) []string {

	client, err := config.CivoAPIClient()
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	k8s, err := client.FindKubernetesCluster(value)
	if err != nil {
		utility.Error("Unable to get kubernetes cluster %s", err)
		os.Exit(1)
	}

	var k8sList []string
	k8sList = append(k8sList, k8s.Name)

	return k8sList

}

func getAllKubernetesClusterName() []string {

	client, err := config.CivoAPIClient()
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	k8s, err := client.ListKubernetesClusters()
	if err != nil {
		utility.Error("Unable to list kubernetes cluster %s", err)
		os.Exit(1)
	}

	var k8sList []string
	for _, v := range k8s.Items {
		k8sList = append(k8sList, v.Name)
	}

	return k8sList
}
