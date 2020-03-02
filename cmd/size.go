package cmd

import (
	"fmt"
	"strconv"

	"github.com/civo/cli/config"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// sizeCmd represents the size command
var sizeCmd = &cobra.Command{
	Use:     "size",
	Aliases: []string{"sizes"},
	Short:   "Details of Civo instance sizes",
}

// sizeListCmd represents the command to list available API keys
var sizeListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List sizes",
	Long: `List all available sizes for instances or Kubernetes nodes.
If you wish to use a custom format, the available fields are:

* Name
* NiceName
* CPUCores
* RAMMegabytes
* DiskGigabytes
* Description
* Selectable

Example: civo size ls -o custom -f "Code: Name (size)"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			return
		}

		sizes, err := client.ListInstanceSizes()
		if err != nil {
			fmt.Printf("Unable to list sizes: %s\n", aurora.Red(err))
			return
		}

		if OutputFormat == "json" {
			fmt.Println(client.LastJSONResponse)
			return
		}

		data := make([][]string, len(sizes))
		for i, size := range sizes {
			var selectableLabel string
			if size.Selectable {
				selectableLabel = "Yes"
			} else {
				selectableLabel = "No"
			}

			data[i] = []string{size.Name, size.Description, strconv.Itoa(size.CPUCores), strconv.Itoa(size.RAMMegabytes), strconv.Itoa(size.DiskGigabytes), selectableLabel}
		}

		outputTable([]string{"Name", "Description", "CPU", "RAM (MB)", "Disk (GB)", "Selectable"}, data)
	},
}

func init() {
	rootCmd.AddCommand(sizeCmd)

	sizeCmd.AddCommand(sizeListCmd)
}
