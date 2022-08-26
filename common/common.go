package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/tcnksm/go-latest"
)

var (
	// OutputFields for custom format output
	OutputFields string
	// OutputFormat for custom format output
	OutputFormat string
	// RegionSet picks the region to connect to, if you use this option it will use it over the default region
	RegionSet string
	// DefaultYes : automatic yes to prompts; assume \"yes\" as answer to all prompts and run non-interactively
	DefaultYes bool
	// PrettySet : Prints the json output in pretty format
	PrettySet bool
	// VersionCli is set from outside using ldflags
	VersionCli = "0.0.0"
	// CommitCli is set from outside using ldflags
	CommitCli = "none"
	// DateCli is set from outside using ldflags
	DateCli = "unknown"
)

func VersionCheck() *latest.CheckResponse {
	githubTag := &latest.GithubTag{
		Owner:             "civo",
		Repository:        "cli",
		FixVersionStrFunc: latest.DeleteFrontV(),
	}
	res, err := latest.Check(githubTag, strings.Replace(VersionCli, "v", "", 1))
	if err != nil {
		fmt.Printf("Checking for a newer version failed with %s \n", err)
		os.Exit(1)
	}
	return res
}
