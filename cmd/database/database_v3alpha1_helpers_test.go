package database

import (
	"testing"
	"time"
)

func TestParseTargetTime(t *testing.T) {
	got, err := parseTargetTime("2026-10-21T12:00:00+02:00")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	want := time.Date(2026, 10, 21, 10, 0, 0, 0, time.UTC)
	if !got.Equal(want) || got.Location() != time.UTC {
		t.Errorf("expected %s in UTC, got %s", want, got)
	}
	if _, err := parseTargetTime("yesterday"); err == nil {
		t.Errorf("expected an error for a non RFC 3339 value")
	}
}

func TestParseExtensions(t *testing.T) {
	got := parseExtensions(" pg_cron, postgis ,,")
	if len(got) != 2 || got[0] != "pg_cron" || got[1] != "postgis" {
		t.Errorf("unexpected extensions %v", got)
	}
	if got := parseExtensions(""); len(got) != 0 {
		t.Errorf("expected no extensions, got %v", got)
	}
}

func TestParseDatabaseParameters(t *testing.T) {
	got, err := parseDatabaseParameters([]string{"max_connections=200", "shared_buffers = 1GB"})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if got["max_connections"] != "200" || got["shared_buffers"] != "1GB" {
		t.Errorf("unexpected parameters %v", got)
	}
	for _, bad := range []string{"max_connections", "=200"} {
		if _, err := parseDatabaseParameters([]string{bad}); err == nil {
			t.Errorf("expected an error for %q", bad)
		}
	}
}

func TestBackupDestinationFromFlags(t *testing.T) {
	reset := func() {
		destinationRegion, s3Endpoint, s3Region, s3Bucket, s3AccessKey, s3SecretKey, s3PathPrefix = "", "", "", "", "", "", ""
	}

	reset()
	if d, err := backupDestinationFromFlags(); err != nil || d != nil {
		t.Errorf("expected no destination without flags, got %v, %v", d, err)
	}

	reset()
	destinationRegion = "fra1"
	d, err := backupDestinationFromFlags()
	if err != nil || d == nil || d.ObjectStore == nil || d.ObjectStore.Region != "FRA1" || d.ExternalS3 != nil {
		t.Errorf("expected an object store destination in FRA1, got %+v, %v", d, err)
	}

	reset()
	s3Endpoint, s3Bucket, s3AccessKey, s3SecretKey, s3PathPrefix = "https://s3.example.com", "bk", "AK", "SK", "db/"
	d, err = backupDestinationFromFlags()
	if err != nil || d == nil || d.ExternalS3 == nil || d.ExternalS3.Bucket != "bk" || d.ExternalS3.PathPrefix != "db/" || d.ObjectStore != nil {
		t.Errorf("expected an external S3 destination, got %+v, %v", d, err)
	}

	reset()
	s3Endpoint = "https://s3.example.com"
	if _, err := backupDestinationFromFlags(); err == nil {
		t.Errorf("expected an error when the bucket is missing")
	}

	reset()
	destinationRegion, s3Endpoint, s3Bucket = "FRA1", "https://s3.example.com", "bk"
	if _, err := backupDestinationFromFlags(); err == nil {
		t.Errorf("expected an error when both destination kinds are given")
	}
	reset()
}

func TestFormatDatabaseStorage(t *testing.T) {
	cases := []struct {
		storage, base int
		want          string
	}{
		{0, 0, ""},
		{40, 40, "40 GB"},
		{40, 0, "40 GB"},
		{120, 40, "120 GB (40 GB plan + 80 GB add-on)"},
	}
	for _, c := range cases {
		if got := formatDatabaseStorage(c.storage, c.base); got != c.want {
			t.Errorf("%d/%d: expected %q, got %q", c.storage, c.base, c.want, got)
		}
	}
}

func TestFormatRestoreWindow(t *testing.T) {
	if got := formatRestoreWindow(nil); got != "not available yet (waiting for the first full backup)" {
		t.Errorf("unexpected text for nil: %q", got)
	}
	ts := time.Date(2026, 10, 21, 0, 0, 0, 0, time.FixedZone("CEST", 2*3600))
	if got := formatRestoreWindow(&ts); got != "2026-10-20T22:00:00Z" {
		t.Errorf("expected UTC RFC 3339, got %q", got)
	}
}
