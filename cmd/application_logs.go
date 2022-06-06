package cmd

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func logCmd(cmd *cobra.Command, args []string) {
	client, err := config.CivoAPIClient()
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	findApp, err := client.FindApplication(args[0])
	if err != nil {
		utility.Error("App %s", err)
		os.Exit(1)
	}

	conf, err := client.GetApplicationLogAuth(findApp.ID)
	if err != nil {
		utility.Error("App %s", err)
		os.Exit(1)
	}

	conf = strings.ReplaceAll(conf, "\"", "")
	byteconf, err := base64.StdEncoding.DecodeString(conf)
	if err != nil {
		utility.Error("Decode failed %s", err)
		os.Exit(1)
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "kubeconfig-")
	if err != nil {
		utility.Error("Creating the file failed with %s", err)
		os.Exit(1)
	}
	defer tmpFile.Close()

	conf = string(byteconf)
	if _, err := tmpFile.WriteString(conf); err != nil {
		utility.Error("Writing the file failed with %s", err)
		os.Exit(1)
	}

	tenantConfig, err := clientcmd.BuildConfigFromFlags("", tmpFile.Name())
	if err != nil {
		utility.Error("Building the config failed with %s", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(tenantConfig)
	if err != nil {
		utility.Error("Creating the client failed with %s", err)
		os.Exit(1)
	}

	ns := strings.Split(conf, "namespace: ")[1]
	ns = strings.Split(ns, "\n")[0]
	pods, err := clientset.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		utility.Error("Failed to get logs application might not be running %s", err)
		os.Exit(1)
	}

	deploymets, err := clientset.AppsV1().Deployments(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		utility.Error("Failed to get logs application might not be running %s", err)
		os.Exit(1)
	}

	//TODO: Print logs
	for _, pod := range pods.Items {
		utility.Info("Pod: %s", pod.Name)
	}
	//TODO: Print logs
	for _, deploy := range deploymets.Items {
		utility.Info("Deployment: %s", deploy.Name)
	}

}
