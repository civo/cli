package cmd

import (
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var kubernetesListVersionCmd = &cobra.Command{
	Use:     "versions",
	Aliases: []string{"version"},
	Short:   "List all kubernetes clusters version",
	Long: `List all kubernetes clusters versions.
If you wish to use a custom format, the available fields are:

	* Version
	* Type
	* Default

Example: civo kubernetes versions -o custom -f "Version: Default"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		kubeVersions, err := client.ListAvailableKubernetesVersions()
		if err != nil {
			utility.Error("Unable to list kubernetes cluster %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, version := range kubeVersions {
			ow.StartLine()

			ow.AppendData("Version", version.Version)
			ow.AppendData("Type", version.Type)
			ow.AppendData("Default", strconv.FormatBool(version.Default))
		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
