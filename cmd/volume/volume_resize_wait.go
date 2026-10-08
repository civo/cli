package volume

import (
	"fmt"
	"time"

	"github.com/civo/civogo"
)

// volumeStatusResizing is the volume status the API reports while the platform carries out a resize.
const volumeStatusResizing = "resizing"

// How a followed resize ended.
const (
	resizeDelivered   = "delivered"   // the volume delivers the requested size
	resizeFailed      = "failed"      // the platform settled the resize without delivering it
	resizeUnconfirmed = "unconfirmed" // the API gave no evidence either way
)

// resizeWait follows one resize request to its outcome by polling the volume.
type resizeWait struct {
	// find reads the volume afresh.
	find func() (*civogo.Volume, error)
	// before is the volume as it was read just before the resize request: the baseline that lets
	// the wait tell an outcome of this request from one left over from a previous resize.
	before *civogo.Volume
	// requested is the size asked for, in GB.
	requested int
	// grace bounds the wait for an API that does not report the delivered size, where the only
	// possible evidence is the transitional status; timeout bounds the whole wait.
	grace, timeout, interval time.Duration
	// maxFindFailures is how many consecutive find errors are tolerated before giving up.
	maxFindFailures int
	sleep           func(time.Duration)
	now             func() time.Time
}

// resizeResult is what a wait returns: the last volume read, how the resize ended, and a
// human-readable detail for the failed and unconfirmed cases.
type resizeResult struct {
	volume *civogo.Volume
	state  string
	detail string
}

// resizeInFlight reports whether the API shows the volume mid-resize: the transitional status, or a
// resize outcome that says in progress.
func resizeInFlight(v *civogo.Volume) bool {
	return v.Status == volumeStatusResizing || (v.Resize != nil && !v.Resize.Settled())
}

// resizeFromThisRequest reports whether v.Resize describes the request being followed rather than an
// earlier resize: it is in flight, or it did not exist before the request, or it changed since.
func (w resizeWait) resizeFromThisRequest(v *civogo.Volume) bool {
	r := v.Resize
	if r == nil {
		return false
	}
	if !r.Settled() {
		return true
	}
	if w.before == nil || w.before.Resize == nil {
		return true
	}
	return !r.UpdatedAt.Equal(w.before.Resize.UpdatedAt)
}

// deliveredNow reports whether the volume delivers the requested size, and did not already do so
// before the request (a no-op would otherwise pass as evidence).
func (w resizeWait) deliveredNow(v *civogo.Volume) bool {
	if v.DeliveredSizeGigabytes <= 0 || v.DeliveredSizeGigabytes < w.requested {
		return false
	}
	return w.before == nil || w.before.DeliveredSizeGigabytes < w.requested
}

// evidence reports whether the volume shows any sign that this request was picked up.
func (w resizeWait) evidence(v *civogo.Volume) bool {
	return resizeInFlight(v) || w.resizeFromThisRequest(v) || w.deliveredNow(v)
}

// reportsDeliveredSize reports whether the API serves delivered_size_gb at all, judged on the
// baseline read. An API that does gives a definite outcome for every resize; one that does not can
// only show the transitional status, and may not even do that for a resize it has not started.
func (w resizeWait) reportsDeliveredSize() bool {
	return w.before != nil && w.before.DeliveredSizeGigabytes > 0
}

// run polls until the resize has an outcome. It returns an error only when the wait itself breaks
// down (the API keeps failing, or the timeout passes); every outcome of the resize, including
// failure, is reported in the result.
func (w resizeWait) run() (resizeResult, error) {
	start := w.now()
	seen := false
	failures := 0
	var last *civogo.Volume
	for {
		v, err := w.find()
		if err != nil {
			failures++
			if failures > w.maxFindFailures {
				return resizeResult{volume: last}, fmt.Errorf("reading the volume failed %d times in a row, last error: %w", failures, err)
			}
			w.sleep(w.interval)
			continue
		}
		failures = 0
		last = v
		elapsed := w.now().Sub(start)

		if w.evidence(v) {
			seen = true
		}
		if seen && !resizeInFlight(v) {
			switch {
			case w.deliveredNow(v):
				return resizeResult{volume: v, state: resizeDelivered}, nil
			case w.resizeFromThisRequest(v) && v.Resize.Succeeded():
				return resizeResult{volume: v, state: resizeDelivered}, nil
			case w.resizeFromThisRequest(v):
				return resizeResult{volume: v, state: resizeFailed, detail: resizeFailureDetail(v, w.requested)}, nil
			case !w.reportsDeliveredSize():
				// The transitional status came and went on an API that reports nothing else.
				return resizeResult{volume: v, state: resizeUnconfirmed, detail: "the resize has settled, but this API does not report the delivered size"}, nil
			}
			// Delivered size still below the request and no outcome of this request yet: the
			// platform is between steps. Keep polling.
		}
		if !seen && !w.reportsDeliveredSize() && elapsed >= w.grace {
			return resizeResult{volume: v, state: resizeUnconfirmed, detail: fmt.Sprintf("the resize was not seen in progress within %s, and this API does not report the delivered size", w.grace)}, nil
		}
		if elapsed >= w.timeout {
			return resizeResult{volume: v}, fmt.Errorf("timed out after %s waiting for the resize to settle (status %q)", w.timeout, v.Status)
		}
		w.sleep(w.interval)
	}
}

// resizeFailureDetail words a settled, unsuccessful resize of this request.
func resizeFailureDetail(v *civogo.Volume, requestedGB int) string {
	r := v.Resize
	msg := "the resize " + r.State
	if r.Reason != "" {
		msg += " (" + r.Reason + ")"
	}
	if r.Message != "" {
		msg += ": " + r.Message
	}
	if v.DeliveredSizeGigabytes > 0 && v.DeliveredSizeGigabytes < requestedGB {
		msg += fmt.Sprintf("; the volume still delivers %d GB of the %d GB requested", v.DeliveredSizeGigabytes, requestedGB)
	}
	return msg
}

// deliveredSizeLabel is the Delivered column of `civo volume ls`: the size the volume currently
// provides, blank when the API does not report it.
func deliveredSizeLabel(v civogo.Volume) string {
	if v.DeliveredSizeGigabytes <= 0 {
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
	case r.Succeeded() || r.Reason == "":
		return r.State
	default:
		return fmt.Sprintf("%s (%s)", r.State, r.Reason)
	}
}

// reportsResize reports whether any volume in the list carries the delivered size or a resize
// outcome, so `civo volume ls` adds the two columns only when there is something to show and the
// custom-output field order stays the same on APIs that do not report them.
func reportsResize(volumes []civogo.Volume) bool {
	for _, v := range volumes {
		if v.DeliveredSizeGigabytes > 0 || v.Resize != nil {
			return true
		}
	}
	return false
}
