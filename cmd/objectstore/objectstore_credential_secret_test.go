package objectstore

import "testing"

func TestResolveAccessKey(t *testing.T) {
	tests := []struct {
		name      string
		accessKey string
		args      []string
		expected  string
	}{
		{
			name:      "no flag and no args",
			accessKey: "",
			args:      []string{},
			expected:  "",
		},
		{
			name:      "positional arg only",
			accessKey: "",
			args:      []string{"abc123"},
			expected:  "abc123",
		},
		{
			name:      "flag only",
			accessKey: "flagkey",
			args:      []string{},
			expected:  "flagkey",
		},
		{
			name:      "flag takes precedence over positional arg",
			accessKey: "flagkey",
			args:      []string{"abc123"},
			expected:  "flagkey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveAccessKey(tt.accessKey, tt.args)
			if result != tt.expected {
				t.Errorf("resolveAccessKey(%q, %v) = %q, want %q", tt.accessKey, tt.args, result, tt.expected)
			}
		})
	}
}
