package common

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/savioxavier/termlink"
	"golang.org/x/mod/semver"
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
	// VersionCheckOptOut can be set to true to disable checking for new CLI versions
	VersionCheckOptOut bool
)

const (
	// OutputFormatHuman for human readable output
	OutputFormatHuman = "human"
	// OutputFormatJSON for JSON output
	OutputFormatJSON = "json"
	// OutputFormatCustom for custom output
	OutputFormatCustom = "custom"
)

// OutputFormats lists the values accepted by the --output flag
var OutputFormats = []string{OutputFormatJSON, OutputFormatHuman, OutputFormatCustom}

// OutputFormatFlag is a flag value that only accepts one of OutputFormats,
// so an unsupported format fails while the flags are parsed
type OutputFormatFlag struct {
	target *string
}

// NewOutputFormatFlag returns an OutputFormatFlag that stores the chosen format in target, starting at defaultFormat
func NewOutputFormatFlag(target *string, defaultFormat string) *OutputFormatFlag {
	*target = defaultFormat
	return &OutputFormatFlag{target: target}
}

func (f *OutputFormatFlag) String() string {
	return *f.target
}

// Set stores value if it is a supported output format
func (f *OutputFormatFlag) Set(value string) error {
	if !slices.Contains(OutputFormats, value) {
		return fmt.Errorf("must be one of %s", strings.Join(OutputFormats, ", "))
	}
	*f.target = value
	return nil
}

// Type is shown next to the flag name in the help output
func (f *OutputFormatFlag) Type() string {
	return "string"
}

// IsVersionCheckDisabled reports whether checking for CLI updates is disabled.
// It returns true if VersionCheckOptOut is true or if either CIVO_NO_VERSION_CHECK
// or CIVO_DISABLE_VERSION_CHECK environment variable is set to a truthy value ("1", "true", "yes", "on").
func IsVersionCheckDisabled() bool {
	if VersionCheckOptOut {
		return true
	}
	for _, env := range []string{"CIVO_NO_VERSION_CHECK", "CIVO_DISABLE_VERSION_CHECK"} {
		if val, ok := os.LookupEnv(env); ok {
			val = strings.TrimSpace(strings.ToLower(val))
			if val != "" && val != "0" && val != "false" {
				return true
			}
		}
	}
	return false
}

// GithubClient Create a Github client
func GithubClient() *github.Client {
	return github.NewClient(nil)
}

// CheckVersionUpdate checks if there's an update to be done
func CheckVersionUpdate() {
	checkVersionUpdateWithClient(GithubClient(), os.Stderr)
}

func checkVersionUpdateWithClient(ghClient *github.Client, stderr io.Writer) {
	if IsVersionCheckDisabled() {
		return
	}

	res, skip := VersionCheck(ghClient)
	if skip || res == nil || res.TagName == nil {
		return
	}

	latest := *res.TagName
	latestV := latest
	if !strings.HasPrefix(latestV, "v") {
		latestV = "v" + latestV
	}

	currentV := VersionCli
	if !strings.HasPrefix(currentV, "v") {
		currentV = "v" + currentV
	}

	// Do not notify for development builds or invalid semver
	if !semver.IsValid(latestV) || !semver.IsValid(currentV) || currentV == "v0.0.0" {
		return
	}

	if semver.Compare(latestV, currentV) > 0 {
		fmt.Fprintf(stderr, "A newer version (%s) is available, please upgrade with \"civo update\"\n", latest)
	}
}

// IssueMessage is the message to be displayed when an error is returned
func IssueMessage() {
	gitIssueLink := termlink.ColorLink("GitHub issue", "https://github.com/civo/cli/issues", "green")
	fmt.Printf("Please check if you are using the latest version of CLI and retry the command \nIf you are still facing issues, please report it on our community slack or open a %s \n", gitIssueLink)
}

// VersionCheck checks if there is a new version of the CLI
func VersionCheck(client *github.Client) (res *github.RepositoryRelease, skip bool) {
	if IsVersionCheckDisabled() {
		return nil, true
	}
	if client == nil || client.Repositories == nil {
		return nil, true
	}

	// Get the last release from GitHub API
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "civo", "cli")
	if _, ok := err.(*github.AbuseRateLimitError); ok {
		fmt.Fprintf(os.Stderr, "hit secondary rate limit try again in %s minute\n", err.(*github.AbuseRateLimitError).RetryAfter)
		return nil, true
	}
	if err != nil {
		return nil, true
	}
	return release, false
}
