package cmd

import (
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
)

var templateListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List templates",
	Long: `List all available templates.
If you wish to use a custom format, the available fields are:

	* ID
	* Code
	* Name
	* ShortDescription
	* Description
	* DefaultUsername

Example: civo template ls -o custom -f "ID: Code (DefaultUsername)"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		templates, err := client.ListTemplates()
		if err != nil {
			utility.Error("Unable to list templates %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		for _, template := range templates {
			ow.StartLine()
			ow.AppendData("ID", template.ID)
			ow.AppendData("Code", template.Code)
			ow.AppendData("Name", template.Name)
			ow.AppendDataWithLabel("ImageID", template.ImageID, "Image ID")
			ow.AppendDataWithLabel("ShortDescription", template.ShortDescription, "Short Description")
			ow.AppendData("Description", template.Description)
			ow.AppendDataWithLabel("DefaultUsername", template.DefaultUsername, "Default Username")
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
