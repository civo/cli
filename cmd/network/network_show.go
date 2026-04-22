package network

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
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
		ow.AppendDataWithLabel("vlan_id", fmt.Sprintf("%d", network.VlanID), "VLAN ID")
		ow.AppendDataWithLabel("physical_interface", network.PhysicalInterface, "Hardware Address")
		ow.AppendDataWithLabel("gateway_ipv4", network.GatewayIPv4, "Gateway IPv4")
		ow.AppendDataWithLabel("allocation_pool_v4_start", network.AllocationPoolV4Start, "Allocation Pool IPv4 Start")
		ow.AppendDataWithLabel("allocation_pool_v4_end", network.AllocationPoolV4End, "Allocation Pool IPv4 End")
		ow.AppendDataWithLabel("nameservers_v4", utility.SliceToString(network.NameserversV4), "Nameservers IPv4")
		ow.AppendDataWithLabel("nameservers_v6", utility.SliceToString(network.NameserversV6), "Nameservers IPv6")

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Println("Network Details:")
			fmt.Printf("ID: %s\n", network.ID)
			fmt.Printf("Name: %s\n", network.Name)
			fmt.Printf("Default: %s\n", utility.BoolToYesNo(network.Default))
			fmt.Printf("CIDR: %s\n", network.CIDR)
			fmt.Printf("Status: %s\n", network.Status)
			fmt.Printf("IPv4 Enabled: %s\n", utility.BoolToYesNo(network.IPv4Enabled))
			fmt.Printf("IPv6 Enabled: %s\n", utility.BoolToYesNo(network.IPv6Enabled))

			if network.VlanID != 0 {
				fmt.Println("\nVLAN Details:")
				fmt.Printf("VLAN ID: %d\n", network.VlanID)
				fmt.Printf("Hardware Address: %s\n", network.PhysicalInterface)
				fmt.Printf("Gateway IPv4: %s\n", network.GatewayIPv4)
				fmt.Printf("Allocation Pool IPv4 Start: %s\n", network.AllocationPoolV4Start)
				fmt.Printf("Allocation Pool IPv4 End: %s\n", network.AllocationPoolV4End)
			} else {
				fmt.Println("\nNo VLAN Configuration")
			}

			if len(network.NameserversV4) > 0 || len(network.NameserversV6) > 0 {
				fmt.Println("\nNameserver Details:")
				if len(network.NameserversV4) > 0 {
					fmt.Printf("Nameservers IPv4: %s\n", utility.SliceToString(network.NameserversV4))
				}
				if len(network.NameserversV6) > 0 {
					fmt.Printf("Nameservers IPv6: %s\n", utility.SliceToString(network.NameserversV6))
				}
			}
		}
	},
}
