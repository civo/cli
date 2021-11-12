package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-latest"
)

var (
	quiet   bool
	verbose bool

	// VersionCli is set from outside using ldflags
	VersionCli = "0.0.0"

	// CommitCli is set from outside using ldflags
	CommitCli = "none"

	// DateCli is set from outside using ldflags
	DateCli = "unknown"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Run: func(cmd *cobra.Command, args []string) {

			githubTag := &latest.GithubTag{
				Owner:             "civo",
				Repository:        "cli",
				FixVersionStrFunc: latest.DeleteFrontV(),
			}
			switch {
			case verbose:
				fmt.Printf("Client version: v%s\n", VersionCli)
				fmt.Printf("Go version (client): %s\n", runtime.Version())
				fmt.Printf("Build date (client): %s\n", DateCli)
				fmt.Printf("Git commit (client): %s\n", CommitCli)
				fmt.Printf("OS/Arch (client): %s/%s\n", runtime.GOOS, runtime.GOARCH)

				res, err := latest.Check(githubTag, strings.Replace(VersionCli, "v", "", 1))
				if err != nil {
					utility.Error("Checking for a newer version failed with %s", err)
					os.Exit(1)
				}

				if res.Outdated {
					utility.RedConfirm("A newer version (v%s) is available, please upgrade with \"civo update\"\n", res.Current)
				}
			case quiet:
				fmt.Printf("v%s\n", VersionCli)
			default:
				fmt.Printf("Civo CLI v%s\n", VersionCli)

				res, err := latest.Check(githubTag, strings.Replace(VersionCli, "v", "", 1))
				if err != nil {
					utility.Error("Checking for a newer version failed with %s", err)
					os.Exit(1)
				}

				if res.Outdated {
					utility.RedConfirm("A newer version (v%s) is available, please upgrade with \"civo update\"\n", res.Current)
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Use quiet output for simple output")
	versionCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Use verbose output to see full information")
}
