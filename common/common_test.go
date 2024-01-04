package common

import (
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
		if release.TagName == nil || *release.TagName != "v1.0.0" {
			t.Errorf("Expected release version to be v1.0.0, got %s", *release.TagName)
		}
	})
}
