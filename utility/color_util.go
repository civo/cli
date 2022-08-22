package utility

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/gookit/color"
	"github.com/savioxavier/termlink"
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

// Orange is the function to convert str to orange in console
func Orange(value string) string {
	newColor := color.FgYellow.Render
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
	res := common.VersionCheck()
	if res.Outdated {
		updateVersion := "civo update"
		gitIssueLink := termlink.ColorLink("GitHub issue", "https://github.com/civo/cli/issues", "italic green")
		fmt.Printf("Please, run %q and retry the command. If you are still facing issues, please report it on our community slack or open a %s \n", updateVersion, gitIssueLink)
	}
	fmt.Fprintf(os.Stderr, "%s: %s\n", color.Red.Sprintf("Error"), fmt.Sprintf(msg, args...))
}

// Info is the function to handler all info messages in the Cli
func Info(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", color.Blue.Sprintf("Info"), fmt.Sprintf(msg, args...))
}

// Warning is the function to handler all warnings in the Cli
func Warning(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", color.Yellow.Sprintf("Warning"), fmt.Sprintf(msg, args...))
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

	switch {
	case status == "ACTIVE":
		returnText = Green(status)
	case status == "SHUTOFF":
		returnText = Red(status)
	case status == "REBOOTING":
		returnText = Yellow(status)
	case status == "BUILDING":
		returnText = Yellow(status)
	case status == "INSTANCE-CREATE":
		returnText = Blue(status)
	case status == "INSTALLING":
		returnText = Blue(status)
	case status == "SCALING":
		returnText = Magenta(status)
	case status == "STOPPING":
		returnText = Yellow(status)
	default:
		returnText = Red("Unknown")
	}

	return returnText
}
