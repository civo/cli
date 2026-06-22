package kubernetes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMarketplaceAppUninstallScript(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics-server/uninstall.sh" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("echo uninstalling"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	tests := []struct {
		name       string
		appName    string
		wantScript string
		wantErr    bool
	}{
		{
			name:       "returns the script for an existing app",
			appName:    "metrics-server",
			wantScript: "echo uninstalling",
			wantErr:    false,
		},
		{
			name:    "returns an error for a missing app",
			appName: "does-not-exist",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getMarketplaceAppUninstallScript(srv.URL, tt.appName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getMarketplaceAppUninstallScript(%q) error = %v, wantErr %v", tt.appName, err, tt.wantErr)
			}
			if !tt.wantErr && string(got) != tt.wantScript {
				t.Errorf("got %q, want %q", string(got), tt.wantScript)
			}
		})
	}
}
