package cmd

import (
	"testing"

	"github.com/civo/cli/common"
)

func TestOutputFlagRejectsUnsupportedFormats(t *testing.T) {
	output := rootCmd.PersistentFlags().Lookup("output")
	if output == nil {
		t.Fatal("expected the root command to have an --output flag")
	}
	t.Cleanup(func() {
		if err := output.Value.Set(common.OutputFormatHuman); err != nil {
			t.Errorf("resetting --output failed: %v", err)
		}
	})

	if err := output.Value.Set("lefafa"); err == nil {
		t.Error("expected --output lefafa to be rejected")
	}
	if err := output.Value.Set(common.OutputFormatJSON); err != nil {
		t.Errorf("expected --output json to be accepted, got %v", err)
	}
	if common.OutputFormat != common.OutputFormatJSON {
		t.Errorf("expected common.OutputFormat to be %q, got %q", common.OutputFormatJSON, common.OutputFormat)
	}
}
