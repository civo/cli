package database

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var dbVersionListCmd = &cobra.Command{
	Use:     "versions",
	Aliases: []string{"version"},
	Example: `civo db versions`,
	Short:   "List all the available database versions",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		dbVersions, err := client.ListDBVersions()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		// Iterate through each database type and all its versions
		for dbName, versionDetails := range dbVersions {
			for _, version := range versionDetails {
				ow.StartLine()
				ow.AppendDataWithLabel("name", dbName, "Name")
				ow.AppendDataWithLabel("version", version.SoftwareVersion, "Version")
				ow.AppendDataWithLabel("default", fmt.Sprintf("%t", version.Default), "Default")
			}
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
