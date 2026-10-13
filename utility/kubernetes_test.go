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

func TestIsStandardSmallOrXSmallKubeSize(t *testing.T) {
	tests := []struct {
		name string
		size string
		want bool
	}{
		{name: "standard xsmall", size: "g4s.kube.xsmall", want: true},
		{name: "standard small", size: "g4s.kube.small", want: true},
		{name: "standard medium is not small/xsmall", size: "g4s.kube.medium", want: false},
		{name: "standard large is not small/xsmall", size: "g4s.kube.large", want: false},
		{name: "performance small is not standard tier", size: "g4p.kube.small", want: false},
		{name: "ram-optimized small is not standard tier", size: "g4m.kube.small", want: false},
		{name: "cpu-optimized small is not standard tier", size: "g4c.kube.small", want: false},
		{name: "legacy non-kube instance size", size: "g3.xsmall", want: false},
		{name: "empty size", size: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStandardSmallOrXSmallKubeSize(tt.size); got != tt.want {
				t.Errorf("IsStandardSmallOrXSmallKubeSize(%q) = %v, want %v", tt.size, got, tt.want)
			}
		})
	}
}

func TestResolveLogsCollectorEnabled(t *testing.T) {
	tests := []struct {
		name          string
		size          string
		explicitlySet bool
		value         bool
		wantEnabled   *bool
		wantMessage   string
	}{
		{
			name:          "explicit true is respected even on xsmall",
			size:          "g4s.kube.xsmall",
			explicitlySet: true,
			value:         true,
			wantEnabled:   boolPtr(true),
			wantMessage:   "",
		},
		{
			name:          "explicit false is respected even on a large node",
			size:          "g4s.kube.large",
			explicitlySet: true,
			value:         false,
			wantEnabled:   boolPtr(false),
			wantMessage:   "",
		},
		{
			name:          "auto-disabled on standard small when not explicit",
			size:          "g4s.kube.small",
			explicitlySet: false,
			value:         true,
			wantEnabled:   boolPtr(false),
			wantMessage:   LogsCollectorDisabledMessage,
		},
		{
			name:          "auto-disabled on standard xsmall when not explicit",
			size:          "g4s.kube.xsmall",
			explicitlySet: false,
			value:         true,
			wantEnabled:   boolPtr(false),
			wantMessage:   LogsCollectorDisabledMessage,
		},
		{
			name:          "left to API default on standard medium when not explicit",
			size:          "g4s.kube.medium",
			explicitlySet: false,
			value:         true,
			wantEnabled:   nil,
			wantMessage:   "",
		},
		{
			name:          "performance-tier small is not auto-disabled",
			size:          "g4p.kube.small",
			explicitlySet: false,
			value:         true,
			wantEnabled:   nil,
			wantMessage:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEnabled, gotMessage := ResolveLogsCollectorEnabled(tt.size, tt.explicitlySet, tt.value)

			if (gotEnabled == nil) != (tt.wantEnabled == nil) {
				t.Fatalf("ResolveLogsCollectorEnabled(%q, %v, %v) enabled = %v, want %v", tt.size, tt.explicitlySet, tt.value, gotEnabled, tt.wantEnabled)
			}
			if gotEnabled != nil && *gotEnabled != *tt.wantEnabled {
				t.Errorf("ResolveLogsCollectorEnabled(%q, %v, %v) enabled = %v, want %v", tt.size, tt.explicitlySet, tt.value, *gotEnabled, *tt.wantEnabled)
			}
			if gotMessage != tt.wantMessage {
				t.Errorf("ResolveLogsCollectorEnabled(%q, %v, %v) message = %q, want %q", tt.size, tt.explicitlySet, tt.value, gotMessage, tt.wantMessage)
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}
