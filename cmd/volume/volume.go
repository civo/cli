package volume

import (
	"errors"

	"github.com/spf13/cobra"
)

// VolumeCmd is the volume command
var VolumeCmd = &cobra.Command{
	Use:     "volume",
	Aliases: []string{"volumes"},
	Short:   "Details of Civo volumes",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	VolumeCmd.AddCommand(volumeCreateCmd)
	VolumeCmd.AddCommand(volumeListCmd)
	VolumeCmd.AddCommand(volumeRemoveCmd)
	// volumeCmd.AddCommand(volumeResizeCmd)
	VolumeCmd.AddCommand(volumeAttachCmd)
	VolumeCmd.AddCommand(volumeDetachCmd)

	volumeCreateCmd.Flags().IntVarP(&createSizeGB, "size-gb", "s", 0, "The new size in GB (required)")
	volumeCreateCmd.Flags().StringVarP(&networkVolumeID, "network", "t", "default", "The network name/ID where the volume will be created")
	volumeCreateCmd.MarkFlagRequired("size-gb")

	volumeResizeCmd.Flags().IntVarP(&newSizeGB, "size-gb", "s", 0, "The new size in GB (required)")
	volumeResizeCmd.MarkFlagRequired("size-gb")

	volumeAttachCmd.Flags().BoolVarP(&waitVolumeAttach, "wait", "w", false, "wait until the volume is attached")

	volumeDetachCmd.Flags().BoolVarP(&waitVolumeDetach, "wait", "w", false, "wait until the volume is detached")
}
