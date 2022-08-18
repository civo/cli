package kubernetes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesAppRemoveCmd = &cobra.Command{
	Use:     "remove",
	Example: "civo kubernetes application remove NAME --cluster CLUSTER_NAME",
	Aliases: []string{"rm", "uninstall"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Remove the marketplace application from a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}
		kube, err := client.FindKubernetesCluster(kubernetesClusterApp)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}
		if !kube.Ready {
			utility.Error("The cluster isn't ready yet, so the KUBECONFIG isn't available.")
			os.Exit(1)
		}

		allApps := strings.Split(args[0], ",")
		tmpFile, err := ioutil.TempFile(os.TempDir(), "kubeconfig-")
		if err != nil {
			utility.Error("Cannot create temporary file", err)
		}
		if _, err = tmpFile.Write([]byte(kube.KubeConfig)); err != nil {
			utility.Error("Failed to write to temporary file", err)
		}
		defer os.Remove(tmpFile.Name())
		for _, split := range allApps {
			appName := split
			// TODO: Ideally this would come from the Civo API, but the Civo API doesn't currently return uninstall.sh
			// https://www.civo.com/api/kubernetes#listing-applications
			filepath := fmt.Sprintf("bash <(curl -s https://raw.githubusercontent.com/civo/kubernetes-marketplace/master/%s/uninstall.sh)", appName)
			cmdConfig := exec.Command("/bin/bash", "-c", filepath)
			var b bytes.Buffer
			cmdConfig.Stdout = &b
			cmdConfig.Stderr = &b
			cmdConfig.Env = os.Environ()
			cmdConfig.Env = append(cmdConfig.Env, "KUBECONFIG="+tmpFile.Name())

			if err := cmdConfig.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					utility.Error("Failed to uninstall application %s (exited with code %d)\n", appName, exitError.ExitCode())
					cmd := exec.Command("curl", "-s", fmt.Sprintf("https://raw.githubusercontent.com/civo/kubernetes-marketplace/master/%s/uninstall.sh", appName))
					output, _ := cmd.CombinedOutput()
					fmt.Println("--------------- Uninstall script ---------------")
					fmt.Println(string(output))
					fmt.Println("--------------- Uninstall output ---------------")
					fmt.Println(b.String())
					os.Exit(1)
				} else {
					utility.Error("Failed to uninstall application %s because of %s\n", appName, err.Error())
					os.Exit(1)
				}
			}
		}

		result := utility.RemoveApplicationFromInstalledList(kube.InstalledApplications, args[0])
		configKubernetes := &civogo.KubernetesClusterConfig{
			Applications: result,
		}

		_, err = client.UpdateKubernetesCluster(kube.ID, configKubernetes)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

	},
}
