package utility

import (
	"io"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/cli/common"
)

// NewSpinner creates a progress spinner for long-running/--wait style
// commands. In --quiet mode, the spinner is created and can still be
// Start()/Stop()'d as normal by the caller, but its output is discarded so
// no progress indicator is printed - this keeps callers simple (no branching
// needed at each call site) while satisfying --quiet's "no progress
// indicators" requirement.
func NewSpinner(charSet []string, refreshRate time.Duration) *spinner.Spinner {
	s := spinner.New(charSet, refreshRate)
	if common.Quiet {
		s.Writer = io.Discard
	} else {
		s.Writer = os.Stderr
	}
	return s
}
