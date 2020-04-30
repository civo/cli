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

	return returnText
}
