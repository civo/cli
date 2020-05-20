package cmd

import "github.com/spf13/cobra"

var templateCmd = &cobra.Command{
	Use:     "template",
	Aliases: []string{"templates"},
	Short:   "Details of Civo templates",
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
	templateCmd.AddCommand(templateUpdateCmd)
	templateCmd.AddCommand(templateCreateCmd)
	templateCmd.AddCommand(templateRemoveCmd)

	templateCreateCmd.Flags().StringVarP(&codeCreate, "code", "c", "", "The code name of the template, this can't change after creation")
	templateCreateCmd.MarkFlagRequired("code")
	templateCreateCmd.Flags().StringVarP(&imageIDCreate, "image-id", "m", "", "The image id for the template")
	templateCreateCmd.Flags().StringVarP(&volumeIDCreate, "volume-id", "v", "", "The volume id for the template")
	templateCreateCmd.Flags().StringVarP(&nameCreate, "name", "n", "", "The name of the template")
	templateCreateCmd.Flags().StringVarP(&shortDescriptionCreate, "short-description", "s", "", "Add a short description")
	templateCreateCmd.Flags().StringVarP(&descriptionCreate, "description", "d", "", "Add a description")
	templateCreateCmd.Flags().StringVarP(&defaultUsernameCreate, "default-username", "u", "", "The default username of the template")
	templateCreateCmd.Flags().StringVarP(&cloudConfigCreate, "cloudconfig", "i", "", "The path of the cloud config")

	templateUpdateCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the template")
	templateUpdateCmd.Flags().StringVarP(&shortDescription, "short-description", "s", "", "Add a short description")
	templateUpdateCmd.Flags().StringVarP(&description, "description", "d", "", "Add a description")
	templateUpdateCmd.Flags().StringVarP(&defaultUsername, "default-username", "u", "", "The default username of the template")
	templateUpdateCmd.Flags().StringVarP(&cloudConfig, "cloudconfig", "i", "", "The cloud config")

}
