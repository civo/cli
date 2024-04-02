package network

import (
	"fmt"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
)

var networkShowCmd = &cobra.Command{
	Use:     "show [NETWORK-NAME/NETWORK-ID]",
	Short:   "Show details of a specific Civo network, including VLAN information if available",
	Aliases: []string{"get", "describe", "inspect"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		networkID := args[0]

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("failed to create Civo API client: %s", err)
			os.Exit(1)
		}

		network, err := client.GetNetwork(networkID)
		if err != nil {
			utility.Error("Failed to retrieve network: %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", network.ID, "ID")
		ow.AppendDataWithLabel("name", network.Name, "Name")
		ow.AppendDataWithLabel("default", utility.BoolToYesNo(network.Default), "Default")
		ow.AppendDataWithLabel("cidr", network.CIDR, "CIDR")
		ow.AppendDataWithLabel("status", network.Status, "Status")
		ow.AppendDataWithLabel("ipv4_enabled", utility.BoolToYesNo(network.IPv4Enabled), "IPv4 Enabled")
		ow.AppendDataWithLabel("ipv6_enabled", utility.BoolToYesNo(network.IPv6Enabled), "IPv6 Enabled")
		if len(network.NameserversV4) > 0 {
			ow.AppendDataWithLabel("nameservers_v4", utility.SliceToString(network.NameserversV4), "Nameservers IPv4")
		}
		if len(network.NameserversV6) > 0 {
			ow.AppendDataWithLabel("nameservers_v6", utility.SliceToString(network.NameserversV6), "Nameservers IPv6")
		}
		// Add VLAN details if available
		if network.VLAN.VlanID != 0 {
			ow.AppendDataWithLabel("vlan_id", fmt.Sprintf("%d", network.VLAN.VlanID), "VLAN ID")
			ow.AppendDataWithLabel("hardware_addr", network.VLAN.HardwareAddr, "Hardware Address")
			ow.AppendDataWithLabel("cidr_v4", network.VLAN.CIDRv4, "VLAN CIDRv4")
			ow.AppendDataWithLabel("gateway_ipv4", network.VLAN.GatewayIPv4, "Gateway IPv4")
			ow.AppendDataWithLabel("allocation_pool_v4_start", network.VLAN.AllocationPoolV4Start, "Allocation Pool Start IPv4")
			ow.AppendDataWithLabel("allocation_pool_v4_end", network.VLAN.AllocationPoolV4End, "Allocation Pool End IPv4")
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
