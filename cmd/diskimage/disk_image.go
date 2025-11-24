package diskimage

import (
	"errors"

	"github.com/spf13/cobra"
)

// DiskImageCmd manages disk images
var DiskImageCmd = &cobra.Command{
	Use:     "diskimage",
	Aliases: []string{"diskimages", "template", "templates"},
	Short:   "Details of Civo disk images",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	DiskImageCmd.AddCommand(diskImageListCmd)
	DiskImageCmd.AddCommand(diskImageFindCmd)
	DiskImageCmd.AddCommand(diskImageCreateCmd)
	DiskImageCmd.AddCommand(diskImageDeleteCmd)

	diskImageCreateCmd.Flags().StringVarP(&createDiskImageName, "name", "n", "", "Name of the disk image")
	diskImageCreateCmd.Flags().StringVarP(&createDiskImageDistribution, "distribution", "d", "", "Distribution name (e.g. ubuntu, centos)")
	diskImageCreateCmd.Flags().StringVarP(&createDiskImageVersion, "version", "v", "", "Version of the distribution")
	diskImageCreateCmd.Flags().StringVarP(&createDiskImagePath, "path", "p", "", "Path to disk image file (.raw/.qcow2)")
	diskImageCreateCmd.Flags().StringVarP(&createOS, "os", "t", "linux", "Operating system type (linux/windows)")
	diskImageCreateCmd.Flags().StringVarP(&createLogoPath, "logo_path", "l", "", "Path to SVG logo file")
	_ = diskImageCreateCmd.MarkFlagRequired("name")
	_ = diskImageCreateCmd.MarkFlagRequired("distribution")
	_ = diskImageCreateCmd.MarkFlagRequired("version")
	_ = diskImageCreateCmd.MarkFlagRequired("path")

	diskImageListCmd.Flags().BoolVar(&showCustomImages, "custom", false, "Show only your custom disk images")
}
