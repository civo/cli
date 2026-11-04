package common

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v57/github"
)

// FILEPATH: /Users/alejandrojnm/Project/go/src/github.com/civo/cli/common/common_test.go

func TestVersionCheck(t *testing.T) {
	// Test when the GitHub API returns a successful response
	t.Run("successful response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(`{"tag_name": "v1.0.0"}`))
		}))
		defer server.Close()

		client := github.NewClient(nil)
		client.BaseURL, _ = url.Parse(fmt.Sprintf("%s/", server.URL))

		release, skip := VersionCheck(client)
		if skip {
			t.Errorf("Expected skip to be false, got true")
		}
		if release == nil || release.TagName == nil || *release.TagName != "v1.0.0" {
			t.Errorf("Expected release version to be v1.0.0, got %v", release)
		}
	})

	t.Run("opt-out via VersionCheckOptOut", func(t *testing.T) {
		oldOptOut := VersionCheckOptOut
		defer func() { VersionCheckOptOut = oldOptOut }()
		VersionCheckOptOut = true

		release, skip := VersionCheck(nil)
		if !skip || release != nil {
			t.Errorf("expected skip=true and release=nil when opted out, got skip=%v, release=%v", skip, release)
		}
	})

	t.Run("opt-out via CIVO_NO_VERSION_CHECK", func(t *testing.T) {
		t.Setenv("CIVO_NO_VERSION_CHECK", "1")
		release, skip := VersionCheck(nil)
		if !skip || release != nil {
			t.Errorf("expected skip=true and release=nil when CIVO_NO_VERSION_CHECK is set")
		}
	})

	t.Run("opt-out via CIVO_DISABLE_VERSION_CHECK", func(t *testing.T) {
		t.Setenv("CIVO_DISABLE_VERSION_CHECK", "true")
		release, skip := VersionCheck(nil)
		if !skip || release != nil {
			t.Errorf("expected skip=true and release=nil when CIVO_DISABLE_VERSION_CHECK is set")
		}
	})

	t.Run("nil client returns skip=true without panic", func(t *testing.T) {
		release, skip := VersionCheck(nil)
		if !skip || release != nil {
			t.Errorf("expected skip=true and release=nil for nil client")
		}
	})
}

func TestCheckVersionUpdate(t *testing.T) {
	tests := []struct {
		name       string
		latestTag  string
		cliVersion string
		optOut     bool
		envVar     string
		wantOutput bool
	}{
		{
			name:       "newer version available prints notice to stderr",
			latestTag:  "v1.1.0",
			cliVersion: "1.0.0",
			wantOutput: true,
		},
		{
			name:       "newer version without v prefix prints notice to stderr",
			latestTag:  "1.1.0",
			cliVersion: "v1.0.0",
			wantOutput: true,
		},
		{
			name:       "same version produces no output",
			latestTag:  "v1.0.0",
			cliVersion: "1.0.0",
			wantOutput: false,
		},
		{
			name:       "binary newer than latest release produces no output",
			latestTag:  "v1.0.0",
			cliVersion: "1.1.0",
			wantOutput: false,
		},
		{
			name:       "dev build 0.0.0 produces no output",
			latestTag:  "v1.2.0",
			cliVersion: "0.0.0",
			wantOutput: false,
		},
		{
			name:       "opt-out via variable suppresses check",
			latestTag:  "v2.0.0",
			cliVersion: "1.0.0",
			optOut:     true,
			wantOutput: false,
		},
		{
			name:       "opt-out via environment variable suppresses check",
			latestTag:  "v2.0.0",
			cliVersion: "1.0.0",
			envVar:     "1",
			wantOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldVersion := VersionCli
			oldOptOut := VersionCheckOptOut
			defer func() {
				VersionCli = oldVersion
				VersionCheckOptOut = oldOptOut
			}()

			VersionCli = tt.cliVersion
			VersionCheckOptOut = tt.optOut
			if tt.envVar != "" {
				t.Setenv("CIVO_NO_VERSION_CHECK", tt.envVar)
			}

			requestMade := false
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				requestMade = true
				rw.Write([]byte(fmt.Sprintf(`{"tag_name": %q}`, tt.latestTag)))
			}))
			defer server.Close()

			client := github.NewClient(nil)
			client.BaseURL, _ = url.Parse(fmt.Sprintf("%s/", server.URL))

			var stderr bytes.Buffer
			checkVersionUpdateWithClient(client, &stderr)

			if tt.optOut || tt.envVar != "" {
				if requestMade {
					t.Errorf("expected no HTTP request when opted out, but request was made")
				}
			}

			output := stderr.String()
			if tt.wantOutput && output == "" {
				t.Errorf("expected update warning in stderr, got empty")
			} else if !tt.wantOutput && output != "" {
				t.Errorf("expected no output, got: %q", output)
			}
		})
	}
}

func TestOutputFormatFlag(t *testing.T) {
	var format string
	flag := NewOutputFormatFlag(&format, OutputFormatHuman)
	if format != OutputFormatHuman {
		t.Errorf("expected default %q, got %q", OutputFormatHuman, format)
	}

	for _, value := range OutputFormats {
		if err := flag.Set(value); err != nil {
			t.Errorf("expected %q to be accepted, got %v", value, err)
		}
		if format != value {
			t.Errorf("expected format to be %q, got %q", value, format)
		}
	}

	// Matching is case-sensitive, so "JSON" is rejected rather than silently printing human output
	before := format
	for _, value := range []string{"", "lefafa", "JSON", "yaml", "table"} {
		if err := flag.Set(value); err == nil {
			t.Errorf("expected %q to be rejected", value)
		}
		if format != before {
			t.Errorf("rejected value %q changed format to %q", value, format)
		}
	}
}
