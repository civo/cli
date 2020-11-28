package utility

import "testing"

func TestWriteCustomOutput(t *testing.T) {
	ow := NewOutputWriter()
	ow.StartLine()
	ow.AppendData("A", "a")
	ow.AppendData("B", "b")
	ow.WriteCustomOutput("B")
}

func ExampleFWriteCustomOutput() {
	ow := NewOutputWriter()
	ow.StartLine()
	ow.AppendData("ID", "1")
	ow.AppendData("Key", "Raspberry")
	ow.AppendData("Desc", "first")
	ow.StartLine()
	ow.AppendData("ID", "2")
	ow.AppendData("Key", "Pi")
	ow.AppendData("Desc", "second")
	ow.WriteCustomOutput("Key,Desc")

	// Output:
	// Raspberry,first
	// Pi,second
}
