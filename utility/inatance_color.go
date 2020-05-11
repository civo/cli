package utility

import (
	"github.com/fatih/color"
)

func ColorStatus(status string) interface{} {

	var returnText interface{}

	if status == "ACTIVE" {
		returnText = color.GreenString(status)
	}

	if status == "SHUTOFF" {
		returnText = color.RedString(status)
	}

	if status == "REBOOTING" {
		returnText = color.YellowString(status)
	}

	if status == "INSTANCE-CREATE" {
		returnText = color.BlueString(status)
	}

	if status == "INSTALLING" {
		returnText = color.MagentaString(status)
	}

	return returnText
}
