package cmd

import (
	"fmt"
	"runtime"

	"github.com/civo/cli/common"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var (
	quiet   bool
	verbose bool

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Run: func(cmd *cobra.Command, args []string) {
			switch {
			case verbose:
				fmt.Printf("Client version: v%s\n", common.VersionCli)
				fmt.Printf("Go version (client): %s\n", runtime.Version())
				fmt.Printf("Build date (client): %s\n", common.DateCli)
				fmt.Printf("Git commit (client): %s\n", common.CommitCli)
				fmt.Printf("OS/Arch (client): %s/%s\n", runtime.GOOS, runtime.GOARCH)
			case quiet:
				fmt.Printf("v%s\n", common.VersionCli)
			default:
				utility.Error("Testing Error")
				fmt.Printf("Civo CLI v%s\n", common.VersionCli)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Use quiet output for simple output")
	versionCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Use verbose output to see full information")
}
