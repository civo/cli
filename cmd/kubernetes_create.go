package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var numTargetNodes int
var rulesFirewall string
var waitKubernetes, saveConfigKubernetes, mergeConfigKubernetes, switchConfigKubernetes, createFirewall bool
var kubernetesVersion, targetNodesSize, clusterName, applications, removeapplications, networkID, existingFirewall, cniPlugin string
var kubernetesCluster *civogo.KubernetesCluster

var kubernetesCreateCmdExample = `civo kubernetes create CLUSTER_NAME [flags]

Notes:
* The '--create-firewall' will open the ports 80,443 and 6443 in the firewall if '--firewall-rules' is not used.
* The '--create-firewall' and '--existing-firewall' flags are mutually exclusive. You can't use them together.
* The '--firewall-rules' flag need to be used with '--create-firewall'.
* The '--firewall-rules' flag can accept:
    * You can pass 'all' to open all ports.
    * An optional end port using 'start_port-end_port' format (e.g. 8000-8100)
    * An optional CIDR notation (e.g. 0.0.0.0/0)
    * When no CIDR notation is provided, the port will get 0.0.0.0/0 (open to public) as default CIDR notation
    * When a CIDR notation is provided without slash and number segment, it will default to /32
    * Within a rule, you can use comma separator for multiple ports to have same CIDR notation
    * To separate between rules, you can use semicolon symbol and wrap everything in double quotes (see below)
    So the following would all be valid:
    * "80,443,6443:0.0.0.0/0;8080:1.2.3.4" (open 80,443,6443 to public and 8080 just for 1.2.3.4/32)
    * "80,443,6443;6000-6500:4.4.4.4/24" (open 80,443,6443 to public and 6000 to 6500 just for 4.4.4.4/24)
`

var kubernetesCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: kubernetesCreateCmdExample,
	Short:   "Create a new Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		check, region, err := utility.CheckAvailability("kubernetes", common.RegionSet)
		if err != nil {
			utility.Error("Error checking availability %s", err)
			os.Exit(1)
		}

		if !check {
			utility.Error("Sorry you can't create a kubernetes cluster in the %s region", region)
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if !strings.Contains(targetNodesSize, ".kube.") {
			k8sSize, err := utility.GetK3sSize()
			if err != nil {
				utility.Error("Error %s", err)
				os.Exit(1)
			}
			utility.Error("You can't create a cluster with the specified size %s. Possible values: %s", targetNodesSize, k8sSize)
			os.Exit(1)
		}

		if !waitKubernetes {
			if saveConfigKubernetes || switchConfigKubernetes || mergeConfigKubernetes {
				utility.Error("you can't use --save, --switch or --merge without --wait")
				os.Exit(1)
			}
		} else {
			if mergeConfigKubernetes && !saveConfigKubernetes {
				utility.Error("you can't use --merge without --save")
				os.Exit(1)
			}
			if switchConfigKubernetes && !saveConfigKubernetes {
				utility.Error("you can't use --switch without --save")
				os.Exit(1)
			}
		}

		if len(args) > 0 {
			if utility.ValidNameLength(args[0]) {
				utility.Warning("the cluster name cannot be longer than 63 characters")
				os.Exit(1)
			}
			clusterName = args[0]

		} else {
			clusterName = utility.RandomName()
		}

		var network = &civogo.Network{}
		if networkID == "default" {
			network, err = client.GetDefaultNetwork()
			if err != nil {
				utility.Error("Network %s", err)
				os.Exit(1)
			}
		} else {
			network, err = client.FindNetwork(networkID)
			if err != nil {
				utility.Error("Network %s", err)
				os.Exit(1)
			}
		}

		cni := strings.ToLower(cniPlugin)
		if cni != "flannel" && cni != "cilium" {
			utility.Error("CNI plugin provided isn't valid/supported")
			os.Exit(1)
		}
		configKubernetes := &civogo.KubernetesClusterConfig{
			Name:            clusterName,
			NumTargetNodes:  numTargetNodes,
			TargetNodesSize: targetNodesSize,
			NetworkID:       network.ID,
			CNIPlugin:       cni,
		}

		if rulesFirewall != "default" && !createFirewall {
			utility.Error("You can't use --firewall-rules without --create-firewall flag")
			os.Exit(1)
		}

		if createFirewall && rulesFirewall == "default" {
			configKubernetes.FirewallRule = "80,443,6443"
		}
		if createFirewall && rulesFirewall != "default" {
			configKubernetes.FirewallRule = rulesFirewall
		}

		if kubernetesVersion != "latest" {
			configKubernetes.KubernetesVersion = kubernetesVersion
		}

		defaultApplications, err := utility.ListDefaultApps()
		if err != nil {
			utility.Error("Error %s", err)
			os.Exit(1)
		}
		apps := InstallApps(defaultApplications, applications, removeapplications)
		for _, app := range apps {
			if !utility.CheckAPPName(app) {
				utility.Error("%q is not a valid application name", app)
				os.Exit(1)
			}
		}
		configKubernetes.Applications = strings.Join(apps, ",")

		// We check if the user don't expecify a firewall we send to create a new one with the default rules
		if !createFirewall && rulesFirewall == "default" {
			configKubernetes.FirewallRule = "80,443,6443"
		}

		if existingFirewall != "" {
			if createFirewall {
				utility.Error("You can't use --create-firewall together with --existing-firewall flag")
				os.Exit(1)
			}

			ef, err := client.FindFirewall(existingFirewall)
			if err != nil {
				utility.Error("Unable to find %q firewall - %s", existingFirewall, err)
				os.Exit(1)
			}

			if ef.NetworkID != network.ID {
				utility.Error("Unable to find firewall %q in %q network", ef.ID, network.Label)
				os.Exit(1)
			}

			configKubernetes.InstanceFirewall = ef.ID
			configKubernetes.FirewallRule = ""
		}

		if !mergeConfigKubernetes && saveConfigKubernetes {
			if utility.UserConfirmedOverwrite("kubernetes config", common.DefaultYes) {
				kubernetesCluster, err = client.NewKubernetesClusters(configKubernetes)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
			} else {
				fmt.Println("Operation aborted.")
				os.Exit(1)
			}
		} else {
			kubernetesCluster, err = client.NewKubernetesClusters(configKubernetes)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
		}

		var executionTime string

		if waitKubernetes {
			startTime := utility.StartTime()

			stillCreating := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = fmt.Sprintf("Creating a %s node k3s cluster of %s instances called %s... ", strconv.Itoa(kubernetesCluster.NumTargetNode), kubernetesCluster.TargetNodeSize, kubernetesCluster.Name)
			s.Start()

			for stillCreating {
				kubernetesCheck, err := client.FindKubernetesCluster(kubernetesCluster.ID)
				if err != nil {
					utility.Error("Kubernetes %s", err)
					os.Exit(1)
				}
				if kubernetesCheck.Ready {
					stillCreating = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}

			executionTime = utility.TrackTime(startTime)
		}

		if saveConfigKubernetes {
			kube, err := client.FindKubernetesCluster(kubernetesCluster.ID)
			if err != nil {
				utility.Error("Kubernetes %s", err)
				os.Exit(1)
			}

			err = utility.ObtainKubeConfig(localPathConfig, kube.KubeConfig, mergeConfigKubernetes, switchConfigKubernetes, kube.Name)
			if err != nil {
				utility.Error("Saving the cluster config failed with %s", err)
				os.Exit(1)
			}

		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": kubernetesCluster.ID, "name": kubernetesCluster.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if executionTime != "" {
				fmt.Printf("The cluster %s (%s) has been created in %s\n", utility.Green(kubernetesCluster.Name), kubernetesCluster.ID, executionTime)
			} else {
				fmt.Printf("The cluster %s (%s) has been created\n", utility.Green(kubernetesCluster.Name), kubernetesCluster.ID)
			}

		}
	},
}

// InstallApps returns the list of applications to install
func InstallApps(defaultApps []string, apps, removeApps string) []string {
	var iApps []string
	if apps != "" {
		iApps = strings.Split(apps, ",")
	}
	iApps = append(defaultApps, iApps...)

	if removeApps != "" {
		for i, v := range iApps {
			for _, v2 := range strings.Split(removeApps, ",") {
				if v == v2 {
					iApps = append(iApps[:i], iApps[i+1:]...)
				}
			}
		}
	}
	return iApps
}
