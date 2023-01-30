package app

// import (
// 	"fmt"

// 	"github.com/civo/cli/common"
// 	"github.com/civo/cli/config"
// 	"github.com/civo/cli/utility"

// 	"os"

// 	"github.com/spf13/cobra"
// )

// var appConfigShowCmd = &cobra.Command{
// 	Use:     "show",
// 	Aliases: []string{"get", "inspect"},
// 	Args:    cobra.MinimumNArgs(1),
// 	Short:   "Show application config",
// 	Example: "civo app config show APP_NAME",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		client, err := config.CivoAPIClient()
// 		if err != nil {
// 			utility.Error("Creating the connection to Civo's API failed with %s", err)
// 			os.Exit(1)
// 		}

// 		app, err := client.FindApplication(args[0])
// 		if err != nil {
// 			utility.Error("%s", err)
// 			os.Exit(1)
// 		}

// 		ow := utility.NewOutputWriter()
// 		for _, config := range app.Config {
// 			fmt.Println(config)
// 		}

// 		switch common.OutputFormat {
// 		case "json":
// 			ow.WriteMultipleObjectsJSON(common.PrettySet)
// 		case "custom":
// 			ow.WriteCustomOutput(common.OutputFields)
// 		default:
// 			ow.WriteKeyValues()
// 		}
// 	},
// }
