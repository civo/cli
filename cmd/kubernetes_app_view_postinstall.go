package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"

	markdown "github.com/MichaelMure/go-term-markdown"
)

var kubernetesAppViewPostInstallCmd = &cobra.Command{
	Use:     "post-install",
	Aliases: []string{"postinstall", "pi"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "View the post install instructions for a given kubernetes applications",
	Long: `Downloads and renders the kubernetes applications post_install.md file to the terminal`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		kubeApps, err := client.ListKubernetesMarketplaceApplications()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		 app := utility.RequestedSplit(kubeApps, args[0])
		 if len(app) == 0 {
		 	utility.Error("Kubernetes application %s not found in the marketplace")
		 	os.Exit(1)
		}

		appCount := len(strings.Split(app, ","))
		if appCount > 1  {
			utility.Error("Please specify a single application to view the post-install instructions for got %d '%s'",appCount, app )
			os.Exit(1)
		}


		url := fmt.Sprintf("https://raw.githubusercontent.com/civo/kubernetes-marketplace/master/%s/post_install.md",
			strings.ToLower(app))
		resp, err := http.Get(url)
		if err != nil {
			utility.Error("Failed to get post_install.md from %s - %s", err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				utility.Error("%s", err)
			}
		}(resp.Body)

		postinstall, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			utility.Error("%s", err)
		}
		result := markdown.Render(string(postinstall), 80, 6)

		fmt.Println(string(result))

	},
}
