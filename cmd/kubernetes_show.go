package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Example: `civo kubernetes show ID/HOSTNAME -o custom -f "ID: Code (DefaultUsername)"`,
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show Kubernetes cluster",
	Long: `Show a specified Kubernetes cluster.
If you wish to use a custom format, the available fields are:

	* ID
	* Code
	* Name
	* Nodes
	* Size
	* CPUCores
	* RAMMegabytes
	* DiskGigabytes
	* Status
	* KubernetesVersion
	* APIEndPoint
	* MasterIP
	* DNSEntry
	* InstalledApplications`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return getAllKubernetesList(), cobra.ShellCompDirectiveNoFileComp
		}
		return getKubernetesList(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		kubernetesCluster, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("Kubernetes %s", err)
			os.Exit(1)
		}

		//Get the firewall name by the id
		firewall, err := client.FindFirewall(kubernetesCluster.FirewallID)
		if err != nil {
			utility.Error("Firewall %s", err)
			os.Exit(1)
		}

		//Get the loadbalancer name by the id
		lbCluster := []civogo.LoadBalancer{}
		loadbalancer, err := client.ListLoadBalancers()
		if err != nil {
			utility.Error("Loadbalancer %s", err)
			os.Exit(1)
		}
		for _, lb := range loadbalancer {
			if lb.ClusterID == kubernetesCluster.ID {
				lbCluster = append(lbCluster, lb)
			}
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		ow.AppendData("ID", kubernetesCluster.ID)
		ow.AppendData("Name", kubernetesCluster.Name)
		ow.AppendData("Region", client.Region)
		ow.AppendData("Nodes", strconv.Itoa(len(kubernetesCluster.Instances)))
		ow.AppendData("Size", kubernetesCluster.TargetNodeSize)
		ow.AppendData("Status", kubernetesCluster.Status)
		ow.AppendData("Firewall", firewall.Name)

		if kubernetesCluster.UpgradeAvailableTo != "" {
			ow.AppendDataWithLabel("KubernetesVersion", utility.Red(kubernetesCluster.KubernetesVersion+" *"), "Version")
		} else {
			ow.AppendDataWithLabel("KubernetesVersion", kubernetesCluster.KubernetesVersion, "Version")
		}

		ow.AppendDataWithLabel("APIEndPoint", kubernetesCluster.APIEndPoint, "API Endpoint")
		ow.AppendDataWithLabel("MasterIP", kubernetesCluster.MasterIP, "External IP")
		ow.AppendDataWithLabel("DNSEntry", kubernetesCluster.DNSEntry, "DNS A record")

		if len(kubernetesCluster.InstalledApplications) > 0 {
			var appsList []string
			for _, app := range kubernetesCluster.InstalledApplications {
				appsList = append(appsList, app.Name)
			}
			ow.AppendDataWithLabel("InstalledApplications", strings.Join(appsList, ", "), "Installed Applications")
		}

		if outputFormat == "json" || outputFormat == "custom" {
			//ow.AppendData("CloudConfig", template.CloudConfig)

			if kubernetesCluster.UpgradeAvailableTo != "" {
				ow.AppendDataWithLabel("KubernetesVersion", kubernetesCluster.KubernetesVersion, "Version")
			} else {
				ow.AppendDataWithLabel("KubernetesVersion", kubernetesCluster.KubernetesVersion, "Version")
			}

			ow.AppendDataWithLabel("UpgradeAvailableTo", kubernetesCluster.UpgradeAvailableTo, "Upgrade Available")

			if outputFormat == "json" {
				ow.ToJSON(kubernetesCluster, prettySet)
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		} else {
			ow.WriteKeyValues()

			if kubernetesCluster.UpgradeAvailableTo != "" {
				var versionsList []string
				kubeVersions, err := client.ListAvailableKubernetesVersions()
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				for _, version := range kubeVersions {
					if version.Type == "deprecated" {
						continue
					}
					if version.Version > kubernetesCluster.Version {
						versionsList = append(versionsList, version.Version)
					}
				}
				fmt.Printf(utility.Red("\n* An upgrade to v%s is available. Learn more about how to upgrade: civo k3s upgrade --help"), strings.Join(versionsList, ", "))
				fmt.Println()
			}

			if len(lbCluster) > 0 {
				fmt.Println()
				ow.WriteHeader("Loadbalancers")
				owLB := utility.NewOutputWriter()
				for _, lb := range lbCluster {
					owLB.StartLine()
					owLB.AppendData("Name", lb.Name)
					owLB.AppendData("Algorithm", lb.Algorithm)
					owLB.AppendData("Public IP", lb.PublicIP)
					owLB.AppendData("Private IP", lb.PrivateIP)
					owLB.AppendData("State", lb.State)
					owLB.AppendData("Firewall", lb.FirewallID)
					owLB.AppendData("DNS Name", fmt.Sprintf("%s.lb.civo.com", lb.ID))
				}
				owLB.WriteTable()
			}

			if len(kubernetesCluster.Instances) > 0 {
				fmt.Println()
				for _, pool := range kubernetesCluster.Pools {
					ow.WriteHeader(fmt.Sprintf("Pool (%s)", pool.ID[:6]))
					owNode := utility.NewOutputWriter()

					for _, instance := range kubernetesCluster.Instances {
						for _, pinstance := range pool.InstanceNames {
							if instance.Hostname != "" && strings.Contains(pinstance, instance.Hostname[5:]) {
								owNode.StartLine()
								owNode.AppendData("Name", instance.Hostname)
								owNode.AppendData("IP", instance.PublicIP)
								owNode.AppendData("Status", instance.Status)
								owNode.AppendData("Size", instance.Size)
								owNode.AppendDataWithLabel("CPUCores", strconv.Itoa(instance.CPUCores), "Cpu Cores")
								owNode.AppendDataWithLabel("RAMMegabytes", strconv.Itoa(instance.RAMMegabytes), "RAM (MB)")
								owNode.AppendDataWithLabel("DiskGigabytes", strconv.Itoa(instance.DiskGigabytes), "SSD disk (GB)")
							}
						}
					}
					owNode.WriteTable()
					fmt.Println()
					ow.WriteHeader("Labels")
					fmt.Printf("kubernetes.civo.com/node-pool=%s\n", pool.ID)
					fmt.Printf("kubernetes.civo.com/node-size=%s\n", pool.Size)
				}

				if len(kubernetesCluster.InstalledApplications) > 0 {
					fmt.Println()
					ow.WriteHeader("Applications")
					owApp := utility.NewOutputWriter()

					for _, app := range kubernetesCluster.InstalledApplications {
						owApp.StartLine()

						owApp.AppendData("Name", app.Name)
						owApp.AppendData("Version", app.Version)
						owApp.AppendData("Installed", strconv.FormatBool(app.Installed))
						owApp.AppendData("Category", app.Category)
					}
					owApp.WriteTable()
					fmt.Println()
				}
			}
		}
	},
}
