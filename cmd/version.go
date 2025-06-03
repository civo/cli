package cmd

import (
	"runtime"

	"github.com/civo/cli/common"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

const logo = `
  _____  _                  _____ ___      _____ 
/ ____| (_)               / ____| | |      |_ _|
| |      _ __   __  ___   | |     | |       | |  
| |     | |\ \ / / / _ \  | |     | |       | |  
| |____ | | \ V / | (_) | | |____ | |____  _| |_ 
 \_____||_|  \_/   \___/   \_____||______||_____|
`

var (
	quiet   bool
	verbose bool

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Run: func(cmd *cobra.Command, args []string) {
			switch {
			case verbose:
				utility.Printf(logo)
				utility.Printf("Client version: v%s\n", common.VersionCli)
				utility.Printf("Go version (client): %s\n", runtime.Version())
				utility.Printf("Build date (client): %s\n", common.DateCli)
				utility.Printf("Git commit (client): %s\n", common.CommitCli)
				utility.Printf("OS/Arch (client): %s/%s\n", runtime.GOOS, runtime.GOARCH)
				common.CheckVersionUpdate()
			case quiet:
				utility.Printf("v%s\n", common.VersionCli)
			default:
				common.CheckVersionUpdate()
				utility.Printf("Civo CLI v%s\n", common.VersionCli)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Use quiet output for simple output")
	versionCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Use verbose output to see full information")
}
