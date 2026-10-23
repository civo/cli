package database

import (
	"strings"
	"testing"

	"github.com/civo/civogo"
)

func TestCheckUpgradeFlags(t *testing.T) {
	cases := []struct {
		engine, postgres string
		minor, wantErr   bool
	}{
		{"", "", false, false},
		{"", "", true, false},
		{"2.9.0", "", false, false},
		{"", "18", false, false},
		{"2.9.0", "", true, true},
		{"", "18", true, true},
	}
	for _, c := range cases {
		err := checkUpgradeFlags(c.engine, c.postgres, c.minor)
		if (err != nil) != c.wantErr {
			t.Errorf("engine=%q postgres=%q minor=%v: err=%v, wantErr=%v", c.engine, c.postgres, c.minor, err, c.wantErr)
		}
	}
}

func TestMinorVersionUpgradeSummary(t *testing.T) {
	plan := &civogo.UpgradeDatabaseMinorVersionResponse{FromVersion: "17.9", ToVersion: "17.10"}
	if got := minorVersionUpgradeSummary("app", plan); got != "app can be upgraded from PostgreSQL 17.9 to 17.10." {
		t.Errorf("dry run: %q", got)
	}
	plan.Started = true
	if got := minorVersionUpgradeSummary("app", plan); !strings.Contains(got, "civo db show app") || !strings.HasPrefix(got, "Upgrading app from PostgreSQL 17.9 to 17.10") {
		t.Errorf("started: %q", got)
	}
}

func TestMinorVersionUpgradeWarning(t *testing.T) {
	w := minorVersionUpgradeWarning("17.9", "17.10")
	for _, want := range []string{"17.9 to 17.10", "switchover", "point-in-time restore are unaffected", "no downgrade"} {
		if !strings.Contains(w, want) {
			t.Errorf("warning misses %q:\n%s", want, w)
		}
	}
}
