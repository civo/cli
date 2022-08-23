package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/civo/cli/common"
	"github.com/civo/cli/utility"
	"github.com/kierdavis/ansi"

	"github.com/spf13/cobra"
	"github.com/tj/go-update"
	"github.com/tj/go-update/progress"
	"github.com/tj/go-update/stores/github"
)

var (
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update the cli to the last version",
		Run: func(cmd *cobra.Command, args []string) {

			ansi.HideCursor()
			defer ansi.ShowCursor()

			// source update
			m := &update.Manager{
				Command: "civo",
				Store: &github.Store{
					Owner:   "civo",
					Repo:    "cli",
					Version: common.VersionCli,
				},
			}

			// fetch the new releases
			releases, err := m.LatestReleases()
			if err != nil {
				utility.Error("error fetching releases: %s", err)
				os.Exit(1)
			}

			// no updates
			if len(releases) == 0 {
				fmt.Printf("%s\n", utility.Green("Your client is up to date"))
				os.Exit(0)
			}

			// latest release
			latest := releases[0]

			s := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
			for _, a := range releases[0].Assets {
				ext := filepath.Ext(a.Name)
				if strings.Contains(a.Name, s) && ext == ".gz" {
					// download tarball to a tmp dir
					tarball, err := a.DownloadProxy(progress.Reader)
					if err != nil {
						utility.Error("error downloading: %s", err)
						os.Exit(1)
					}

					// install it
					if err := m.Install(tarball); err != nil {
						utility.Error("error installing: %s", err)
						os.Exit(1)
					}

					fmt.Printf("Updated to %s\n", utility.Green(latest.Version))

				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
}
