package cmd

import (
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var templateListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo template ls`,
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
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		templates, err := client.ListTemplates()
		if err != nil {
			utility.Error("%s", err)
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
