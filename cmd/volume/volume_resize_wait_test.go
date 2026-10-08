package volume

import (
	"strings"
	"testing"
	"time"

	"github.com/civo/civogo"
)

// fakeClock drives waitForResize without real sleeps: every sleep advances the clock.
type fakeClock struct{ t time.Time }

func (c *fakeClock) now() time.Time        { return c.t }
func (c *fakeClock) sleep(d time.Duration) { c.t = c.t.Add(d) }

func script(volumes ...*civogo.Volume) func() (*civogo.Volume, error) {
	i := 0
	return func() (*civogo.Volume, error) {
		v := volumes[i]
		if i < len(volumes)-1 {
			i++
		}
		return v, nil
	}
}

func TestWaitForResize(t *testing.T) {
	inProgress := &civogo.VolumeResize{State: civogo.VolumeResizeInProgress, Reason: "Expanding"}
	completed := &civogo.VolumeResize{State: civogo.VolumeResizeCompleted, Reason: "Completed"}
	failed := &civogo.VolumeResize{State: civogo.VolumeResizeFailed, Reason: "InsufficientStorage", Message: "pool is full"}

	cases := []struct {
		name          string
		reads         []*civogo.Volume
		wantDelivered int
		wantErr       string
		wantOutcome   string
	}{
		{
			name: "online resize seen in flight then delivered",
			reads: []*civogo.Volume{
				{Status: "attached", SizeGigabytes: 6, DeliveredSizeGigabytes: 5, Resize: completed},
				{Status: "resizing", SizeGigabytes: 6, DeliveredSizeGigabytes: 5, Resize: inProgress},
				{Status: "attached", SizeGigabytes: 6, DeliveredSizeGigabytes: 6, Resize: completed},
			},
			wantDelivered: 6,
		},
		{
			name: "refused resize settles below the request with a reason",
			reads: []*civogo.Volume{
				{Status: "resizing", SizeGigabytes: 300, DeliveredSizeGigabytes: 6, Resize: inProgress},
				{Status: "attached", SizeGigabytes: 300, DeliveredSizeGigabytes: 6, Resize: failed},
			},
			wantDelivered: 6,
			wantOutcome:   "the resize failed (InsufficientStorage): pool is full; the volume still delivers 6 GB of the 300 GB requested",
		},
		{
			name: "delivered before the first poll",
			reads: []*civogo.Volume{
				{Status: "attached", SizeGigabytes: 6, DeliveredSizeGigabytes: 6, Resize: completed},
			},
			wantDelivered: 6,
		},
		{
			name: "API without the delivered size settles after the grace window",
			reads: []*civogo.Volume{
				{Status: "attached", SizeGigabytes: 6},
			},
			wantDelivered: 0,
		},
		{
			name: "stuck in flight times out",
			reads: []*civogo.Volume{
				{Status: "resizing", SizeGigabytes: 6, DeliveredSizeGigabytes: 5, Resize: inProgress},
			},
			wantDelivered: 5,
			wantErr:       "timed out after 10m0s",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clock := &fakeClock{t: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)}
			got, err := waitForResize(script(tc.reads...), tc.reads[0].SizeGigabytes, time.Minute, 10*time.Minute, 2*time.Second, clock.sleep, clock.now)
			switch {
			case tc.wantErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)):
				t.Fatalf("error = %v, want it to contain %q", err, tc.wantErr)
			}
			if got.DeliveredSizeGigabytes != tc.wantDelivered {
				t.Errorf("delivered = %d, want %d", got.DeliveredSizeGigabytes, tc.wantDelivered)
			}
			if err == nil {
				if outcome := resizeOutcome(got, tc.reads[0].SizeGigabytes); outcome != tc.wantOutcome {
					t.Errorf("outcome = %q, want %q", outcome, tc.wantOutcome)
				}
			}
		})
	}
}

func TestResizeOutcomeWithoutDeliveredSize(t *testing.T) {
	// An API that predates delivered_size_gb reports nothing to judge: no complaint.
	if got := resizeOutcome(&civogo.Volume{Status: "attached", SizeGigabytes: 6}, 6); got != "" {
		t.Errorf("outcome = %q, want none", got)
	}
	// A settled failure is reported even when the API predates the delivered size.
	v := &civogo.Volume{Status: "attached", SizeGigabytes: 6, Resize: &civogo.VolumeResize{State: civogo.VolumeResizeCancelled, Reason: "Cancelled"}}
	if got := resizeOutcome(v, 6); got != "the resize cancelled (Cancelled)" {
		t.Errorf("outcome = %q", got)
	}
}
