package utility

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/cli/common"
)

func TestNewSpinner_WriterSelection(t *testing.T) {
	originalQuiet := common.Quiet
	defer func() { common.Quiet = originalQuiet }()

	t.Run("quiet mode discards spinner output", func(t *testing.T) {
		common.Quiet = true
		s := NewSpinner(spinner.CharSets[9], 100*time.Millisecond)
		if s.Writer != io.Discard {
			t.Fatal("expected spinner Writer to be io.Discard when --quiet is set")
		}
	})

	t.Run("non-quiet mode writes to stderr as before", func(t *testing.T) {
		common.Quiet = false
		s := NewSpinner(spinner.CharSets[9], 100*time.Millisecond)
		if s.Writer != os.Stderr {
			t.Fatal("expected spinner Writer to be os.Stderr when --quiet is not set")
		}
	})
}
