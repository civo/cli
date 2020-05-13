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

func Green(value string) string {
	color := color.New(color.FgGreen).SprintFunc()
	return color(value)
}

func Red(value string) string {
	color := color.New(color.FgRed).SprintFunc()
	return color(value)
}

func Error(msg string, args ...interface{}) {
	fmt.Fprintf(color.Output, "%s: %s\n", color.RedString("Error"), fmt.Sprintf(msg, args...))
}

func YellowConfirm(msg string, args ...interface{}) {
	fmt.Fprintf(color.Output, "%s: %s", color.YellowString("Warning"), fmt.Sprintf(msg, args...))
}

func RedConfirm(msg string, args ...interface{}) {
	fmt.Fprintf(color.Output, "%s: %s", color.RedString("IMPORTANT"), fmt.Sprintf(msg, args...))
}
