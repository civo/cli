package diskimage

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var diskImageFindCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "search", "find"},
	Example: `civo diskimage show`,
	Short:   "Finds a disk image by either part of the ID or part of the name",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		diskImg, err := client.FindDiskImage(args[0])
		if err != nil {
			utility.Error("Disk Image %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		ow.AppendDataWithLabel("id", diskImg.ID, "ID")
		ow.AppendDataWithLabel("name", diskImg.Name, "Name")
		ow.AppendDataWithLabel("version", diskImg.Version, "Version")
		ow.AppendDataWithLabel("state", diskImg.State, "State")
		ow.AppendDataWithLabel("distribution", diskImg.Distribution, "Distribution")
		ow.AppendDataWithLabel("initial_user", diskImg.InitialUser, "Initial User")
		ow.AppendDataWithLabel("os", diskImg.OS, "OS")
		ow.AppendDataWithLabel("description", diskImg.Description, "Description")
		ow.AppendDataWithLabel("label", diskImg.Label, "Label")
		ow.AppendDataWithLabel("disk_image_url", diskImg.DiskImageURL, "Disk Image URL")
		ow.AppendDataWithLabel("disk_image_size_bytes", fmt.Sprintf("%d", diskImg.DiskImageSizeBytes), "Disk Image Size (bytes)")
		ow.AppendDataWithLabel("logo_url", diskImg.LogoURL, "Logo URL")
		ow.AppendDataWithLabel("created_at", diskImg.CreatedAt.String(), "Created At")
		ow.AppendDataWithLabel("created_by", diskImg.CreatedBy, "Created By")
		ow.AppendDataWithLabel("distribution_default", fmt.Sprintf("%t", diskImg.DistributionDefault), "Distribution Default")

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteKeyValues()
		}
	},
}
