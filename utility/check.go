package utility

import (
	"fmt"
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
