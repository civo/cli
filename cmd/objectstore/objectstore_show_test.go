package objectstore

import (
	"testing"
)

func TestHumanSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{
			name:     "Zero bytes",
			bytes:    0,
			expected: "0 B",
		},
		{
			name:     "Bytes only",
			bytes:    512,
			expected: "512 B",
		},
		{
			name:     "Exactly 1 KB",
			bytes:    1024,
			expected: "1.00 KB",
		},
		{
			name:     "Kilobytes",
			bytes:    1536,
			expected: "1.50 KB",
		},
		{
			name:     "Exactly 1 MB",
			bytes:    1024 * 1024,
			expected: "1.00 MB",
		},
		{
			name:     "Megabytes",
			bytes:    1572864,
			expected: "1.50 MB",
		},
		{
			name:     "Exactly 1 GB",
			bytes:    1024 * 1024 * 1024,
			expected: "1.00 GB",
		},
		{
			name:     "Gigabytes",
			bytes:    1610612736,
			expected: "1.50 GB",
		},
		{
			name:     "Exactly 1 TB",
			bytes:    1024 * 1024 * 1024 * 1024,
			expected: "1.00 TB",
		},
		{
			name:     "Terabytes",
			bytes:    1649267441664,
			expected: "1.50 TB",
		},
		{
			name:     "Exactly 1 PB",
			bytes:    1024 * 1024 * 1024 * 1024 * 1024,
			expected: "1.00 PB",
		},
		{
			name:     "Petabytes",
			bytes:    1688849860263936,
			expected: "1.50 PB",
		},
		{
			name:     "Large number with decimals",
			bytes:    12356 * 1024,
			expected: "12.07 MB",
		},
		{
			name:     "Small KB value",
			bytes:    2048,
			expected: "2.00 KB",
		},
		{
			name:     "Fractional MB",
			bytes:    2621440,
			expected: "2.50 MB",
		},
		{
			name:     "Maximum int64 value",
			bytes:    9223372036854775807,
			expected: "8192.00 PB",
		},
		{
			name:     "Negative value (edge case)",
			bytes:    -1024,
			expected: "-1024 B",
		},
		{
			name:     "Just under 1 KB",
			bytes:    1023,
			expected: "1023 B",
		},
		{
			name:     "Just over 1 KB",
			bytes:    1025,
			expected: "1.00 KB",
		},
		{
			name:     "Just under 1 MB",
			bytes:    1048575,
			expected: "1024.00 KB",
		},
		{
			name:     "Just over 1 MB",
			bytes:    1048577,
			expected: "1.00 MB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := humanSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("humanSize(%d) = %s, want %s", tt.bytes, result, tt.expected)
			}
		})
	}
}
