package volume

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
)

func TestVolumeListLookups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v2/volumes":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[
				{
					"id": "vol-1",
					"name": "vol-a-unmatched-plain",
					"network_id": "net-missing",
					"cluster_id": "",
					"instance_id": "inst-missing",
					"size_gb": 10,
					"mountpoint": "/mnt/vol1",
					"status": "attached"
				},
				{
					"id": "vol-2",
					"name": "vol-b-unmatched-cluster-inst",
					"network_id": "net-1",
					"cluster_id": "cluster-1",
					"instance_id": "inst-cluster-missing",
					"size_gb": 20,
					"mountpoint": "/mnt/vol2",
					"status": "attached"
				},
				{
					"id": "vol-3",
					"name": "vol-c-matched",
					"network_id": "net-1",
					"cluster_id": "",
					"instance_id": "inst-1",
					"size_gb": 30,
					"mountpoint": "/mnt/vol3",
					"status": "attached"
				}
			]`))
		case "/v2/networks":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[
				{"id": "net-1", "name": "net-1", "label": "network-one"},
				{"id": "net-2", "name": "net-2", "label": "network-two-last"}
			]`))
		case "/v2/instances":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"page": 1, "per_page": 20, "pages": 1, "items": [
				{"id": "inst-1", "hostname": "host-one"},
				{"id": "inst-2", "hostname": "host-two-last"}
			]}`))
		case "/v2/kubernetes/clusters":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"page": 1, "per_page": 20, "pages": 1, "items": [
				{
					"id": "cluster-1",
					"name": "cluster-one",
					"instances": [
						{"id": "inst-cluster-1", "hostname": "k8s-node-1"}
					]
				}
			]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	config.Current.Meta.URL = server.URL
	config.Current.Meta.CurrentAPIKey = "test-key"
	config.Current.APIKeys = map[string]string{"test-key": "dummy-token"}
	config.Current.Meta.DefaultRegion = "LON1"

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	common.OutputFormat = "custom"
	common.OutputFields = "id,name,network_id,cluster_id,instance_id,status"

	volumeListCmd.Run(volumeListCmd, []string{})

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected 3 lines of output, got %d:\n%s", len(lines), output)
	}

	// Line 1: vol-a-unmatched-plain
	if !strings.Contains(lines[0], "net-missing") {
		t.Errorf("expected row 1 network_id to fall back to 'net-missing', line: %q", lines[0])
	}
	if strings.Contains(lines[0], "network-two-last") {
		t.Errorf("row 1 incorrectly fell back to last network label: %q", lines[0])
	}
	if !strings.Contains(lines[0], "inst-missing") {
		t.Errorf("expected row 1 instance_id to fall back to 'inst-missing', line: %q", lines[0])
	}
	if strings.Contains(lines[0], "host-two-last") {
		t.Errorf("row 1 incorrectly fell back to last instance hostname: %q", lines[0])
	}

	// Line 2: vol-b-unmatched-cluster-inst
	if !strings.Contains(lines[1], "network-one") {
		t.Errorf("expected row 2 network_id to be 'network-one', line: %q", lines[1])
	}
	if !strings.Contains(lines[1], "cluster-1") {
		t.Errorf("expected row 2 cluster_id to be 'cluster-1', line: %q", lines[1])
	}
	if !strings.Contains(lines[1], "inst-cluster-missing") {
		t.Errorf("expected row 2 instance_id to be 'inst-cluster-missing', line: %q", lines[1])
	}

	// Line 3: vol-c-matched
	if !strings.Contains(lines[2], "network-one") {
		t.Errorf("expected row 3 network_id to be 'network-one', line: %q", lines[2])
	}
	if !strings.Contains(lines[2], "host-one") {
		t.Errorf("expected row 3 instance_id to be 'host-one', line: %q", lines[2])
	}
}
