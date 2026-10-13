package utility

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/civo/cli/common"
)

// captureStderr redirects os.Stderr for the duration of fn and returns what was written to it.
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	originalStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stderr = w

	fn()

	w.Close()
	os.Stderr = originalStderr

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestInfoWarning_SuppressedWhenQuiet(t *testing.T) {
	originalQuiet := common.Quiet
	defer func() { common.Quiet = originalQuiet }()

	common.Quiet = true
	out := captureStderr(t, func() {
		Info("some info")
		Warning("some warning")
		YellowConfirm("some confirm banner")
		RedConfirm("some important banner")
	})
	if out != "" {
		t.Fatalf("expected no output for Info/Warning/YellowConfirm/RedConfirm when --quiet is set, got %q", out)
	}
}

func TestInfoWarning_ShownWhenNotQuiet(t *testing.T) {
	originalQuiet := common.Quiet
	defer func() { common.Quiet = originalQuiet }()

	common.Quiet = false
	out := captureStderr(t, func() {
		Info("some info")
	})
	if out == "" {
		t.Fatal("expected Info to print when --quiet is not set")
	}
}

func TestError_NeverSuppressedEvenWhenQuiet(t *testing.T) {
	originalQuiet := common.Quiet
	defer func() { common.Quiet = originalQuiet }()

	common.Quiet = true
	out := captureStderr(t, func() {
		Error("something failed")
	})
	if out == "" {
		t.Fatal("expected Error to always print, even when --quiet is set")
	}
}
