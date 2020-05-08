package cmd

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var (
	shortDescriptionCreate                                                  string
	nameCreate, imageIDCreate, volumeIDCreate                               string
	codeCreate, descriptionCreate, defaultUsernameCreate, cloudConfigCreate string
)

var templateCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add", "new"},
	Short:   "Create a template",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
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
				fmt.Printf("Unable to read the cloud config file: %s\n", aurora.Red(err))
				os.Exit(1)
			}

			configTemplate.CloudConfig = string(data)
		}

		_, err = client.NewTemplate(configTemplate)
		if err != nil {
			fmt.Printf("Unable to create the template: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		template, err := client.GetTemplateByCode(codeCreate)
		if err != nil {
			fmt.Printf("Unable to find the template: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": template.ID, "Name": template.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created template with name %s with ID %s\n", aurora.Green(template.Name), aurora.Green(template.ID))
		}
	},
}
