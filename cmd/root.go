package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// {
//   "apikeys":{
//     "andy@civo.com":"DH48mUq3VIg0drpbkiO9RBsAJxuQKnSCTcZ2eGWYEPw7yFM6lv",
//     "andy@andyjeffries.co.uk":"mfFNbE0yR8qIMrzSbLHzENYB7QrfhkA1XZdQ08BBITvQTXJpro",
//     "demo@civo.com":"0G3EtRnkBxQdvzhsmACeLKlP4OUFNJVI51Zf67a9jMucY2qXiT"
//   },
//   "meta":{
//     "admin":false,
//     "current_apikey":"demo@civo.com",
//     "default_region":"lon1",
//     "latest_release_check":"2019-02-20T20:33:02Z",
//     "url":"https://api.civo.com"
//   }
// }

// Config describes the configuration for Civo's CLI
type Config struct {
	APIKeys map[string]string `json:"apikeys"`
	Meta    struct {
		Admin              bool      `json:"admin"`
		CurrentAPIKey      string    `json:"current_apikey"`
		DefaultRegion      string    `json:"default_region"`
		LatestReleaseCheck time.Time `json:"latest_release_check"`
		URL                string    `json:"url"`
	} `json:"meta"`
}

var cfgFile string
var CurrentConfig Config
var OutputFields string
var OutputFormat string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "civo",
	Short: "CLI to manage cloud resources at Civo.com",
	Long: `civo is a CLI library for managing cloud resources such
as instances and Kubernetes clusters at Civo.com.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", "config file (default is $HOME/.civo.json)")
	rootCmd.PersistentFlags().StringVarP(&OutputFields, "fields", "f", "", "output fields (use -h to determine fields)")
	rootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "human", "output format (json/human/custom)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
