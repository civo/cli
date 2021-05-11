package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:     "domain",
	Aliases: []string{"domains"},
	Short:   "Details of Civo domains",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

var domainRecordCmd = &cobra.Command{
	Use:     "record",
	Aliases: []string{"records"},
	Short:   "Details of Civo domains records",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("command is required")
	},
}

func init() {
	rootCmd.AddCommand(domainCmd)
	domainCmd.AddCommand(domainListCmd)
	domainCmd.AddCommand(domainCreateCmd)
	domainCmd.AddCommand(domainRemoveCmd)

	// Domains record cmd
	domainCmd.AddCommand(domainRecordCmd)
	domainRecordCmd.AddCommand(domainRecordListCmd)
	domainRecordCmd.AddCommand(domainRecordCreateCmd)
	domainRecordCmd.AddCommand(domainRecordShowCmd)
	domainRecordCmd.AddCommand(domainRecordRemoveCmd)

	/*
		Flags for domain record create cmd
	*/
	domainRecordCreateCmd.Flags().StringVarP(&recordName, "name", "n", "", "the name of the record")
	domainRecordCreateCmd.Flags().StringVarP(&recordType, "type", "e", "", "type of the record (A, CNAME, TXT, SRV, MX)")
	domainRecordCreateCmd.Flags().StringVarP(&recordValue, "value", "v", "", "the value of the record")
	domainRecordCreateCmd.Flags().IntVarP(&recordTTL, "ttl", "t", 600, "The TTL of the record")
	domainRecordCreateCmd.Flags().IntVarP(&recordPriority, "priority", "p", 0, "the priority of record only for SRV and MX record")

}
