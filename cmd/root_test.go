package cmd

import (
	"slices"
	"testing"

	"github.com/civo/cli/common"
	"github.com/spf13/cobra"
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

func TestOutputFlagCompletesFormats(t *testing.T) {
	// Subcommands inherit --output from the root command, so completion has to work there too
	for _, cmd := range []*cobra.Command{rootCmd, versionCmd} {
		complete, ok := cmd.GetFlagCompletionFunc("output")
		if !ok {
			t.Fatalf("expected %q to have a completion function for --output", cmd.Name())
		}
		formats, directive := complete(cmd, nil, "")
		if !slices.Equal(formats, common.OutputFormats) {
			t.Errorf("expected %q to complete --output with %v, got %v", cmd.Name(), common.OutputFormats, formats)
		}
		if directive != cobra.ShellCompDirectiveNoFileComp {
			t.Errorf("expected %q to skip file completion for --output, got directive %d", cmd.Name(), directive)
		}
	}
}
