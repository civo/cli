package cmd

import (
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

type DiskImage struct {
	ID           string
	Name         string
	Version      string
	State        string
	Distribution string
}

var diskImageListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo diskimage ls`,
	Short:   "List diskimages",
	Long: `List all available diskimages.
If you wish to use a custom format, the available fields are:

	* id
	* name
	* version
	* state
	* distribution

Example: civo diskimage ls -o=custom -f=id,name`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		diskImageList := []DiskImage{}

		diskImages, err := client.ListDiskImages()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		for _, v := range diskImages {
			diskImageList = append(diskImageList, DiskImage{
				ID:           v.ID,
				Name:         v.Name,
				Version:      v.Version,
				State:        v.State,
				Distribution: v.Distribution,
			})
		}

		ow := utility.NewOutputWriter()

		for _, diskImage := range diskImageList {
			ow.StartLine()
			ow.AppendDataWithLabel("id", diskImage.ID, "ID")
			ow.AppendDataWithLabel("name", diskImage.Name, "Name")
			ow.AppendDataWithLabel("version", diskImage.Version, "Version")
			ow.AppendDataWithLabel("state", diskImage.State, "State")
			ow.AppendDataWithLabel("distribution", diskImage.Distribution, "Distribution")
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
