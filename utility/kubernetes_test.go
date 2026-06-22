package utility

import (
	"testing"

	"github.com/civo/civogo"
)

func TestRemoveApplicationFromInstalledListSimpleName(t *testing.T) {
	current := []civogo.KubernetesInstalledApplication{
		{
			Name: "mysql",
		},
	}
	uninstall := "mysql"

	ret := RemoveApplicationFromInstalledList(current, uninstall)

	if ret != "" {
		t.Errorf("expected '', got '%s'", ret)
	}
}

func TestRemoveApplicationFromInstalledListMissing(t *testing.T) {
	current := []civogo.KubernetesInstalledApplication{
		{
			Name: "mysql",
		},
	}
	uninstall := "postgresql"

	ret := RemoveApplicationFromInstalledList(current, uninstall)

	if ret != "mysql" {
		t.Errorf("expected 'mysql', got '%s'", ret)
	}
}

func TestRemoveApplicationFromInstalledListWithMultiple(t *testing.T) {
	current := []civogo.KubernetesInstalledApplication{
		{
			Name: "mysql",
		},
		{
			Name: "postgresql",
		},
		{
			Name: "redis",
		},
	}
	uninstall := "postgresql,mysql"

	ret := RemoveApplicationFromInstalledList(current, uninstall)

	if ret != "redis" {
		t.Errorf("expected 'redis', got '%s'", ret)
	}
}

// func TestRemoveApplicationFromInstalledListWithPlan(t *testing.T) {
// 	current := []civogo.KubernetesInstalledApplication{
// 		{
// 			Name: "mysql",
// 		},
// 	}
// 	uninstall := "mysql"

// 	ret := RemoveApplicationFromInstalledList(current, uninstall)

// 	if ret != "mysql" {
// 		t.Errorf("expected 'mysql', got '%s'", ret)
// 	}
// }

func TestIsAppCompatibleWithClusterType(t *testing.T) {
	tests := []struct {
		name         string
		requestedApp string
		clusterType  string
		want         bool
	}{
		{
			name:         "metrics-server is not compatible with talos",
			requestedApp: "metrics-server",
			clusterType:  "talos",
			want:         false,
		},
		{
			name:         "metrics-server is compatible with k3s",
			requestedApp: "metrics-server",
			clusterType:  "k3s",
			want:         true,
		},
		{
			name:         "other app is compatible with talos",
			requestedApp: "traefik2-nodeport",
			clusterType:  "talos",
			want:         true,
		},
		{
			name:         "other app is compatible with k3s",
			requestedApp: "traefik2-nodeport",
			clusterType:  "k3s",
			want:         true,
		},
		{
			name:         "unknown cluster type has no restrictions",
			requestedApp: "metrics-server",
			clusterType:  "",
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAppCompatibleWithClusterType(tt.requestedApp, tt.clusterType); got != tt.want {
				t.Errorf("isAppCompatibleWithClusterType(%q, %q) = %v, want %v", tt.requestedApp, tt.clusterType, got, tt.want)
			}
		})
	}
}
