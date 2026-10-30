package database

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/civo/cli/common"

	"github.com/civo/civogo"
)

func TestFailedVersionUpgradeChecks(t *testing.T) {
	checks := []civogo.DatabaseVersionUpgradeCheck{
		{Name: "database_ready", Passed: true},
		{Name: "engine_supports_target", Passed: false, Message: "upgrade the database engine first"},
		{Name: "extensions_available", Passed: false},
	}
	got := failedVersionUpgradeChecks(checks)
	if len(got) != 2 || got[0].Name != "engine_supports_target" || got[1].Name != "extensions_available" {
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

func TestVersionUpgradeSummary(t *testing.T) {
	got := versionUpgradeSummary("app", &civogo.UpgradeDatabaseVersionResponse{FromVersion: "17", ToVersion: "18", Started: true, BackupName: "pre-upgrade-pg18-20261029120000"})
	if !strings.Contains(got, "pre-upgrade-pg18-20261029120000") {
		t.Errorf("expected the backup name the API reported, got %q", got)
	}
	got = versionUpgradeSummary("app", &civogo.UpgradeDatabaseVersionResponse{FromVersion: "17", ToVersion: "18", Started: true})
	if strings.Contains(got, "pre-upgrade") || !strings.Contains(got, "A backup is taken first") {
		t.Errorf("expected no guessed backup name, got %q", got)
	}
	got = versionUpgradeSummary("app", &civogo.UpgradeDatabaseVersionResponse{FromVersion: "17", ToVersion: "18"})
	if !strings.Contains(got, "did not start") {
		t.Errorf("expected a not-started message, got %q", got)
	}
}

func TestMajorVersionUpgradeWarning(t *testing.T) {
	got := majorVersionUpgradeWarning("17", "18")
	for _, want := range []string{"PostgreSQL 17 to 18", "a backup is taken first and kept", "stops for a few minutes", "replicas re-sync", "first backup on PostgreSQL 18", "no downgrade"} {
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

// upgradeVersionServer fakes the upgrade-version endpoint: validate-only
// requests get checks, real requests start the upgrade. It records every
// request body.
func upgradeVersionServer(t *testing.T, checks []civogo.DatabaseVersionUpgradeCheck) (*civogo.Client, *[]civogo.UpgradeDatabaseVersionRequest) {
	t.Helper()
	var requests []civogo.UpgradeDatabaseVersionRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3alpha1/databases/db-1/upgrade-version" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var req civogo.UpgradeDatabaseVersionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decoding the request: %v", err)
		}
		requests = append(requests, req)
		resp := map[string]any{"fromVersion": "17", "toVersion": "18", "checks": checks}
		if !req.ValidateOnly {
			resp["started"] = true
			resp["backupName"] = "pre-upgrade-pg18"
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)
	client, err := civogo.NewClientForTestingWithServer(srv)
	if err != nil {
		t.Fatal(err)
	}
	return client, &requests
}

func TestMajorVersionUpgrade(t *testing.T) {
	outputFormat := common.OutputFormat
	common.OutputFormat = "human"
	t.Cleanup(func() { common.OutputFormat = outputFormat })

	passed := []civogo.DatabaseVersionUpgradeCheck{{Name: "database_ready", Passed: true}}
	failed := []civogo.DatabaseVersionUpgradeCheck{{Name: "database_ready", Passed: true}, {Name: "extensions_available", Passed: false}}
	db := &civogo.Database{ID: "db-1", Name: "app"}

	cases := []struct {
		name         string
		checks       []civogo.DatabaseVersionUpgradeCheck
		dryRun       bool
		confirm      bool
		wantErr      error
		wantRequests int
		wantAsked    bool
	}{
		{name: "failed checks never send the upgrade", checks: failed, confirm: true, wantErr: errVersionUpgradeChecksFailed, wantRequests: 1},
		{name: "a dry run with failed checks fails", checks: failed, dryRun: true, confirm: true, wantErr: errVersionUpgradeChecksFailed, wantRequests: 1},
		{name: "a dry run never sends the upgrade", checks: passed, dryRun: true, confirm: true, wantRequests: 1},
		{name: "declining never sends the upgrade", checks: passed, confirm: false, wantRequests: 1, wantAsked: true},
		{name: "confirming sends the upgrade", checks: passed, confirm: true, wantRequests: 2, wantAsked: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client, requests := upgradeVersionServer(t, c.checks)
			asked := false
			confirm := func(*civogo.Database, *civogo.UpgradeDatabaseVersionResponse) bool {
				asked = true
				return c.confirm
			}
			err := majorVersionUpgrade(client, db, " 18.2 ", c.dryRun, confirm)
			if !errors.Is(err, c.wantErr) {
				t.Fatalf("expected error %v, got %v", c.wantErr, err)
			}
			if len(*requests) != c.wantRequests {
				t.Fatalf("expected %d requests, got %d: %+v", c.wantRequests, len(*requests), *requests)
			}
			if asked != c.wantAsked {
				t.Errorf("expected asked=%t, got %t", c.wantAsked, asked)
			}
			if first := (*requests)[0]; !first.ValidateOnly || first.TargetVersion != "18.2" {
				t.Errorf("expected a validate-only check of the flag first, got %+v", first)
			}
			if c.wantRequests == 2 {
				if real := (*requests)[1]; real.ValidateOnly || real.TargetVersion != "18" {
					t.Errorf("expected the real request for the major the API planned, got %+v", real)
				}
			}
		})
	}
}
