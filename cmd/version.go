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
	quiet      = false
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
			res, _ := latest.Check(githubTag, strings.Replace(VersionCli, "v", "", 1))

			if !quiet {
				fmt.Printf("Client version: %s\n", VersionCli)
				fmt.Printf("Go version (client): %s\n", runtime.Version())
				fmt.Printf("Git commit (client): %s\n", CommitCli)
				fmt.Printf("OS/Arch (client): %s/%s\n", runtime.GOOS, runtime.GOARCH)

				if res.Outdated {
					utility.YellowConfirm("A newer version (v%s) is available, please upgrade\n", res.Current)
				}

			} else {
				fmt.Printf("%s\n", VersionCli)
			}

		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Use quiet output for version information.")
}
