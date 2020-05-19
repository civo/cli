package cmd

import (
	"github.com/spf13/cobra"
)

var volumeCmd = &cobra.Command{
	Use:     "volume",
	Aliases: []string{"volumes"},
	Short:   "Details of Civo volume",
}

func init() {
	rootCmd.AddCommand(volumeCmd)
	volumeCmd.AddCommand(volumeListCmd)
	volumeCmd.AddCommand(volumeCreateCmd)
	volumeCmd.AddCommand(volumeResizeCmd)
	volumeCmd.AddCommand(volumeRemoveCmd)
	volumeCmd.AddCommand(volumeAttachCmd)
	volumeCmd.AddCommand(volumeDetachCmd)

	/*
		Flags for the create cmd
	*/
	volumeCreateCmd.Flags().BoolVarP(&bootableVolume, "bootable", "b", false, "Mark the volume as bootable")
	volumeCreateCmd.Flags().IntVarP(&createSizeGB, "size-gb", "s", 0, "The new size in GB (required)")
	volumeCreateCmd.MarkFlagRequired("size-gb")

	/*
		Flags for the resize cmd
	*/
	volumeResizeCmd.Flags().IntVarP(&newSizeGB, "size-gb", "s", 0, "The new size in GB (required)")
	volumeResizeCmd.MarkFlagRequired("size-gb")

	/*
		Flags for the attach cmd
	*/
	volumeAttachCmd.Flags().BoolVarP(&waitVolumeAttach, "wait", "w", false, "wait until the volume is attached")

	/*
		Flags for the attach cmd
	*/
	volumeDetachCmd.Flags().BoolVarP(&waitVolumeDetach, "wait", "w", false, "wait until the volume is detached")

}
