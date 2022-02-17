package cmd

import (
	"reflect"
	"testing"
)

type InstallApplicationsArgs struct {
	defaultApps []string
	apps        string
	removeApps  string
}

func TestInstallApplications(t *testing.T) {
	var tests = []struct {
		name string
		args InstallApplicationsArgs
		want []string
	}{
		{
			name: "Test Removal of Default Apps",
			args: InstallApplicationsArgs{defaultApps: []string{"metrics-server", "traefik"}, apps: "", removeApps: "traefik"},
			want: []string{"metrics-server"},
		},
		{
			name: "Test Installation of Default Apps and other Apps",
			args: InstallApplicationsArgs{defaultApps: []string{"metrics-server", "traefik"}, apps: "foo", removeApps: ""},
			want: []string{"metrics-server", "traefik", "foo"},
		},
		{
			name: "Test Removal of Default Apps and Installation of other Apps",
			args: InstallApplicationsArgs{defaultApps: []string{"metrics-server", "traefik"}, apps: "foo", removeApps: "traefik"},
			want: []string{"metrics-server", "foo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InstallApps(tt.args.defaultApps, tt.args.apps, tt.args.removeApps)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstallApps() = %v, want %v", got, tt.want)
			}
		})
	}
}
