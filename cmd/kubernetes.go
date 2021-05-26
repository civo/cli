package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	kubemartutils "github.com/kubemart/kubemart-cli/pkg/utils"
	"github.com/spf13/cobra"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var debug bool
var kubernetesClusterApp string

var legacyMarketplaceWarning = "This command is only available for the new version of marketplace. Your current cluster is running on legacy marketplace.\nYou can launch a new cluster to start using the new marketplace."

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
	kubernetesApplicationsCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "print verbose logs when running command")
	cobra.OnInitialize(setLogLevelEnvIfFlagIsTrue)
	kubernetesCmd.AddCommand(kubernetesApplicationsCmd)
	kubernetesApplicationsCmd.AddCommand(kubernetesAppListCmd)
	kubernetesApplicationsCmd.AddCommand(kubernetesAppAddCmd)
	kubernetesApplicationsCmd.AddCommand(kubernetesAppShowCmd)
	kubernetesApplicationsCmd.AddCommand(kubernetesAppInstalledCmd)
	kubernetesApplicationsCmd.AddCommand(kubernetesAppUninstallCmd)
	kubernetesApplicationsCmd.AddCommand(kubernetesAppUpdateCmd)
	kubernetesAppAddCmd.Flags().StringVarP(&kubernetesClusterApp, "cluster", "c", "", "the name of the cluster")
	kubernetesAppInstalledCmd.Flags().StringVarP(&kubernetesClusterApp, "cluster", "c", "", "the name of the cluster")
	kubernetesAppUninstallCmd.Flags().StringVarP(&kubernetesClusterApp, "cluster", "c", "", "the name of the cluster")
	kubernetesAppUpdateCmd.Flags().StringVarP(&kubernetesClusterApp, "cluster", "c", "", "the name of the cluster")
	kubernetesAppListCmd.Flags().StringVarP(&kubernetesClusterApp, "cluster", "c", "", "the name of the cluster")

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

// setLogLevelEnvIfFlagIsTrue will set LOGLEVEL env variable
// when user use '--debug' or '-d' flag
func setLogLevelEnvIfFlagIsTrue() {
	if debug {
		os.Setenv("LOGLEVEL", "debug")
	}
}

// isKubemartCluster will return true if the cluster has
// Kubemart ConfigMap (inside "kubemart-system" namespace)
func isKubemartCluster(clusterIdentifier string) (bool, error) {
	namespace := "kubemart-system"
	if clusterIdentifier == "" {
		return false, fmt.Errorf("cluster is not set")
	}

	kubemartutils.DebugPrintf("Creating Civo API client\n")
	client, err := config.CivoAPIClient()
	if regionSet != "" {
		client.Region = regionSet
	}
	if err != nil {
		return false, err
	}

	kubemartutils.DebugPrintf("Finding Civo Kubernetes cluster\n")
	kubernetesCluster, err := client.FindKubernetesCluster(clusterIdentifier)
	if err != nil {
		return false, err
	}

	kubemartutils.DebugPrintf("Finding kubeconfig\n")
	kubeconfigStr := kubernetesCluster.KubeConfig
	kcBytes := []byte(kubeconfigStr)

	kubemartutils.DebugPrintf("Getting REST config from kubeconfig\n")
	rc, err := clientcmd.RESTConfigFromKubeConfig(kcBytes)
	if err != nil {
		return false, err
	}

	kubemartutils.DebugPrintf("Creating Kubernetes clientset from REST config\n")
	cs, err := kubernetes.NewForConfig(rc)
	if err != nil {
		return false, err
	}

	kubemartutils.DebugPrintf("Creating ConfigMap client from Kubernetes clientset\n")
	cmClient := cs.CoreV1().ConfigMaps(namespace)

	kubemartutils.DebugPrintf("Checking Kubemart ConfigMap\n")
	_, err = cmClient.Get(context.Background(), "kubemart-config", metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			kubemartutils.DebugPrintf("Kubemart ConfigMap not found\n")
			return false, nil
		}
		kubemartutils.DebugPrintf("Error occured when checking Kubemart ConfigMap\n")
		return false, err
	}

	kubemartutils.DebugPrintf("Kubemart ConfigMap found\n")
	return true, nil
}

func syncKubemartApps() error {
	err := kubemartutils.CloneAppFilesIfNotExist()
	if err != nil {
		return err
	}

	_, err = kubemartutils.UpdateAppsCacheIfStale()
	if err != nil {
		return err
	}

	return nil
}
