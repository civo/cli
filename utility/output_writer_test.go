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
