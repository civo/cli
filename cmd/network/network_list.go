package network

import (
	"strconv"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var networkListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo network ls -o custom -f "ID: Name (CIDR)"`,
	Short:   "List networks",
	Long: `List all available networks.
If you wish to use a custom format, the available fields are:

	* id
	* label
	* region
	* default
	* status`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			return
		}

		networks, err := client.ListNetworks()
		if err != nil {
			utility.Error("%s", err)
			return
		}

		ow := utility.NewOutputWriter()

		for _, network := range networks {
			ow.StartLine()
			ow.AppendDataWithLabel("id", network.ID, "ID")
			ow.AppendDataWithLabel("label", network.Label, "Label")
			ow.AppendDataWithLabel("region", client.Region, "Region")
			ow.AppendDataWithLabel("ipv4_enabled", strconv.FormatBool(network.IPv4Enabled), "IPv4 Enabled")
			ow.AppendDataWithLabel("ipv6_enabled", strconv.FormatBool(network.IPv6Enabled), "IPv6 Enabled")
			ow.AppendDataWithLabel("default", strconv.FormatBool(network.Default), "Default")
			ow.AppendDataWithLabel("status", network.Status, "Status")
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
