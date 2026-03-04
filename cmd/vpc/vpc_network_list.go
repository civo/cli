package vpc

import (
	"os"
	"strconv"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var vpcNetworkListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo vpc network ls -o custom -f "ID: Label"`,
	Short:   "List VPC networks",
	Long: `List all available VPC networks.
If you wish to use a custom format, the available fields are:

	* id
	* label
	* default
	* cidr
	* status
	* ipv4_enabled
	* ipv6_enabled
	* free_ip_count`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		networks, err := client.ListVPCNetworks()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		for _, network := range networks {
			ow.StartLine()
			ow.AppendDataWithLabel("id", network.ID, "ID")
			ow.AppendDataWithLabel("label", network.Label, "Label")
			ow.AppendDataWithLabel("default", strconv.FormatBool(network.Default), "Default")
			ow.AppendDataWithLabel("cidr", network.CIDR, "CIDR")
			ow.AppendDataWithLabel("status", network.Status, "Status")
			ow.AppendDataWithLabel("ipv4_enabled", strconv.FormatBool(network.IPv4Enabled), "IPv4 Enabled")
			ow.AppendDataWithLabel("ipv6_enabled", strconv.FormatBool(network.IPv6Enabled), "IPv6 Enabled")
			ow.AppendDataWithLabel("free_ip_count", strconv.Itoa(network.FreeIPCount), "Free IP Count")
		}

		ow.FinishAndPrintOutput()
	},
}
