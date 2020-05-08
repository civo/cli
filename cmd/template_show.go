package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
)

var templateShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show template",
	Long: `Show your current template.
If you wish to use a custom format, the available fields are:

	* ID
	* Code
	* Name
	* AccountID
	* ImageID
	* VolumeID
	* ShortDescription
	* Description
	* DefaultUsername
	* CloudConfig

Example: civo template show CODE -o custom -f "ID: Code (DefaultUsername)"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		template, err := client.GetTemplateByCode(args[0])
		if err != nil {
			fmt.Printf("Unable to search template: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		ow.AppendData("ID", template.ID)
		ow.AppendData("Code", template.Code)
		ow.AppendData("Name", template.Name)
		ow.AppendDataWithLabel("AccountID", template.AccountID, "Account ID")
		ow.AppendDataWithLabel("ImageID", template.ImageID, "Image ID")
		ow.AppendDataWithLabel("VolumeID", template.VolumeID, "Volume ID")
		ow.AppendDataWithLabel("ShortDescription", template.ShortDescription, "Short Description")
		ow.AppendData("Description", template.Description)
		ow.AppendDataWithLabel("DefaultUsername", template.DefaultUsername, "Default Username")

		if outputFormat == "json" || outputFormat == "custom" {
			ow.AppendData("CloudConfig", template.CloudConfig)
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		} else {
			ow.WriteKeyValues()

			if len(template.CloudConfig) > 0 {
				fmt.Println()
				ow.WriteSubheader("Cloud Config")
				fmt.Println(template.CloudConfig)
			}
		}

	},
}
