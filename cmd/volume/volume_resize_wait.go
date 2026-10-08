package volume

import (
	"fmt"
	"time"

	"github.com/civo/civogo"
)

// Volume states the API reports while a resize is being carried out.
const (
	volumeStatusResizing  = "resizing"
	volumeStatusMigrating = "migrating"
)

// resizeInFlight reports whether the API shows the volume mid-resize: either the volume's status
// is transitional, or the resize outcome says in_progress.
func resizeInFlight(v *civogo.Volume) bool {
	if v.Status == volumeStatusResizing || v.Status == volumeStatusMigrating {
		return true
	}
	return v.Resize != nil && !v.Resize.Settled()
}

// waitForResize polls find until the resize of a volume to requestedGB has settled, and returns
// the last volume read. It settles when the volume reports it delivers at least the requested
// size, or when the resize was seen in flight and has since left that state, or when the grace
// window passes without the resize ever being seen in flight (the API accepted the request before
// the platform picked it up, or the API predates the delivered size and there is nothing more to
// learn). It gives up with an error after timeout.
func waitForResize(find func() (*civogo.Volume, error), requestedGB int, grace, timeout, interval time.Duration, sleep func(time.Duration), now func() time.Time) (*civogo.Volume, error) {
	start := now()
	seenInFlight := false
	for {
		v, err := find()
		if err != nil {
			return nil, err
		}
		elapsed := now().Sub(start)
		switch {
		case resizeInFlight(v):
			seenInFlight = true
		case v.DeliveredSizeGigabytes >= requestedGB && v.DeliveredSizeGigabytes > 0:
			return v, nil
		case seenInFlight:
			return v, nil
		case elapsed >= grace:
			return v, nil
		}
		if elapsed >= timeout {
			return v, fmt.Errorf("timed out after %s waiting for the volume to finish resizing (status %q)", timeout, v.Status)
		}
		sleep(interval)
	}
}

// resizeOutcome judges a settled volume against the size that was requested. It returns an empty
// string when the resize delivered, otherwise a message explaining what the volume reports.
func resizeOutcome(v *civogo.Volume, requestedGB int) string {
	delivered := v.DeliveredSizeGigabytes
	if r := v.Resize; r != nil && r.Settled() && !r.Succeeded() && (delivered == 0 || delivered < requestedGB) {
		msg := fmt.Sprintf("the resize %s (%s)", r.State, r.Reason)
		if r.Message != "" {
			msg += ": " + r.Message
		}
		if delivered > 0 {
			msg += fmt.Sprintf("; the volume still delivers %d GB of the %d GB requested", delivered, requestedGB)
		}
		return msg
	}
	if delivered > 0 && delivered < requestedGB {
		return fmt.Sprintf("the volume still delivers %d GB of the %d GB requested", delivered, requestedGB)
	}
	return ""
}

// deliveredSizeLabel is the Delivered column of `civo volume ls`: the size the volume currently
// provides, blank when the API does not report it.
func deliveredSizeLabel(v civogo.Volume) string {
	if v.DeliveredSizeGigabytes == 0 {
		return ""
	}
	return fmt.Sprintf("%d GB", v.DeliveredSizeGigabytes)
}

// resizeLabel is the Resize column of `civo volume ls`: the state of the most recent resize and,
// when it did not complete, the platform's reason. Blank when no resize was ever attempted.
func resizeLabel(v civogo.Volume) string {
	r := v.Resize
	switch {
	case r == nil:
		return ""
	case r.Succeeded():
		return r.State
	default:
		return fmt.Sprintf("%s (%s)", r.State, r.Reason)
	}
}
