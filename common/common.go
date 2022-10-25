package common

import (
	"fmt"
	"strings"

	"github.com/google/go-github/github"
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

// VersionCheck checks if there is a new version of the CLI
func VersionCheck() (res *latest.CheckResponse, skip bool) {
	githubTag := &latest.GithubTag{
		Owner:             "civo",
		Repository:        "cli",
		FixVersionStrFunc: latest.DeleteFrontV(),
	}
	res, err := latest.Check(githubTag, strings.Replace(VersionCli, "v", "", 1))
	if err != nil {
		if IsGHRatelimitError(err) {
			return nil, true
		}
		fmt.Printf("Checking for a newer version failed with %s \n", err)
		return nil, true
	}
	return res, false
}

// IsGHRatelimitError checks if the error is a github rate limit error
func IsGHRatelimitError(err error) bool {
	_, ok := err.(*github.RateLimitError)
	return ok
}
