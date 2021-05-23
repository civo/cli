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
	shortDescriptionCreate                                                  string
	nameCreate, imageIDCreate, volumeIDCreate                               string
	codeCreate, descriptionCreate, defaultUsernameCreate, cloudConfigCreate string
)

var templateCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add", "new"},
	Example: "civo template create [flags]",
	Short:   "Create a template",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		configTemplate := &civogo.Template{}

		if nameCreate != "" {
			configTemplate.Name = nameCreate
		}

		if codeCreate != "" {
			configTemplate.Code = codeCreate
		}

		if imageIDCreate != "" {
			configTemplate.ImageID = imageIDCreate
		}

		if volumeIDCreate != "" {
			configTemplate.VolumeID = volumeIDCreate
		}

		if imageIDCreate != "" {
			configTemplate.ImageID = imageIDCreate
		}

		if shortDescriptionCreate != "" {
			configTemplate.ShortDescription = shortDescriptionCreate
		}

		if descriptionCreate != "" {
			configTemplate.Description = descriptionCreate
		}

		if defaultUsernameCreate != "" {
			configTemplate.DefaultUsername = defaultUsernameCreate
		}

		if cloudConfigCreate != "" {
			// reading the file
			data, err := ioutil.ReadFile(cloudConfigCreate)
			if err != nil {
				utility.Error("Reading the cloud config file failed with %s", err)
				os.Exit(1)
			}

			configTemplate.CloudConfig = string(data)
		}

		_, err = client.NewTemplate(configTemplate)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		template, err := client.GetTemplateByCode(codeCreate)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": template.ID, "name": template.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created the template with name %s with ID %s\n", utility.Green(template.Name), utility.Green(template.ID))
		}
	},
}
