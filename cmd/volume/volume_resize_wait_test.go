package volume

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/civo/civogo"
)

// fakeClock drives resizeWait without real sleeps: every sleep advances the clock.
type fakeClock struct{ t time.Time }

func (c *fakeClock) now() time.Time        { return c.t }
func (c *fakeClock) sleep(d time.Duration) { c.t = c.t.Add(d) }

// script serves the given reads in order and repeats the last one; a nil entry is a read error.
func script(reads ...*civogo.Volume) func() (*civogo.Volume, error) {
	i := 0
	return func() (*civogo.Volume, error) {
		v := reads[i]
		if i < len(reads)-1 {
			i++
		}
		if v == nil {
			return nil, errors.New("502 bad gateway")
		}
		return v, nil
	}
}

const (
	testGrace    = time.Minute
	testTimeout  = 10 * time.Minute
	testInterval = 2 * time.Second
)

func newWait(before *civogo.Volume, requested int, clock *fakeClock, reads ...*civogo.Volume) resizeWait {
	return resizeWait{
		find: script(reads...), before: before, requested: requested,
		grace: testGrace, timeout: testTimeout, interval: testInterval, maxFindFailures: 3,
		sleep: clock.sleep, now: clock.now,
	}
}

func TestResizeWait(t *testing.T) {
	t0 := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	at := func(s int) time.Time { return t0.Add(time.Duration(s) * time.Second) }
	oldFailed := &civogo.VolumeResize{State: civogo.VolumeResizeFailed, Reason: "RequiresDetach", UpdatedAt: at(-3600)}
	oldDone := &civogo.VolumeResize{State: civogo.VolumeResizeCompleted, Reason: "Completed", UpdatedAt: at(-3600)}
	inProgress := &civogo.VolumeResize{State: civogo.VolumeResizeInProgress, Reason: "Expanding", UpdatedAt: at(5)}
	done := &civogo.VolumeResize{State: civogo.VolumeResizeCompleted, Reason: "Completed", UpdatedAt: at(40)}
	failed := &civogo.VolumeResize{State: civogo.VolumeResizeFailed, Reason: "InsufficientStorage", Message: "pool is full", UpdatedAt: at(17)}
	cancelled := &civogo.VolumeResize{State: civogo.VolumeResizeCancelled, UpdatedAt: at(17)}

	cases := []struct {
		name       string
		before     *civogo.Volume
		requested  int
		reads      []*civogo.Volume
		wantState  string
		wantDetail string
		wantErr    string
		quick      bool // must settle before the grace window
	}{
		{
			name:      "online resize seen in flight then delivered",
			before:    &civogo.Volume{Status: "attached", SizeGigabytes: 5, DeliveredSizeGigabytes: 5, Resize: oldDone},
			requested: 6,
			reads: []*civogo.Volume{
				{Status: "attached", SizeGigabytes: 6, DeliveredSizeGigabytes: 5, Resize: oldDone},
				{Status: "resizing", SizeGigabytes: 6, DeliveredSizeGigabytes: 5, Resize: inProgress},
				{Status: "attached", SizeGigabytes: 6, DeliveredSizeGigabytes: 6, Resize: done},
			},
			wantState: resizeDelivered, quick: true,
		},
		{
			name:      "refused resize settles below the request with the platform's reason",
			before:    &civogo.Volume{Status: "attached", SizeGigabytes: 6, DeliveredSizeGigabytes: 6, Resize: oldDone},
			requested: 300,
			reads: []*civogo.Volume{
				{Status: "resizing", SizeGigabytes: 300, DeliveredSizeGigabytes: 6, Resize: inProgress},
				{Status: "attached", SizeGigabytes: 300, DeliveredSizeGigabytes: 6, Resize: failed},
			},
			wantState:  resizeFailed,
			wantDetail: "the resize failed (InsufficientStorage): pool is full; the volume still delivers 6 GB of the 300 GB requested",
			quick:      true,
		},
		{
			name:      "a failed resize left over from an earlier request is not blamed on this one",
			before:    &civogo.Volume{Status: "available", SizeGigabytes: 20, DeliveredSizeGigabytes: 20, Resize: oldFailed},
			requested: 60,
			reads: []*civogo.Volume{
				{Status: "available", SizeGigabytes: 60, DeliveredSizeGigabytes: 20, Resize: oldFailed},
				{Status: "available", SizeGigabytes: 60, DeliveredSizeGigabytes: 20, Resize: oldFailed},
				{Status: "resizing", SizeGigabytes: 60, DeliveredSizeGigabytes: 20, Resize: inProgress},
				{Status: "available", SizeGigabytes: 60, DeliveredSizeGigabytes: 60, Resize: done},
			},
			wantState: resizeDelivered, quick: true,
		},
		{
			name:      "slow pickup with no evidence is reported as unconfirmed after the grace window",
			before:    &civogo.Volume{Status: "available", SizeGigabytes: 20, DeliveredSizeGigabytes: 20},
			requested: 60,
			reads:     []*civogo.Volume{{Status: "available", SizeGigabytes: 60, DeliveredSizeGigabytes: 20}},
			wantState: resizeUnconfirmed, wantDetail: "the resize was not seen in progress within 1m0s",
		},
		{
			name:      "the delivered size alone is not evidence: the API serves the request when the platform has not reported",
			before:    &civogo.Volume{Status: "attached", SizeGigabytes: 20, DeliveredSizeGigabytes: 20},
			requested: 60,
			reads:     []*civogo.Volume{{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 60}},
			wantState: resizeUnconfirmed, wantDetail: "the resize was not seen in progress within 1m0s",
		},
		{
			name:      "delivered before the first poll",
			before:    &civogo.Volume{Status: "attached", SizeGigabytes: 5, DeliveredSizeGigabytes: 5},
			requested: 6,
			reads:     []*civogo.Volume{{Status: "attached", SizeGigabytes: 6, DeliveredSizeGigabytes: 6, Resize: done}},
			wantState: resizeDelivered, quick: true,
		},
		{
			name:      "no-op request on a volume that already delivers the size is judged by the resize outcome, not the size",
			before:    &civogo.Volume{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 60, Resize: oldDone},
			requested: 60,
			reads: []*civogo.Volume{
				{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 60, Resize: oldDone},
				{Status: "attached", SizeGigabytes: 60, DeliveredSizeGigabytes: 60, Resize: done},
			},
			wantState: resizeDelivered, quick: true,
		},
		{
			name:      "cancelled with no reason",
			before:    &civogo.Volume{Status: "attached", SizeGigabytes: 6, DeliveredSizeGigabytes: 6},
			requested: 10,
			reads: []*civogo.Volume{
				{Status: "attached", SizeGigabytes: 10, DeliveredSizeGigabytes: 6, Resize: cancelled},
			},
			wantState: resizeFailed, wantDetail: "the resize cancelled; the volume still delivers 6 GB of the 10 GB requested", quick: true,
		},
		{
			name:      "API without the delivered size: transitional status seen and cleared is unconfirmed",
			before:    &civogo.Volume{Status: "available", SizeGigabytes: 5},
			requested: 6,
			reads: []*civogo.Volume{
				{Status: "resizing", SizeGigabytes: 6},
				{Status: "available", SizeGigabytes: 6},
			},
			wantState: resizeUnconfirmed, wantDetail: "the resize has settled, but this API does not report the delivered size", quick: true,
		},
		{
			name:      "API without the delivered size: never seen in progress is unconfirmed after the grace window",
			before:    &civogo.Volume{Status: "available", SizeGigabytes: 5},
			requested: 6,
			reads:     []*civogo.Volume{{Status: "available", SizeGigabytes: 6}},
			wantState: resizeUnconfirmed, wantDetail: "the resize was not seen in progress within 1m0s",
		},
		{
			name:      "transient read errors are tolerated",
			before:    &civogo.Volume{Status: "attached", SizeGigabytes: 5, DeliveredSizeGigabytes: 5},
			requested: 6,
			reads: []*civogo.Volume{
				nil, nil,
				{Status: "resizing", SizeGigabytes: 6, DeliveredSizeGigabytes: 5, Resize: inProgress},
				nil,
				{Status: "attached", SizeGigabytes: 6, DeliveredSizeGigabytes: 6, Resize: done},
			},
			wantState: resizeDelivered, quick: true,
		},
		{
			name:      "persistent read errors give up",
			before:    &civogo.Volume{Status: "attached", SizeGigabytes: 5, DeliveredSizeGigabytes: 5},
			requested: 6,
			reads:     []*civogo.Volume{nil},
			wantErr:   "reading the volume failed 4 times in a row",
		},
		{
			name:      "stuck in flight times out",
			before:    &civogo.Volume{Status: "attached", SizeGigabytes: 5, DeliveredSizeGigabytes: 5},
			requested: 6,
			reads:     []*civogo.Volume{{Status: "resizing", SizeGigabytes: 6, DeliveredSizeGigabytes: 5, Resize: inProgress}},
			wantErr:   "timed out after 10m0s",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clock := &fakeClock{t: t0}
			got, err := newWait(tc.before, tc.requested, clock, tc.reads...).run()
			switch {
			case tc.wantErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)):
				t.Fatalf("error = %v, want it to contain %q", err, tc.wantErr)
			case tc.wantErr != "":
				return
			}
			if got.state != tc.wantState {
				t.Errorf("state = %q, want %q (detail %q)", got.state, tc.wantState, got.detail)
			}
			if got.detail != tc.wantDetail {
				t.Errorf("detail = %q, want %q", got.detail, tc.wantDetail)
			}
			if elapsed := clock.t.Sub(t0); tc.quick && elapsed >= testGrace {
				t.Errorf("took %s, want a settled outcome before the %s grace window", elapsed, testGrace)
			}
		})
	}
}

func TestResizeInFlight(t *testing.T) {
	cases := []struct {
		name string
		v    civogo.Volume
		want bool
	}{
		{"resizing status alone", civogo.Volume{Status: "resizing"}, true},
		{"in-progress outcome alone", civogo.Volume{Status: "attached", Resize: &civogo.VolumeResize{State: civogo.VolumeResizeInProgress}}, true},
		{"migrating alone is not a resize", civogo.Volume{Status: "migrating"}, false},
		{"settled outcome", civogo.Volume{Status: "attached", Resize: &civogo.VolumeResize{State: civogo.VolumeResizeCompleted}}, false},
		{"nothing", civogo.Volume{Status: "available"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := resizeInFlight(&tc.v); got != tc.want {
				t.Errorf("resizeInFlight = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestResizeLabels(t *testing.T) {
	failedNoReason := civogo.Volume{Resize: &civogo.VolumeResize{State: civogo.VolumeResizeFailed}}
	if got := resizeLabel(failedNoReason); got != "failed" {
		t.Errorf("label without reason = %q, want %q", got, "failed")
	}
	failed := civogo.Volume{Resize: &civogo.VolumeResize{State: civogo.VolumeResizeFailed, Reason: "Timeout"}}
	if got := resizeLabel(failed); got != "failed (Timeout)" {
		t.Errorf("label = %q", got)
	}
	if got := resizeLabel(civogo.Volume{}); got != "" {
		t.Errorf("label with no resize = %q, want empty", got)
	}
	if got := deliveredSizeLabel(civogo.Volume{DeliveredSizeGigabytes: 7}); got != "7 GB" {
		t.Errorf("delivered label = %q", got)
	}
	if reportsResize([]civogo.Volume{{SizeGigabytes: 5}, {SizeGigabytes: 6}}) {
		t.Error("an API without the fields must not add the columns")
	}
	if !reportsResize([]civogo.Volume{{SizeGigabytes: 5}, {SizeGigabytes: 6, DeliveredSizeGigabytes: 6}}) {
		t.Error("a delivered size on any volume must add the columns")
	}
}
