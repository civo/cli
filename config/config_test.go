package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/civo/cli/common"
)

func TestLoadConfig_VersionCheckOptOut(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name       string
		jsonConfig string
		wantOptOut bool
	}{
		{
			name:       "default config keeps opt-out false",
			jsonConfig: `{"apikeys":{},"meta":{"current_apikey":""}}`,
			wantOptOut: false,
		},
		{
			name:       "disable_version_check at root sets opt-out true",
			jsonConfig: `{"disable_version_check":true}`,
			wantOptOut: true,
		},
		{
			name:       "version_check false at root sets opt-out true",
			jsonConfig: `{"version_check":false}`,
			wantOptOut: true,
		},
		{
			name:       "disable_version_check in meta sets opt-out true",
			jsonConfig: `{"meta":{"disable_version_check":true}}`,
			wantOptOut: true,
		},
		{
			name:       "version_check false in meta sets opt-out true",
			jsonConfig: `{"meta":{"version_check":false}}`,
			wantOptOut: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			common.VersionCheckOptOut = false
			SkipAPIInitialization = true
			defer func() {
				common.VersionCheckOptOut = false
				SkipAPIInitialization = false
			}()

			cfgPath := filepath.Join(tempDir, "config.json")
			if err := os.WriteFile(cfgPath, []byte(tt.jsonConfig), 0600); err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			loadConfig(cfgPath)

			if common.VersionCheckOptOut != tt.wantOptOut {
				t.Errorf("got VersionCheckOptOut = %v, want %v", common.VersionCheckOptOut, tt.wantOptOut)
			}
		})
	}
}
