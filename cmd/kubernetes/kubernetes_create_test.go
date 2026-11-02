package kubernetes

import "testing"

func TestKubernetesCreateArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{name: "name is optional", args: []string{}},
		{name: "one name", args: []string{"my-cluster"}},
		{name: "extra names are rejected", args: []string{"test2", "test3", "test4"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := kubernetesCreateCmd.ValidateArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateArgs(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}
		})
	}
}
