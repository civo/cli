package utility

import (
	"reflect"
	"testing"
)

func TestOutputWriter_WriteCustomOutput(t *testing.T) {
	ow := NewOutputWriter()
	ow.StartLine()
	ow.AppendData("A", "a")
	ow.AppendData("B", "b")
	ow.WriteCustomOutput("B")
}

// ExampleOutputWriter_WriteCustomOutput shows how to use the OutputWriter.WriteCustomOutput function
func ExampleOutputWriter_WriteCustomOutput() {
	ow := NewOutputWriter()
	// Write 3 lines and assert correct result
	ow.StartLine()
	ow.AppendData("ID", "1")
	ow.AppendData("Key", "Raspberry")
	ow.AppendData("Desc", "first")
	ow.StartLine()
	ow.AppendData("ID", "2")
	ow.AppendData("Key", "Pi")
	ow.AppendData("Desc", "second")
	ow.StartLine()
	ow.AppendData("ID", "3")
	ow.AppendData("Key", "Zero")
	ow.AppendData("Desc", "third")
	ow.WriteCustomOutput("Key,Desc")

	// Output:
	// Raspberry,first
	// Pi,second
	// Zero,third
}

func TestInvalidCustomOutputFields(t *testing.T) {
	// keys + labels as WriteCustomOutput passes them (instance ls style)
	names := []string{"id", "hostname", "public_ip", "ID", "Hostname", "Public IP"}

	tests := []struct {
		name  string
		fields string
		names []string
		want  []string
	}{
		{name: "empty fields", fields: "", names: names, want: nil},
		{name: "whitespace fields", fields: "  ", names: names, want: nil},
		{name: "exact key", fields: "hostname", names: names, want: nil},
		{name: "comma-separated keys", fields: "id,hostname", names: names, want: nil},
		{name: "spaces around keys", fields: " id , hostname ", names: names, want: nil},
		{name: "documented key template", fields: "id: hostname", names: names, want: nil},
		{name: "documented label template", fields: "ID: Hostname (Public IP)", names: names, want: nil},
		{name: "label alone", fields: "ID", names: names, want: nil},
		{name: "unknown single field", fields: "lefufu", names: names, want: []string{"lefufu"}},
		{name: "unknown with commas", fields: "lefufu,bogus", names: names, want: []string{"lefufu,bogus"}},
		{name: "mixed list still has a real key", fields: "id,lefufu,hostname", names: names, want: nil},
		{name: "no available names", fields: "id", names: nil, want: []string{"id"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InvalidCustomOutputFields(tt.fields, tt.names)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvalidCustomOutputFields(%q, %v) = %v, want %v", tt.fields, tt.names, got, tt.want)
			}
		})
	}
}
