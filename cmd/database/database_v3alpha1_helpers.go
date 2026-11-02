package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/civo/civogo"
)

// Flags shared by the commands that talk to the v3alpha1 database API.
var (
	// update / fork
	storageGB  int
	extensions string
	parameters []string

	// restore / fork
	targetTime string

	// fork
	forkName, forkSize, forkNetworkID, forkFirewallID string
	forkNodes                                         int
	forkWait                                          bool

	// backup destination
	destinationRegion                                                      string
	s3Endpoint, s3Region, s3Bucket, s3AccessKey, s3SecretKey, s3PathPrefix string
)

// parseTargetTime parses an RFC 3339 timestamp ("2026-10-21T10:00:00Z" or
// with an offset) and returns it in UTC, which is what the API expects.
func parseTargetTime(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time %q, use RFC 3339 such as 2026-10-21T10:00:00Z", value)
	}
	return t.UTC(), nil
}

// parseExtensions splits a comma separated list, trimming blanks and
// dropping empties, so "pg_cron, postgis" and "" both parse.
func parseExtensions(value string) []string {
	out := []string{}
	for _, e := range strings.Split(value, ",") {
		if e = strings.TrimSpace(e); e != "" {
			out = append(out, e)
		}
	}
	return out
}

// parseDatabaseParameters turns repeated key=value flags into the map the
// API takes.
func parseDatabaseParameters(pairs []string) (map[string]string, error) {
	out := map[string]string{}
	for _, pair := range pairs {
		key, value, found := strings.Cut(pair, "=")
		key = strings.TrimSpace(key)
		if !found || key == "" {
			return nil, fmt.Errorf("invalid parameter %q, use key=value", pair)
		}
		out[key] = strings.TrimSpace(value)
	}
	return out, nil
}

// backupDestinationFromFlags builds the backup destination from the
// --destination-region and --s3-* flags, or returns nil when none is set so
// the backup lands in this region's Civo Object Store as before.
func backupDestinationFromFlags() (*civogo.DatabaseBackupDestination, error) {
	s3Set := s3Endpoint != "" || s3Region != "" || s3Bucket != "" || s3AccessKey != "" || s3SecretKey != "" || s3PathPrefix != ""
	switch {
	case destinationRegion != "" && s3Set:
		return nil, fmt.Errorf("use either --destination-region or the --s3-* flags, not both")
	case destinationRegion != "":
		return &civogo.DatabaseBackupDestination{
			ObjectStore: &civogo.DatabaseBackupObjectStoreDestination{Region: strings.ToUpper(destinationRegion)},
		}, nil
	case s3Set:
		if s3Endpoint == "" || s3Bucket == "" {
			return nil, fmt.Errorf("--s3-endpoint and --s3-bucket are required for an external S3 destination")
		}
		return &civogo.DatabaseBackupDestination{
			ExternalS3: &civogo.DatabaseBackupExternalS3Destination{
				Endpoint:        s3Endpoint,
				Region:          s3Region,
				Bucket:          s3Bucket,
				AccessKeyID:     s3AccessKey,
				SecretAccessKey: s3SecretKey,
				PathPrefix:      s3PathPrefix,
			},
		}, nil
	default:
		return nil, nil
	}
}

// formatDatabaseStorage renders the per-node storage with the plan/add-on
// split when the database has grown beyond its plan.
func formatDatabaseStorage(storageGB, baseStorageGB int) string {
	if storageGB <= 0 {
		return ""
	}
	if baseStorageGB > 0 && storageGB > baseStorageGB {
		return fmt.Sprintf("%d GB (%d GB plan + %d GB add-on)", storageGB, baseStorageGB, storageGB-baseStorageGB)
	}
	return fmt.Sprintf("%d GB", storageGB)
}

// formatRestoreWindow renders the earliest restorable moment, or a hint
// while the first full backup has not completed yet.
func formatRestoreWindow(earliest *time.Time) string {
	if earliest == nil {
		return "not available yet (waiting for the first full backup)"
	}
	return earliest.UTC().Format(time.RFC3339)
}
