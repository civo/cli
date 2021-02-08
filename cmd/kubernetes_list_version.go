package cmd

import (
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var kubernetesListVersionCmd = &cobra.Command{
	Use:     "versions",
	Aliases: []string{"version"},
	Example: `civo kubernetes versions`,
	Short:   "List all Kubernetes clusters version",
	Long: `List all Kubernetes clusters versions.
If you wish to use a custom format, the available fields are:

	* Version
	* Type
	* Default`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		kubeVersions, err := client.ListAvailableKubernetesVersions()
		if err != nil {
			utility.Error("%s", err)
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
