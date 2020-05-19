package utility

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/shiena/ansicolor"
	"os"
)

func init() {
	color.Output = ansicolor.NewAnsiColorWriter(os.Stderr)
}

// Green is the function to convert str to green in console
func Green(value string) string {
	newColor := color.New(color.FgGreen).SprintFunc()
	return newColor(value)
}

// Red is the function to convert str to red in console
func Red(value string) string {
	newColor := color.New(color.FgRed).SprintFunc()
	return newColor(value)
}

// Error is the function to handler all error in the Cli
func Error(msg string, args ...interface{}) {
	fmt.Fprintf(color.Output, "%s: %s\n", color.RedString("Error"), fmt.Sprintf(msg, args...))

}

// YellowConfirm is the function to handler all delete confirm
func YellowConfirm(msg string, args ...interface{}) {
	fmt.Fprintf(color.Output, "%s: %s", color.YellowString("Warning"), fmt.Sprintf(msg, args...))
}

// RedConfirm is the function to handler the new version of the Cli
func RedConfirm(msg string, args ...interface{}) {
	fmt.Fprintf(color.Output, "%s: %s", color.RedString("IMPORTANT"), fmt.Sprintf(msg, args...))
}
