package utility

import (
	"github.com/logrusorgru/aurora"
)

func ColorStatus(status string) interface{} {

	var returnText interface{}

	if status == "ACTIVE" {
		returnText = aurora.BrightGreen(status)
	}

	if status == "SHUTOFF" {
		returnText = aurora.BrightRed(status)
	}

	if status == "REBOOTING" {
		returnText = aurora.BrightYellow(status)
	}

	if status == "INSTANCE-CREATE" {
		returnText = aurora.BrightBlue(status)
	}

	if status == "INSTALLING" {
		returnText = aurora.BrightMagenta(status)
	}

	return returnText
}
