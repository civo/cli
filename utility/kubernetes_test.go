package utility

import (
	"slices"
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

func TestRemoveNodePool(t *testing.T) {
	pools := func(ids ...string) []civogo.KubernetesClusterPoolConfig {
		s := []civogo.KubernetesClusterPoolConfig{}
		for _, id := range ids {
			s = append(s, civogo.KubernetesClusterPoolConfig{ID: id})
		}
		return s
	}

	tests := []struct {
		name       string
		pools      []civogo.KubernetesClusterPoolConfig
		remove     []string
		wantPools  []string
		wantRemove []string
	}{
		{
			name:       "removes the first pool",
			pools:      pools("aaaaaa-1", "bbbbbb-2", "cccccc-3"),
			remove:     []string{"aaaaaa"},
			wantPools:  []string{"bbbbbb-2", "cccccc-3"},
			wantRemove: []string{"aaaaaa-1"},
		},
		{
			name:       "removes a middle pool",
			pools:      pools("aaaaaa-1", "bbbbbb-2", "cccccc-3"),
			remove:     []string{"bbbbbb"},
			wantPools:  []string{"aaaaaa-1", "cccccc-3"},
			wantRemove: []string{"bbbbbb-2"},
		},
		{
			name:       "removes the last pool",
			pools:      pools("aaaaaa-1", "bbbbbb-2", "cccccc-3"),
			remove:     []string{"cccccc"},
			wantPools:  []string{"aaaaaa-1", "bbbbbb-2"},
			wantRemove: []string{"cccccc-3"},
		},
		{
			name:       "empty slice does not panic",
			pools:      pools(),
			remove:     []string{"aaaaaa"},
			wantPools:  []string{},
			wantRemove: []string{},
		},
		{
			name:       "unknown ID leaves pools unchanged",
			pools:      pools("aaaaaa-1", "bbbbbb-2"),
			remove:     []string{"zzzzzz"},
			wantPools:  []string{"aaaaaa-1", "bbbbbb-2"},
			wantRemove: []string{},
		},
		{
			name:       "same ID twice removes only the matching pool",
			pools:      pools("aaaaaa-1", "bbbbbb-2", "cccccc-3"),
			remove:     []string{"bbbbbb", "bbbbbb"},
			wantPools:  []string{"aaaaaa-1", "cccccc-3"},
			wantRemove: []string{"bbbbbb-2"},
		},
		{
			name:       "ID fragment in the middle of another ID does not match",
			pools:      pools("aaaaaa-bbbbbb", "cccccc-1"),
			remove:     []string{"bbbbbb"},
			wantPools:  []string{"aaaaaa-bbbbbb", "cccccc-1"},
			wantRemove: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("RemoveNodePool panicked: %v", r)
				}
			}()

			got := tt.pools
			names := []string{}
			for _, id := range tt.remove {
				got, names = RemoveNodePool(got, id, names)
			}

			gotIDs := []string{}
			for _, p := range got {
				gotIDs = append(gotIDs, p.ID)
			}
			slices.Sort(gotIDs)

			if !slices.Equal(gotIDs, tt.wantPools) {
				t.Errorf("remaining pools = %v, want %v", gotIDs, tt.wantPools)
			}
			if !slices.Equal(names, tt.wantRemove) {
				t.Errorf("removed names = %v, want %v", names, tt.wantRemove)
			}
		})
	}
}

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
