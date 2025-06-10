package instance

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/pkg/browser"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

// maxAttempts represents the max number of attempts (each one every 7s) to connect to the console URL
const maxAttempts = 5

var duration string

var instanceConsoleCmd = &cobra.Command{
	Use:     "console",
	Aliases: []string{"connect"},
	Example: "civo instance console INSTANCE-ID/NAME [--duration 2h]",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Connect to the console of an instance",
	Long: `Enable and access the console (through the default browser) on an instance with optional duration.
Duration follows Go's duration format (e.g. "30m", "1h", "24h")`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Failed to connect to Civo's API: %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("Unable to find instance with ID/Name '%s': %s", args[0], err)
			os.Exit(1)
		}

		var vnc civogo.InstanceVnc
		if duration != "" {
			vnc, err = client.GetInstanceVnc(instance.ID, duration)
		} else {
			vnc, err = client.GetInstanceVnc(instance.ID)
		}
		if err != nil {
			utility.Error("Failed to enable console access on instance '%s': %s", instance.Hostname, err)
			os.Exit(1)
		}

		utility.Info("Console access successfully enabled for instance: %s", instance.Hostname)
		utility.Info("Console access URL: %s", vnc.URI)
		utility.Info("We're preparing console access. This may take a while...")

		exchangeTokenResp, err := client.ExchangeAuthToken(&civogo.ExchangeAuthTokenRequest{})
		if err != nil {
			utility.Error("Failed to exchange your apikey with a valid Civo JWT '%s': %s", instance.Hostname, err)
			os.Exit(1)
		}

		vnc.URI = fmt.Sprintf("%s&token=%s", vnc.URI, exchangeTokenResp.AccessToken)

		err = waitEndpointReady(vnc.URI)
		if err != nil {
			utility.Error("The console URL is not reachable: %s", err)
			os.Exit(1)
		}

		utility.Info("Opening the console in your default browser...")
		time.Sleep(3 * time.Second)

		err = browser.OpenInBrowser(vnc.URI)
		if err != nil {
			utility.Error("Failed to open the console access URL in the browser: %s", err)
		} else {
			utility.Info("The console access session is now active. You can access your instance's graphical interface.")
		}
	},
}

func endpointReady(url string) bool {
	utility.Info("New attempt to reach the console URL...")
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func waitEndpointReady(url string) error {
	var attempt int
	for {
		attempt++
		if endpointReady(url) {
			return nil
		}
		if attempt == maxAttempts {
			return fmt.Errorf("max num of attempts reached: console endpoint not ready - please contact Civo support")
		}
		time.Sleep(7 * time.Second) // Wait for 7 seconds before the next attempt
	}
}
