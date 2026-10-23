package database

import (
	"strings"
	"testing"

	"github.com/civo/civogo"
)

func TestFailedVersionUpgradeChecks(t *testing.T) {
	checks := []civogo.DatabaseVersionUpgradeCheck{
		{Name: "database_ready", Passed: true},
		{Name: "storage_headroom", Passed: false, Message: "grow the storage first"},
		{Name: "extensions_available", Passed: false},
	}
	got := failedVersionUpgradeChecks(checks)
	if len(got) != 2 || got[0].Name != "storage_headroom" || got[1].Name != "extensions_available" {
		t.Errorf("unexpected failed checks %+v", got)
	}
	if got := failedVersionUpgradeChecks(nil); len(got) != 0 {
		t.Errorf("expected no failed checks, got %+v", got)
	}
}

func TestVersionUpgradeVerdict(t *testing.T) {
	plan := &civogo.UpgradeDatabaseVersionResponse{FromVersion: "17", ToVersion: "18"}
	cases := []struct {
		checks []civogo.DatabaseVersionUpgradeCheck
		want   string
	}{
		{[]civogo.DatabaseVersionUpgradeCheck{{Passed: true}}, "All checks passed: app can be upgraded from PostgreSQL 17 to 18."},
		{[]civogo.DatabaseVersionUpgradeCheck{{Passed: true}, {Passed: false}}, "1 check failed: fix it before upgrading app from PostgreSQL 17 to 18."},
		{[]civogo.DatabaseVersionUpgradeCheck{{Passed: false}, {Passed: false}}, "2 checks failed: fix them before upgrading app from PostgreSQL 17 to 18."},
	}
	for _, c := range cases {
		plan.Checks = c.checks
		if got := versionUpgradeVerdict("app", plan); got != c.want {
			t.Errorf("expected %q, got %q", c.want, got)
		}
	}
}

func TestPreUpgradeBackupName(t *testing.T) {
	if got := preUpgradeBackupName(&civogo.UpgradeDatabaseVersionResponse{ToVersion: "18"}); got != "pre-upgrade-pg18" {
		t.Errorf("expected the default name, got %q", got)
	}
	if got := preUpgradeBackupName(&civogo.UpgradeDatabaseVersionResponse{ToVersion: "18", BackupName: "custom"}); got != "custom" {
		t.Errorf("expected the reported name, got %q", got)
	}
}

func TestMajorVersionUpgradeWarning(t *testing.T) {
	got := majorVersionUpgradeWarning("17", "18", "pre-upgrade-pg18")
	for _, want := range []string{"PostgreSQL 17 to 18", "pre-upgrade-pg18 is taken first and kept", "stops for a few minutes", "replicas re-sync", "first backup on PostgreSQL 18", "no downgrade"} {
		if !strings.Contains(got, want) {
			t.Errorf("warning misses %q:\n%s", want, got)
		}
	}
}

func TestFormatVersionUpgrade(t *testing.T) {
	cases := []struct {
		in   *civogo.DatabaseVersionUpgrade
		want string
	}{
		{nil, ""},
		{&civogo.DatabaseVersionUpgrade{}, ""},
		{&civogo.DatabaseVersionUpgrade{Phase: "Upgrading", FromVersion: "16", TargetVersion: "17", BackupName: "pre-upgrade-pg17"}, "Upgrading (PostgreSQL 16 to 17, backup pre-upgrade-pg17)"},
		{&civogo.DatabaseVersionUpgrade{Phase: "Failed", FromVersion: "16", TargetVersion: "17", Message: "pg_upgrade --check failed"}, "Failed (PostgreSQL 16 to 17): pg_upgrade --check failed"},
	}
	for _, c := range cases {
		if got := formatVersionUpgrade(c.in); got != c.want {
			t.Errorf("expected %q, got %q", c.want, got)
		}
	}
}
