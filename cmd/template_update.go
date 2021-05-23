package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var (
	shortDescription                                string
	name, description, defaultUsername, cloudConfig string
)

var templateUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"change", "modify"},
	Example: "civo template update TEMPLATE_CODE [flags]",
	Short:   "Update a template",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		template, err := client.GetTemplateByCode(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		configTemplateUpdate := &civogo.Template{
			ImageID:  template.ImageID,
			VolumeID: template.VolumeID,
		}

		if name != "" {
			configTemplateUpdate.Name = name
		} else {
			configTemplateUpdate.Name = template.Name
		}

		if shortDescription != "" {
			configTemplateUpdate.ShortDescription = shortDescription
		} else {
			configTemplateUpdate.ShortDescription = template.ShortDescription
		}

		if description != "" {
			configTemplateUpdate.Description = description
		} else {
			configTemplateUpdate.Description = template.Description
		}

		if defaultUsername != "" {
			configTemplateUpdate.DefaultUsername = defaultUsername
		} else {
			configTemplateUpdate.DefaultUsername = template.DefaultUsername
		}

		if cloudConfig != "" {
			data, err := ioutil.ReadFile(cloudConfig)
			if err != nil {
				utility.Error("Reading the cloud config file failed with %s", err)
				os.Exit(1)
			}

			configTemplateUpdate.CloudConfig = string(data)
		} else {
			configTemplateUpdate.CloudConfig = template.CloudConfig
		}

		templateUpdate, err := client.UpdateTemplate(template.ID, configTemplateUpdate)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": templateUpdate.ID, "name": templateUpdate.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Updated template with name %s with ID %s\n", utility.Green(templateUpdate.Name), utility.Green(templateUpdate.ID))
		}
	},
}
