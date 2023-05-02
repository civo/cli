// Package kfcluster is the root command for Civo KFCluster
package kfcluster

import (
	"errors"

	"github.com/spf13/cobra"
)

// KFClusterCmd is the root command for the kfcluster subcommand
var KFClusterCmd = &cobra.Command{
	Use:     "kfcluster",
	Aliases: []string{"kfclusters", "kf", "kfc", "kfcs", "kubeflow", "kfaas"},
	Short:   "Manage Civo Kubeflow Clusters",
	Long:    `Create, update, delete, and list Civo Kubeflow Clusters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	KFClusterCmd.AddCommand(kfcListCmd)
	KFClusterCmd.AddCommand(kfcCreateCmd)
	KFClusterCmd.AddCommand(kfcUpdateCmd)
	KFClusterCmd.AddCommand(kfcDeleteCmd)
	KFClusterCmd.AddCommand(kfcSizeCmd)

	kfcCreateCmd.Flags().StringVarP(&firewallID, "firewall", "", "", "the firewall to use for the kubeflow cluster")
	kfcCreateCmd.Flags().StringVarP(&networkID, "network", "n", "", "the network to use for the kubeflow cluster")
	kfcCreateCmd.Flags().StringVarP(&size, "size", "s", "g3.kf.small", "the size of the kubeflow cluster.")
	kfcCreateCmd.MarkFlagRequired("network")

	kfcUpdateCmd.Flags().StringVarP(&updatedName, "name", "n", "", "the new name for the kubeflow cluster")
}
