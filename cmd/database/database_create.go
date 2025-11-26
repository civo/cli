package database

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var rulesFirewall string
var waitDatabase bool

var dbCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo db create <DATABASE-NAME> --size <SIZE> --software <SOFTWARE_NAME> --version <SOFTWARE_VERSION>",
	Short:   "Create a new database",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		check, region, err := utility.CheckAvailability("dbaas", common.RegionSet)
		if err != nil {
			utility.Error("Error checking availability %s", err)
			os.Exit(1)
		}

		if !check {
			utility.Error("Sorry you can't create a database in the %s region", region)
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		var network *civogo.Network
		if networkID != "" {
			network, err = client.FindNetwork(networkID)
			if err != nil {
				utility.Error("the network %s doesn't exist", networkID)
				os.Exit(1)
			}
		} else {
			network, err = client.GetDefaultNetwork()
			if err != nil {
				utility.Error("Unable to find the default network: %s", err)
				os.Exit(1)
			}
		}

		if firewallID != "" {
			_, err = client.FindFirewall(firewallID)
			if err != nil {
				utility.Error("the firewall %s doesn't exist", firewallID)
				os.Exit(1)
			}
		}

		sizes, err := client.ListInstanceSizes()
		if err != nil {
			utility.Error("Unable to list sizes %s", err)
			os.Exit(1)
		}

		sizeIsValid := false
		if size != "" {
			for _, s := range sizes {
				if s.Name == size && utility.SizeType(s.Name) == "Database" {
					sizeIsValid = true
					break
				}
			}
			if !sizeIsValid {
				utility.Error("The provided size is not valid")
				os.Exit(1)
			}
		}

		dbVersions, err := client.ListDBVersions()
		if err != nil {
			utility.Error("Failed to fetch database versions: %s", err)
			os.Exit(1)
		}

		// Set default software to MySQL if not specified
		if software == "" {
			software = "postgresql"
		}

		if software == "mysql" {
			utility.Warning("The provided software will be deprecated soon. Please check the documentation https://www.civo.com/docs/database/mysql/dump-mysql.")
		}
		software = strings.ToLower(software)

		validSoftwares := map[string][]string{
			"mysql":      {"mysql"},
			"postgresql": {"postgresql", "psql"},
		}

		apiSoftwareNames := map[string]string{
			"mysql":      "MySQL",
			"postgresql": "PostgreSQL",
		}

		canonicalSoftwareName := ""
		softwareIsValid := false
		for swName, aliases := range validSoftwares {
			for _, alias := range aliases {
				if alias == software {
					softwareIsValid = true
					canonicalSoftwareName = apiSoftwareNames[swName]
					break
				}
			}
		}

		if !softwareIsValid {
			utility.Error("The provided software name is not valid. Valid options are mysql, psql or postgresql")
			os.Exit(1)
		}

		// If version is not specified, get the latest version
		if softwareVersion == "" {
			versions := dbVersions[canonicalSoftwareName]
			if len(versions) > 0 {
				// Assuming versions are sorted with latest first
				softwareVersion = versions[0].SoftwareVersion
			} else {
				utility.Error("No versions available for %s", canonicalSoftwareName)
				os.Exit(1)
			}
		} else {
			// Verify the specified version is valid
			versionValid := false
			for _, v := range dbVersions[canonicalSoftwareName] {
				if v.SoftwareVersion == softwareVersion {
					versionValid = true
					break
				}
			}
			if !versionValid {
				utility.Error("The provided software version is not valid. Please check the available versions for the specified software.")
				os.Exit(1)
			}
		}

		configDB := civogo.CreateDatabaseRequest{
			Name:            args[0],
			Size:            size,
			NetworkID:       network.ID,
			Nodes:           nodes,
			FirewallID:      firewallID,
			FirewallRules:   rulesFirewall,
			Software:        canonicalSoftwareName,
			SoftwareVersion: softwareVersion,
			Region:          client.Region,
		}

		db, err := client.NewDatabase(&configDB)
		if err != nil {
			utility.Error("Error creating database %s", err)
			os.Exit(1)
		}

		var executionTime string

		if waitDatabase {
			startTime := utility.StartTime()

			stillCreating := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Writer = os.Stderr
			s.Prefix = fmt.Sprintf("Create a database called %s ", db.Name)
			s.Start()

			for stillCreating {
				databaseCheck, err := client.FindDatabase(db.ID)
				if err != nil {
					utility.Error("Database %s", err)
					os.Exit(1)
				}
				if databaseCheck.Status == "Ready" {
					stillCreating = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}

			executionTime = utility.TrackTime(startTime)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": db.ID, "name": db.Name})
		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if executionTime != "" {
				fmt.Printf("Database (%s) type %s version %s with ID %s and size %s has been created in %s\nTo get fetch the database credentials use the command:\n\ncivo database credentials %s\n", utility.Green(db.Name), strings.ToLower(db.Software), db.SoftwareVersion, db.ID, db.Size, executionTime, db.Name)
			} else {
				fmt.Printf("Database (%s) type %s version %s with ID %s and size %s has been created\nTo get fetch the database credentials use the command:\n\ncivo database credentials %s\n", utility.Green(db.Name), strings.ToLower(db.Software), db.SoftwareVersion, db.ID, db.Size, db.Name)
			}
		}
	},
}
