package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"net"
	"os"

	"github.com/spf13/cobra"
)

const (
	DNSName string = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
)

var rxDNSName = regexp.MustCompile(DNSName)

var appDomainAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create", "new"},
	Short:   "Add a new domain for your application.",
	Args:    cobra.MinimumNArgs(2),
	Example: "civo app domain add APP-NAME your-app-name.example.com",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		findApp, err := client.FindApplication(args[0])
		if err != nil {
			utility.Error("App %s", err)
			os.Exit(1)
		}

		if !IsDNSName(args[1]) {
			utility.Error("%s is not a valid domain name", args[1])
			os.Exit(1)
		}

		application := &civogo.UpdateApplicationRequest{
			Domains: append(findApp.Domains, args[1]),
		}

		app, err := client.UpdateApplication(findApp.ID, application)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": app.ID, "name": app.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("\nYour application %s is now available at:\n %s", utility.Green(app.Name), utility.Green(strings.Join(app.Domains, ", ")))
		}
	},
}

// IsDNSName will validate the given string as a DNS name
func IsDNSName(str string) bool {
	if str == "" || len(strings.Replace(str, ".", "", -1)) > 255 {
		// constraints already violated
		return false
	}
	return !IsIP(str) && rxDNSName.MatchString(str)
}

// IsIP checks if a string is either IP version 4 or 6. Alias for `net.ParseIP`
func IsIP(str string) bool {
	return net.ParseIP(str) != nil
}
