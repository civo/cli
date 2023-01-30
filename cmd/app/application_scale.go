package app

// import (
// 	"fmt"
// 	"os"

// 	"github.com/civo/civogo"
// 	"github.com/civo/cli/common"
// 	"github.com/civo/cli/config"
// 	"github.com/civo/cli/utility"
// 	"github.com/spf13/cobra"
// )

// var processType string
// var processCount int

// var appScaleCmd = &cobra.Command{
// 	Use:     "scale",
// 	Aliases: []string{"change", "modify", "upgrade"},
// 	Example: "civo app scale APP-NAME --process-type=web --process-count=3",
// 	Short:   "Scale processes of your application",
// 	Args:    cobra.MinimumNArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {

// 		client, err := config.CivoAPIClient()
// 		if err != nil {
// 			utility.Error("Creating the connection to Civo's API failed with %s", err)
// 			os.Exit(1)
// 		}

// 		findApp, err := client.FindApplication(args[0])
// 		if err != nil {
// 			utility.Error("App %s", err)
// 			os.Exit(1)
// 		}

// 		processInfo := civogo.ProcessInfo{
// 			ProcessType:  processType,
// 			ProcessCount: processCount,
// 		}

// 		application := &civogo.UpdateApplicationRequest{
// 			ProcessInfo: append(findApp.ProcessInfo, processInfo),
// 		}

// 		app, err := client.UpdateApplication(findApp.ID, application)
// 		if err != nil {
// 			utility.Error("%s", err)
// 			os.Exit(1)
// 		}

// 		ow := utility.NewOutputWriterWithMap(map[string]string{"id": app.ID, "name": app.Name})

// 		switch common.OutputFormat {
// 		case "json":
// 			ow.WriteSingleObjectJSON(common.PrettySet)
// 		case "custom":
// 			ow.WriteCustomOutput(common.OutputFields)
// 		default:
// 			fmt.Printf("The application %s has been updated.\n", utility.Green(app.Name))
// 		}
// 	},
// }
