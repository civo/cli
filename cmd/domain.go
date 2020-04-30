package cmd

// domain list -- list all domain records for a domain [ls, all]
// domain create DOMAIN -- create a new domain name called DOMAIN [new]
// domain remove ID -- remove the domain with ID (or name) [delete, destroy, rm]

// domain_record list DOMAIN_ID -- list all entries for DOMAIN_ID (or name) [ls, all]
// domain_record show RECORD_ID -- show full information for record RECORD_ID (or full DNS name) [get, inspect]
// domain_record create RECORD TYPE VALUE -- create a new domain record called RECORD TYPE(a/alias, cname/canonical, mx/mail, txt/text) VALUE [new]
// domain_record remove ID -- remove the domain record with ID [delete, destroy, rm]

import (
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:     "domain",
	Aliases: []string{"domains"},
	Short:   "Details of Civo domains",
}

func init() {
	rootCmd.AddCommand(domainCmd)
	domainCmd.AddCommand(domainListCmd)
	domainCmd.AddCommand(domainCreateCmd)
	domainCmd.AddCommand(domainRemoveCmd)
}
