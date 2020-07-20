package utility

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

// Green is the function to convert str to green in console
func Green(value string) string {
	newColor := color.FgGreen.Render
	return newColor(value)
}

// Yellow is the function to convert str to yellow in console
func Yellow(value string) string {
	newColor := color.New(color.FgYellow).Render
	return newColor(value)
}

// Blue is the function to convert str to blue in console
func Blue(value string) string {
	newColor := color.New(color.FgBlue).Render
	return newColor(value)
}

// Magenta is the function to convert str to magenta in console
func Magenta(value string) string {
	newColor := color.New(color.FgMagenta).Render
	return newColor(value)
}

// Red is the function to convert str to red in console
func Red(value string) string {
	newColor := color.New(color.FgRed).Render
	return newColor(value)
}

// Error is the function to handler all error in the Cli
func Error(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", color.Error.Sprintf("Error"), fmt.Sprintf(msg, args...))

}

// YellowConfirm is the function to handler all delete confirm
func YellowConfirm(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s", color.Warn.Sprintf("Warning"), fmt.Sprintf(msg, args...))
}

// RedConfirm is the function to handler the new version of the Cli
func RedConfirm(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s", color.Red.Sprintf("IMPORTANT"), fmt.Sprintf(msg, args...))
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
