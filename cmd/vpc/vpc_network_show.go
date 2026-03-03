package vpc

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var vpcNetworkShowCmd = &cobra.Command{
	Use:     "show [NETWORK-NAME/NETWORK-ID]",
	Short:   "Show details of a specific VPC network",
	Aliases: []string{"get", "describe", "inspect"},
	Args:    cobra.ExactArgs(1),
	Example: "civo vpc network show NETWORK_NAME",
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

		network, err := client.FindVPCNetwork(args[0])
		if err != nil {
			utility.Error("Network %s", err)
			os.Exit(1)
		}

		network, err = client.GetVPCNetwork(network.ID)
		if err != nil {
			utility.Error("Failed to retrieve network: %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", network.ID, "ID")
		ow.AppendDataWithLabel("label", network.Label, "Label")
		ow.AppendDataWithLabel("default", utility.BoolToYesNo(network.Default), "Default")
		ow.AppendDataWithLabel("cidr", network.CIDR, "CIDR")
		ow.AppendDataWithLabel("status", network.Status, "Status")
		ow.AppendDataWithLabel("ipv4_enabled", utility.BoolToYesNo(network.IPv4Enabled), "IPv4 Enabled")
		ow.AppendDataWithLabel("ipv6_enabled", utility.BoolToYesNo(network.IPv6Enabled), "IPv6 Enabled")

		switch common.OutputFormat {
		case "json":
			ow.ToJSON(network, common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Println("VPC Network Details:")
			fmt.Printf("ID: %s\n", network.ID)
			fmt.Printf("Name: %s\n", network.Label)
			fmt.Printf("Default: %s\n", utility.BoolToYesNo(network.Default))
			fmt.Printf("CIDR: %s\n", network.CIDR)
			fmt.Printf("Status: %s\n", network.Status)
			fmt.Printf("IPv4 Enabled: %s\n", utility.BoolToYesNo(network.IPv4Enabled))
			fmt.Printf("IPv6 Enabled: %s\n", utility.BoolToYesNo(network.IPv6Enabled))

			if len(network.NameserversV4) > 0 {
				fmt.Printf("Nameservers IPv4: %s\n", utility.SliceToString(network.NameserversV4))
			}
			if len(network.NameserversV6) > 0 {
				fmt.Printf("Nameservers IPv6: %s\n", utility.SliceToString(network.NameserversV6))
			}
		}
	},
}
