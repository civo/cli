package database

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var dbForkCmd = &cobra.Command{
	Use:     "fork",
	Aliases: []string{"clone", "copy"},
	Short:   "Fork a database into a new, independent database",
	Long: `Creates a new database from the backups of an existing one. By default the
fork reflects the source at its latest backup; pass --target-time to fork the
source as it was at that moment (point in time, no earlier than the restore
window shown by "civo db show").

The new database keeps the source's size, node count, network and extensions
unless overridden with the matching flags. --extensions replaces the source's
extensions with the listed ones; an empty value still keeps the source's, so to
fork without some of them, run "civo db update <NEW-NAME> --extensions ..."
once the fork is ready (an empty value there disables them all).`,
	Example: `civo db fork <SOURCE-NAME/ID> --name <NEW-NAME>
civo db fork <SOURCE-NAME/ID> --name <NEW-NAME> --target-time 2026-10-21T10:00:00Z
civo db fork <SOURCE-NAME/ID> --name <NEW-NAME> --size g3.db.medium --nodes 3`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		if forkName == "" {
			utility.Error("A name for the new database must be specified with --name")
			os.Exit(1)
		}
		source, err := client.FindDatabase(args[0])
		if err != nil {
			utility.Error("Database %s", err)
			os.Exit(1)
		}

		// The service does not derive the engine from the source, so send it
		// along or the fork would be created with the default engine.
		req := &civogo.ForkDatabaseRequest{
			Name:             forkName,
			Size:             source.Size,
			Nodes:            source.Nodes,
			Software:         source.Software,
			SoftwareVersion:  source.SoftwareVersion,
			SourceDatabaseID: source.ID,
			NetworkID:        source.NetworkID,
			Region:           client.Region,
		}
		if forkSize != "" {
			req.Size = forkSize
		}
		if forkNodes > 0 {
			req.Nodes = forkNodes
		}
		if forkNetworkID != "" {
			network, err := client.FindNetwork(forkNetworkID)
			if err != nil {
				utility.Error("the network %s doesn't exist", forkNetworkID)
				os.Exit(1)
			}
			req.NetworkID = network.ID
		}
		if forkFirewallID != "" {
			firewall, err := client.FindFirewall(forkFirewallID)
			if err != nil {
				utility.Error("the firewall %s doesn't exist", forkFirewallID)
				os.Exit(1)
			}
			req.FirewallID = firewall.ID
		}
		if storageGB > 0 {
			gb := storageGB
			req.StorageGB = &gb
		}
		if cmd.Flags().Changed("extensions") {
			req.Extensions = parseExtensions(extensions)
		}

		forkPoint := "its latest backup"
		if targetTime != "" {
			t, err := parseTargetTime(targetTime)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
			req.SourceTargetTime = &t
			forkPoint = fmt.Sprintf("its state at %s", utility.Green(t.Format(time.RFC3339)))
		}

		db, err := client.ForkDatabase(req)
		if err != nil {
			utility.Error("Error forking database %s", err)
			os.Exit(1)
		}

		var executionTime string
		if forkWait {
			startTime := utility.StartTime()
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Writer = os.Stderr
			s.Prefix = fmt.Sprintf("Forking %s into %s ", source.Name, db.Name)
			s.Start()
			for {
				check, err := client.FindDatabase(db.ID)
				if err != nil {
					s.Stop()
					utility.Error("Database %s", err)
					os.Exit(1)
				}
				if check.Status == "Ready" {
					break
				}
				if check.Status == "Failed" {
					s.Stop()
					utility.Error("The fork %s failed; check the database in the dashboard or with civo db show %s", db.Name, db.Name)
					os.Exit(1)
				}
				time.Sleep(2 * time.Second)
			}
			s.Stop()
			executionTime = utility.TrackTime(startTime)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": db.ID, "name": db.Name, "source_id": source.ID})
		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if executionTime != "" {
				fmt.Printf("Database %s has been forked from %s (%s) in %s\n", utility.Green(db.Name), utility.Green(source.Name), forkPoint, executionTime)
			} else {
				fmt.Printf("Database %s is being forked from %s (%s)\n", utility.Green(db.Name), utility.Green(source.Name), forkPoint)
			}
			fmt.Printf("To get the credentials, run : civo db credential %s\n", db.Name)
		}
	},
}
