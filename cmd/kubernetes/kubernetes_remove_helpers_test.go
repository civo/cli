package kubernetes

import (
	"strings"
	"testing"
)

func TestShouldRemoveKubeconfigContext(t *testing.T) {
	if !shouldRemoveKubeconfigContext(false) {
		t.Fatal("expected kubeconfig context removal by default")
	}
	if shouldRemoveKubeconfigContext(true) {
		t.Fatal("expected --keep-kubeconfig to retain the context")
	}
}

func TestShouldDeleteFirewall(t *testing.T) {
	if !shouldDeleteFirewall(false) {
		t.Fatal("expected firewall deletion by default")
	}
	if shouldDeleteFirewall(true) {
		t.Fatal("expected --keep-firewalls to retain the firewall")
	}
}

func TestVolumeNames(t *testing.T) {
	got := volumeNames([]pendingVolume{
		{ID: "vol-1", Name: "data"},
		{ID: "vol-2", Name: ""},
	})
	want := []string{"data", "vol-2"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("volumeNames() = %v, want %v", got, want)
	}
}

func TestLeftBehindVolumesInfo(t *testing.T) {
	if got := leftBehindVolumesInfo("demo", nil); got != "" {
		t.Fatalf("empty volumes should produce no message, got %q", got)
	}

	got := leftBehindVolumesInfo("demo", []pendingVolume{
		{ID: "abc", Name: "pvc-data"},
	})
	for _, want := range []string{
		"Volumes left behind for Kubernetes cluster demo:",
		"- pvc-data (abc)",
		"--delete-volumes",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("message missing %q:\n%s", want, got)
		}
	}
}

func TestFirewallSafeToDelete(t *testing.T) {
	if !firewallSafeToDelete(1, 0, 0) {
		t.Fatal("exclusive cluster firewall should be deletable")
	}
	if !firewallSafeToDelete(0, 0, 0) {
		t.Fatal("firewall with no other attachments should be deletable")
	}
	if firewallSafeToDelete(2, 0, 0) {
		t.Fatal("firewall shared by two clusters must not be auto-deleted")
	}
	if firewallSafeToDelete(1, 1, 0) {
		t.Fatal("firewall still attached to an instance must not be auto-deleted")
	}
	if firewallSafeToDelete(1, 0, 1) {
		t.Fatal("firewall still attached to a load balancer must not be auto-deleted")
	}
}
