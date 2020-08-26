package utility

import (
	"fmt"
	"math"
	"runtime"
)

// CheckOS is a function to check the OS of the user
func CheckOS() string {
	os := runtime.GOOS
	var returnValue string
	switch os {
	case "windows":
		returnValue = "windows"
	case "darwin":
		returnValue = "darwin"
	case "linux":
		returnValue = "linux"
	default:
		fmt.Printf("%s.\n", os)
	}

	return returnValue
}

// CheckQuotaPercent function to check the percent of the quota
func CheckQuotaPercent(limit int, usage int) string {

	var returnText string

	calculation := float64(usage) / float64(limit) * 100
	percent := math.Round(calculation)

	switch {
	case percent >= 80 && percent < 100:
		returnText = Orange(fmt.Sprintf("%d/%d", usage, limit))
	case percent == 100:
		returnText = Red(fmt.Sprintf("%d/%d", usage, limit))
	default:
		returnText = Green(fmt.Sprintf("%d/%d", usage, limit))
	}

	return returnText
}
