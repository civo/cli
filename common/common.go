package common

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
	"github.com/savioxavier/termlink"
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

const (
	// OutputFormatHuman for human readable output
	OutputFormatHuman = "human"
	// OutputFormatJSON for JSON output
	OutputFormatJSON = "json"
	// OutputFormatCustom for custom output
	OutputFormatCustom = "custom"
)

// GithubClient Create a Github client
func GithubClient() *github.Client {
	return github.NewClient(nil)
}

// CheckVersionUpdate checks if there's an update to be done
func CheckVersionUpdate() {
	ghClient := GithubClient()
	res, skip := VersionCheck(ghClient)
	if skip {
		return
	}

	// Check if the version is different from the one in the binary
	if res.TagName != nil && *res.TagName != fmt.Sprintf("v%s", VersionCli) {
		if res.TagName != nil && *res.TagName != VersionCli {
			fmt.Printf("A newer version (%s) is available, please upgrade with \"civo update\"\n", *res.TagName)
		}
	}
}

// IssueMessage is the message to be displayed when an error is returned
func IssueMessage() {
	gitIssueLink := termlink.ColorLink("GitHub issue", "https://github.com/civo/cli/issues", "green")
	fmt.Printf("Please check if you are using the latest version of CLI and retry the command \nIf you are still facing issues, please report it on our community slack or open a %s \n", gitIssueLink)
}

// VersionCheck checks if there is a new version of the CLI
func VersionCheck(client *github.Client) (res *github.RepositoryRelease, skip bool) {
	// Get the last release from GitHub API
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "civo", "cli")
	if _, ok := err.(*github.AbuseRateLimitError); ok {
		fmt.Printf("hit secondary rate limit try again in %s minute", err.(*github.AbuseRateLimitError).RetryAfter)
		return nil, true
	}
	if err != nil {
		return nil, true
	}
	return release, false
}
