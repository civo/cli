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
	Args:    cobra.MinimumNArgs(0), // Change from 1 to 0 to make the name argument optional
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

		software = strings.ToLower(software)
		softwareIsValid := false
		softwareVersionIsValid := false

		validSoftwares := map[string][]string{
			"mysql":      {"mysql"},
			"postgresql": {"postgresql", "psql"},
		}

		apiSoftwareNames := map[string]string{
			"mysql":      "MySQL",
			"postgresql": "PostgreSQL",
		}

		canonicalSoftwareName := ""
		for swName, aliases := range validSoftwares {
			for _, alias := range aliases {
				if alias == software {
					softwareIsValid = true
					canonicalSoftwareName = apiSoftwareNames[swName]
					for _, v := range dbVersions[canonicalSoftwareName] {
						if v.SoftwareVersion == softwareVersion {
							softwareVersionIsValid = true
							break
						}
					}
				}
			}
		}

		if !softwareIsValid {
			utility.Error("The provided software name is not valid. valid options are mysql, psql or postgresql")
			os.Exit(1)
		}

		if !softwareVersionIsValid {
			if softwareVersion == "" {
				utility.Error("No version specified for %s. Please provide a version using --version flag. For example, civo database create db-psql --software psql --version 14", canonicalSoftwareName)
			} else {
				utility.Error("The provided software version is not valid. Please check the available versions for the specified software")
			}
			os.Exit(1)
		}

		var dbName string
		if len(args) > 0 {
			if utility.ValidNameLength(args[0]) {
				utility.Warning("the database name cannot be longer than 63 characters")
				os.Exit(1)
			}
			dbName = args[0]
		} else {
			dbName = utility.RandomName()
		}

		configDB := civogo.CreateDatabaseRequest{
			Name:            dbName,
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
				fmt.Printf("Database %s (%s) has been created in %s\n", utility.Green(db.Name), db.ID, executionTime)
			} else {
				fmt.Printf("Database (%s) with ID %s has been created\n", utility.Green(db.Name), db.ID)
			}
		}
	},
}
