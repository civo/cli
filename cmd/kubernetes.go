package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var kubernetesCmd = &cobra.Command{
	Use:     "kubernetes",
	Aliases: []string{"kubernetes", "kube"},
	Short:   "Details of Civo kubernetes",
}

var kubernetesApplicationsCmd = &cobra.Command{
	Use:     "applications",
	Aliases: []string{"application", "app"},
	Short:   "Details of Civo kubernetes applications",
}

func init() {
	rootCmd.AddCommand(kubernetesCmd)
	kubernetesCmd.AddCommand(kubernetesListCmd)
	kubernetesCmd.AddCommand(kubernetesListVersionCmd)
	kubernetesCmd.AddCommand(kubernetesShowCmd)
	kubernetesCmd.AddCommand(kubernetesConfigCmd)
	kubernetesCmd.AddCommand(kubernetesCreateCmd)
	kubernetesCmd.AddCommand(kubernetesRenameCmd)
	kubernetesCmd.AddCommand(kubernetesUpgradeCmd) // TODO: Check this with @andy
	kubernetesCmd.AddCommand(kubernetesScaleCmd)
	kubernetesCmd.AddCommand(kubernetesRemoveCmd)

	/*
		Flags for kubernetes config
	*/
	home, _ := os.UserHomeDir()
	kubernetesConfigCmd.Flags().BoolVarP(&saveConfig, "save", "s", false, "save the config")
	kubernetesConfigCmd.Flags().BoolVarP(&mergeConfig, "merge", "m", false, "merge the config with existing kubeconfig if it already exists.")
	kubernetesConfigCmd.Flags().StringVarP(&localPathConfig, "local-path", "p", fmt.Sprintf("%s/.kube/config", home), "local path to save the kubeconfig file")

	/*
		Flags for kubernetes create
	*/
	kubernetesCreateCmd.Flags().StringVarP(&TargetNodesSize, "size", "s", "g2.medium", "the size of nodes to create.")
	kubernetesCreateCmd.Flags().IntVarP(&NumTargetNodes, "nodes", "n", 3, "the number of nodes to create (the master also acts as a node).")
	kubernetesCreateCmd.Flags().StringVarP(&KubernetesVersion, "version", "v", "latest", "the k3s version to use on the cluster. Defaults to the latest.")
	kubernetesCreateCmd.Flags().BoolVarP(&waitKubernetes, "wait", "w", false, "a simple flag (e.g. --wait) that will cause the CLI to spin and wait for the cluster to be ACTIVE")

	/*
		Flags for kubernetes rename
	*/
	kubernetesRenameCmd.Flags().StringVarP(&KubernetesNewName, "name", "n", "", "the new name for the cluster.")

	/*
		Flags for kubernetes upgrade
	*/
	kubernetesUpgradeCmd.Flags().StringVarP(&KubernetesNewVersion, "version", "v", "", "change the version of the cluster.")
	/*
		Flags for kubernetes scale
	*/
	kubernetesScaleCmd.Flags().IntVarP(&KubernetesNewNodes, "nodes", "n", 3, "change the total nodes of the cluster.")
	kubernetesScaleCmd.Flags().BoolVarP(&waitKubernetesNodes, "wait", "w", false, "a simple flag (e.g. --wait) that will cause the CLI to spin and wait for the cluster to be ACTIVE")

	/*
		Kube application
	*/
	kubernetesCmd.AddCommand(kubernetesApplicationsCmd)
	kubernetesApplicationsCmd.AddCommand(kubernetesAppListCmd)
	// TODO: show command
	kubernetesApplicationsCmd.AddCommand(kubernetesAppAddCmd)
	kubernetesAppAddCmd.Flags().StringVarP(&KubernetesClusterApp, "cluster", "c", "", "the name of the cluster to install the app.")

}
