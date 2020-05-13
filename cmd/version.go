package cmd

import (
	"fmt"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-latest"
	"runtime"
	"strings"
)

var (
	quiet      bool
	verbose    bool
	VersionCli = "dev"
	CommitCli  = "none"
	DateCli    = "unknown"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Run: func(cmd *cobra.Command, args []string) {

			githubTag := &latest.GithubTag{
				Owner:             "civo",
				Repository:        "cli-go",
				FixVersionStrFunc: latest.DeleteFrontV(),
			}
			switch {
			case verbose:
				fmt.Printf("Client version: v%s\n", VersionCli)
				fmt.Printf("Go version (client): %s\n", runtime.Version())
				fmt.Printf("Git commit (client): %s\n", CommitCli)
				fmt.Printf("OS/Arch (client): %s/%s\n", runtime.GOOS, runtime.GOARCH)

				res, _ := latest.Check(githubTag, strings.Replace(VersionCli, "v", "", 1))

				if res.Outdated {
					utility.RedConfirm("A newer version (v%s) is available, please upgrade\n", res.Current)
				}
			case quiet:
				fmt.Printf("%s\n", VersionCli)
			default:
				fmt.Printf("Civo CLI v%s\n", VersionCli)

				res, _ := latest.Check(githubTag, strings.Replace(VersionCli, "v", "", 1))

				if res.Outdated {
					utility.RedConfirm("A newer version (v%s) is available, please upgrade\n", res.Current)
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Use quiet output for simple output.")
	versionCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Use verbose output to see full information.")
}
