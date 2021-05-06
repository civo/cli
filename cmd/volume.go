package cmd

import (
	"github.com/spf13/cobra"
)

var volumeCmd = &cobra.Command{
	Use:     "volume",
	Aliases: []string{"volumes"},
	Short:   "Details of Civo volumes",
}

func init() {
	// rootCmd.AddCommand(volumeCmd)
	volumeCmd.AddCommand(volumeListCmd)
	volumeCmd.AddCommand(volumeCreateCmd)
	volumeCmd.AddCommand(volumeResizeCmd)
	volumeCmd.AddCommand(volumeRemoveCmd)
	volumeCmd.AddCommand(volumeAttachCmd)
	volumeCmd.AddCommand(volumeDetachCmd)

	volumeCreateCmd.Flags().BoolVarP(&bootableVolume, "bootable", "b", false, "Mark the volume as bootable")
	volumeCreateCmd.Flags().IntVarP(&createSizeGB, "size-gb", "s", 0, "The new size in GB (required)")
	volumeCreateCmd.Flags().StringVarP(&networkVolumeID, "network", "t", "default", "The network where the volume will be created")
	volumeCreateCmd.MarkFlagRequired("size-gb")

	volumeResizeCmd.Flags().IntVarP(&newSizeGB, "size-gb", "s", 0, "The new size in GB (required)")
	volumeResizeCmd.MarkFlagRequired("size-gb")

	volumeAttachCmd.Flags().BoolVarP(&waitVolumeAttach, "wait", "w", false, "wait until the volume is attached")

	volumeDetachCmd.Flags().BoolVarP(&waitVolumeDetach, "wait", "w", false, "wait until the volume is detached")
}
