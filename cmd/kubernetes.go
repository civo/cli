package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesCmd = &cobra.Command{
	Use:     "kubernetes",
	Aliases: []string{"k3s", "k8s", "kube"},
	Short:   "Details of Civo Kubernetes clusters",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

var kubernetesNodePoolCmd = &cobra.Command{
	Use:     "node-pool",
	Aliases: []string{"pool", "node-pool"},
	Short:   "Details of Civo Kubernetes applications",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

var kubernetesApplicationsCmd = &cobra.Command{
	Use:     "applications",
	Aliases: []string{"application", "app", "apps", "app", "application", "addon", "addons", "marketplace", "k8s-apps", "k8s-app", "k3s-apps", "k3s-app"},
	Short:   "Details of Civo Kubernetes applications",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	rootCmd.AddCommand(kubernetesCmd)
	kubernetesCmd.AddCommand(kubernetesListCmd)
	kubernetesCmd.AddCommand(kubernetesSizeCmd)
	kubernetesCmd.AddCommand(kubernetesListVersionCmd)
	kubernetesCmd.AddCommand(kubernetesShowCmd)
	kubernetesCmd.AddCommand(kubernetesConfigCmd)
	kubernetesCmd.AddCommand(kubernetesCreateCmd)
	kubernetesCmd.AddCommand(kubernetesRenameCmd)
	kubernetesCmd.AddCommand(kubernetesUpgradeCmd)
	kubernetesCmd.AddCommand(kubernetesRemoveCmd)
	kubernetesCmd.AddCommand(kubernetesRecycleCmd)

	home, err := os.UserHomeDir()
	if err != nil {
		utility.Error("%s", err)
		os.Exit(1)
	}

	kubernetesConfigCmd.Flags().BoolVarP(&saveConfig, "save", "s", false, "save the config")
	kubernetesConfigCmd.Flags().BoolVarP(&switchConfig, "switch", "i", false, "switch context to newly-created cluster")
	kubernetesConfigCmd.Flags().BoolVarP(&mergeConfig, "merge", "m", false, "merge the config with existing kubeconfig if it already exists.")
	kubernetesConfigCmd.Flags().StringVarP(&localPathConfig, "local-path", "p", fmt.Sprintf("%s/.kube/config", home), "local path to save the kubeconfig file")

	kubernetesCreateCmd.Flags().StringVarP(&targetNodesSize, "size", "s", "g3.k3s.medium", "the size of nodes to create.")
	kubernetesCreateCmd.Flags().StringVarP(&networkID, "network", "t", "default", "the name of the network to use in the creation")
	kubernetesCreateCmd.Flags().IntVarP(&numTargetNodes, "nodes", "n", 3, "the number of nodes to create (the master also acts as a node).")
	kubernetesCreateCmd.Flags().StringVarP(&kubernetesVersion, "version", "v", "latest", "the k3s version to use on the cluster. Defaults to the latest.")
	kubernetesCreateCmd.Flags().StringVarP(&applications, "applications", "a", "", "optional, use names shown by running 'civo kubernetes applications ls'")
	kubernetesCreateCmd.Flags().StringVarP(&removeapplications, "remove-applications", "r", "", "optional, remove default application names shown by running  'civo kubernetes applications ls'")
	kubernetesCreateCmd.Flags().BoolVarP(&waitKubernetes, "wait", "w", false, "a simple flag (e.g. --wait) that will cause the CLI to spin and wait for the cluster to be ACTIVE")
	kubernetesCreateCmd.Flags().BoolVarP(&saveConfigKubernetes, "save", "", false, "save the config")
	kubernetesCreateCmd.Flags().BoolVarP(&mergeConfigKubernetes, "merge", "m", false, "merge the config with existing kubeconfig if it already exists.")
	kubernetesCreateCmd.Flags().BoolVarP(&switchConfigKubernetes, "switch", "", false, "switch context to newly-created cluster")

	kubernetesRenameCmd.Flags().StringVarP(&kubernetesNewName, "name", "n", "", "the new name for the cluster.")

	kubernetesUpgradeCmd.Flags().StringVarP(&kubernetesNewVersion, "version", "v", "", "change the version of the cluster.")
	kubernetesUpgradeCmd.MarkFlagRequired("version")

	kubernetesRecycleCmd.Flags().StringVarP(&kubernetesNode, "node", "n", "", "the node that needs to be recycled.")
	kubernetesRecycleCmd.MarkFlagRequired("node")

	// Kubernetes Applications
	kubernetesCmd.AddCommand(kubernetesApplicationsCmd)
	kubernetesApplicationsCmd.AddCommand(kubernetesAppListCmd)
	kubernetesApplicationsCmd.AddCommand(kubernetesAppAddCmd)
	kubernetesApplicationsCmd.AddCommand(kubernetesAppShowCmd)

	kubernetesAppAddCmd.Flags().StringVarP(&kubernetesClusterApp, "cluster", "c", "", "the name of the cluster to install the app.")

	// Kubernetes NodePool
	kubernetesCmd.AddCommand(kubernetesNodePoolCmd)
	kubernetesNodePoolCmd.AddCommand(kubernetesNodePoolCreateCmd)
	kubernetesNodePoolCreateCmd.Flags().StringVarP(&targetNodesPoolSize, "size", "s", "g3.k3s.medium", "the size of nodes to create.")
	kubernetesNodePoolCreateCmd.Flags().IntVarP(&numTargetNodesPool, "nodes", "n", 3, "the number of nodes to create for the pool.")

	kubernetesNodePoolCmd.AddCommand(kubernetesNodePoolDeleteCmd)
	kubernetesNodePoolCmd.AddCommand(kubernetesNodePoolScaleCmd)
	kubernetesNodePoolScaleCmd.Flags().IntVarP(&numTargetNodesPoolScale, "nodes", "n", 3, "the number of nodes to scale for the pool.")

}

func getKubernetesList(value string) []string {
	client, err := config.CivoAPIClient()
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	cluster, err := client.FindKubernetesCluster(value)
	if err != nil {
		utility.Error("Unable to list domains %s", err)
		os.Exit(1)
	}

	var clusterList []string
	clusterList = append(clusterList, cluster.Name)

	return clusterList

}

func getAllKubernetesList() []string {
	client, err := config.CivoAPIClient()
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	cluster, err := client.ListKubernetesClusters()
	if err != nil {
		utility.Error("Unable to list kubernetes cluster %s", err)
		os.Exit(1)
	}

	var clusterList []string
	for _, v := range cluster.Items {
		clusterList = append(clusterList, v.Name)
	}

	return clusterList

}
