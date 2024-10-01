package instance

import (
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/pkg/browser"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
)

var instanceVncCmd = &cobra.Command{
	Use:     "vnc",
	Example: "civo instance vnc INSTANCE-ID/NAME",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Enable and access VNC on an instance",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		// Create the API client
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Failed to connect to Civo's API: %s", err)
			os.Exit(1)
		}

		// Set the region if specified by the user
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		// Locate the instance
		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("Unable to find instance with ID/Name '%s': %s", args[0], err)
			os.Exit(1)
		}

		// Enable VNC for the instance
		vnc, err := client.GetInstanceVnc(instance.ID)
		if err != nil {
			utility.Error("Failed to enable VNC on instance '%s': %s", instance.Hostname, err)
			os.Exit(1)
		}

		// Display VNC details
		utility.Info("VNC has been successfully enabled for instance: %s", instance.Hostname)
		utility.Info("VNC URL: %s", vnc.URI)
		utility.Info("Opening VNC in your default browser...")

		// Open VNC in the browser
		err = browser.OpenInBrowser(vnc.URI)
		if err != nil {
			utility.Error("Failed to open VNC URL in the browser: %s", err)
		} else {
			utility.Info("VNC session is now active. You can access your instance's graphical interface.")
		}
	},
}
