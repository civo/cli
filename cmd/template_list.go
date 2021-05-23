package cmd

import (
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

type TemplateDisk struct {
	ID      string
	Name    string
	Version string
	Label   string
}

var templateListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo template ls`,
	Short:   "List templates",
	Long: `List all available templates.
If you wish to use a custom format, the available fields are:

	* id
	* name
	* version
	* label

Example: civo template ls -o custom -f "ID: Code (DefaultUsername)"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		templateDiskList := []TemplateDisk{}

		if client.Region == "SVG1" {
			templates, err := client.ListTemplates()
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}

			for _, v := range templates {
				templateDiskList = append(templateDiskList, TemplateDisk{ID: v.ID, Name: v.Name, Version: v.Code, Label: v.ShortDescription})
			}
		} else {
			diskImage, err := client.ListDiskImages()
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}

			for _, v := range diskImage {
				templateDiskList = append(templateDiskList, TemplateDisk{ID: v.ID, Name: v.Name, Version: v.Version, Label: v.Label})
			}

		}

		ow := utility.NewOutputWriter()

		for _, template := range templateDiskList {
			ow.StartLine()
			ow.AppendDataWithLabel("id", template.ID, "ID")
			ow.AppendDataWithLabel("name", template.Name, "Name")
			ow.AppendDataWithLabel("version", template.Version, "Version")
			ow.AppendDataWithLabel("label", template.Label, "Label")
		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
