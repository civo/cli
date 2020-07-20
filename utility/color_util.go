package utility

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/shiena/ansicolor"
)

func init() {
	color.Output = ansicolor.NewAnsiColorWriter(os.Stderr)
}

// Green is the function to convert str to green in console
func Green(value string) string {
	newColor := color.New(color.FgGreen).SprintFunc()
	return newColor(value)
}

// Yellow is the function to convert str to yellow in console
func Yellow(value string) string {
	newColor := color.New(color.FgYellow).SprintFunc()
	return newColor(value)
}

// Blue is the function to convert str to blue in console
func Blue(value string) string {
	newColor := color.New(color.FgBlue).SprintFunc()
	return newColor(value)
}

// Magenta is the function to convert str to magenta in console
func Magenta(value string) string {
	newColor := color.New(color.FgMagenta).SprintFunc()
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

// ColorStatus is to print the status of the Instance or k8s Cluster
func ColorStatus(status string) string {

	var returnText string

	if status == "ACTIVE" {
		returnText = Green(status)
	}

	if status == "SHUTOFF" {
		returnText = Red(status)
	}

	if status == "REBOOTING" {
		returnText = Yellow(status)
	}

	if status == "INSTANCE-CREATE" {
		returnText = Blue(status)
	}

	if status == "INSTALLING" {
		returnText = Magenta(status)
	}

	return returnText
}
