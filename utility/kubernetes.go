package utility

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/civo/civogo"
)

// ObtainKubeConfig is the function to get the kubeconfig from the cluster
// and save to the file or merge with the existing one
func ObtainKubeConfig(KubeconfigFilename string, civoConfig string, merge bool, switchContext bool, clusterName string) error {

	kubeConfig := []byte(civoConfig)

	if merge {
		var err error
		kubeConfig, err = mergeConfigs(KubeconfigFilename, kubeConfig, switchContext, clusterName)
		if err != nil {
			return err
		}
	}

	if writeErr := writeConfig(KubeconfigFilename, kubeConfig, false, merge, switchContext, clusterName); writeErr != nil {
		return writeErr
	}

	if merge && switchContext {
		_, err := switchKubernetesContext(clusterName)
		if err != nil {
			YellowConfirm("%s\n", err)
			os.Exit(1)
		}
	}

	return nil
}

func mergeConfigs(localKubeconfigPath string, k3sconfig []byte, switchContext bool, clusterName string) ([]byte, error) {
	// Create a temporary kubeconfig to store the config of the newly create k3s cluster
	file, err := ioutil.TempFile(os.TempDir(), "civo-temp-*")
	if err != nil {
		return nil, fmt.Errorf("could not generate a temporary file to store the kuebeconfig: %s", err)
	}
	defer file.Close()

	if writeErr := writeConfig(file.Name(), k3sconfig, true, true, switchContext, clusterName); writeErr != nil {
		return nil, writeErr
	}

	fmt.Printf("Merged with main kubernetes config: %s\n", Green(localKubeconfigPath))

	// Merge the two kubeconfigs and read the output into 'data'
	osResult := CheckOS()
	var cmd *exec.Cmd

	if osResult == "windows" {
		os.Setenv("KUBECONFIG", fmt.Sprintf("%s;%s", localKubeconfigPath, file.Name()))
		cmd = exec.Command("powershell", "kubectl", "config", "view", "--merge", "--flatten")
	} else {
		// Append KUBECONFIGS in ENV Vars
		appendKubeConfigENV := fmt.Sprintf("KUBECONFIG=%s:%s", localKubeconfigPath, file.Name())
		cmd = exec.Command("kubectl", "config", "view", "--merge", "--flatten")
		cmd.Env = append(os.Environ(), appendKubeConfigENV)
	}

	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("could not merge kubeconfigs: %s", err)
	}

	// To be able to remove the file, it aught to be closed.
	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close temporary kubeconfig file: %s, %s", file.Name(), err)
	}

	// Remove the temporarily generated file
	err = os.Remove(file.Name())
	if err != nil {
		return nil, fmt.Errorf("could not remove temporary kubeconfig file: %s, %s", file.Name(), err)
	}

	return data, nil
}

// Generates config files give the path to file: string and the data: []byte
func writeConfig(path string, data []byte, suppressMessage bool, mergeConfigs bool, switchConfig bool, clusterName string) error {
	if !suppressMessage {
		fmt.Print("\nAccess your cluster with:\n")
		if mergeConfigs {
			if !switchConfig {
				fmt.Printf("kubectl config use-context %s\n", strings.ToLower(clusterName))
			}
			fmt.Println("kubectl get node")
		} else {
			if strings.Contains(path, ".kube") {
				fmt.Print("kubectl get node\n")
			} else {
				fmt.Printf("KUBECONFIG=%s kubectl get node\n", path)
			}
		}
	}

	// we check if the .kube dir is present
	checkKubeDir()

	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			Error(err.Error())
		}
		defer file.Close()
	}

	writeErr := ioutil.WriteFile(path, data, 0600)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

// checkAppPlan is the function to verify if the application to be installed in the cluster
// has a plan or not, in case it has a plan but does not specify it,
// we use the first one in the list
func checkAppPlan(appList []civogo.KubernetesMarketplaceApplication, requested string) (string, error) {
	foundIndex := -1
	parts := strings.SplitN(requested, ":", 2)
	appName := parts[0]

	var planName string
	if len(parts) > 1 {
		planName = parts[1]
	}

	for i, app := range appList {
		if strings.Contains(app.Name, appName) {
			if foundIndex != -1 {
				fmt.Printf("unable to find %s because there were multiple matches", appName)
			}
			foundIndex = i
		}
	}
	if foundIndex == -1 {
		YellowConfirm("you are trying to install the application %s, but this application does not exist\n", appName)
		os.Exit(1)
	}

	if len(appList[foundIndex].Plans) > 0 {
		allPlan := []string{}

		for pk := range appList[foundIndex].Plans {
			allPlan = append(allPlan, appList[foundIndex].Plans[pk].Label)
		}

		_, found := find(allPlan, planName)
		if !found {
			YellowConfirm("the plan you gave doesn't exist for %s; we've picked a default one for you\n", appName)
			return fmt.Sprintf("%s:%s", appName, appList[foundIndex].Plans[0].Label), nil
		}

		return requested, nil
	}

	if planName != "" {
		YellowConfirm("you are trying to install %s application with a plan but this application has no plans\n", appName)
		os.Exit(1)
	}

	return requested, nil
}

// RequestedSplit is a function to split all app requested to be install
func RequestedSplit(appList []civogo.KubernetesMarketplaceApplication, requested string) string {
	allsplit := strings.Split(requested, ",")
	allApp := []string{}

	for i := range allsplit {
		checkApp, err := checkAppPlan(appList, allsplit[i])
		if err != nil {
			fmt.Print(err)
		}

		allApp = append(allApp, checkApp)
	}

	return strings.Join(allApp, ", ")
}

// function to search if the string is in the slice
func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// switchKubernetesContext function to change automatic to the context
func switchKubernetesContext(context string) (bool, error) {
	var cmd *exec.Cmd

	if CheckOS() == "windows" {
		cmd = exec.Command("powershell", "kubectl", "config", "use-context", context)
	} else {
		cmd = exec.Command("kubectl", "config", "use-context", context)
	}

	_, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("could not change to the context: (%s) %s", context, err)
	}
	return true, nil
}

func checkKubeDir() {

	home, err := os.UserHomeDir()
	if err != nil {
		Error("%s", err)
		os.Exit(1)
	}

	if _, err := os.Stat(fmt.Sprintf("%s/.kube/", home)); os.IsNotExist(err) {
		os.Mkdir(fmt.Sprintf("%s/.kube/", home), 0755)
	}
}

// RemoveNodePool is a utility function to remove node pool from a kuberentes cluster
func RemoveNodePool(s []civogo.KubernetesClusterPoolConfig, id string) []civogo.KubernetesClusterPoolConfig {
	key := 0
	for k, v := range s {
		if strings.Contains(v.ID, id) {
			key = k
			break
		}
	}

	s[len(s)-1], s[key] = s[key], s[len(s)-1]
	return s[:len(s)-1]
}

// UpdateNodePool is a utility function to update node pool from a kuberentes cluster
func UpdateNodePool(s []civogo.KubernetesClusterPoolConfig, id string, count int) []civogo.KubernetesClusterPoolConfig {
	for k, v := range s {
		if strings.Contains(v.ID, id) {
			s[k].Count = count
			break
		}
	}
	return s
}
