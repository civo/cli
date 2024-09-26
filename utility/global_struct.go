package utility

import (
	"fmt"
	"os/exec"
	"runtime"
)

// ObjecteList contains ID and Name
type ObjecteList struct {
	ID, Name string
}

// OpenInBrowser attempts to open the specified URL in the default browser.
// Returns an error if the command fails or the platform is unsupported.
func OpenInBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		return fmt.Errorf("failed to open URL: %w", err)
	}
	return nil
}
