package cmd

import (
	"strconv"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var quotaCmd = &cobra.Command{
	Use:     "quota",
	Aliases: []string{"quotas"},
	Example: `civo quota show -o custom -f "InstanceCountUsage/InstanceCountLimit"`,
	Short:   "Show quota",
	Long: `Show your current quota and usage.
If you wish to use a custom format, the available fields are:

	* InstanceCountLimit
	* InstanceCountUsage
	* CPUCoreLimit
	* CPUCoreUsage
	* RAMMegabytesLimit
	* RAMMegabytesUsage
	* DiskGigabytesLimit
	* DiskGigabytesUsage
	* DiskVolumeCountLimit
	* DiskVolumeCountUsage
	* DiskSnapshotCountLimit
	* DiskSnapshotCountUsage
	* PublicIPAddressLimit
	* PublicIPAddressUsage
	* SubnetCountLimit
	* SubnetCountUsage
	* NetworkCountLimit
	* NetworkCountUsage
	* SecurityGroupLimit
	* SecurityGroupUsage
	* SecurityGroupRuleLimit
	* SecurityGroupRuleUsage
	* LoadBalancerCountLimit
	* LoadBalancerCountUsage`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			return
		}

		quota, err := client.GetQuota()
		if err != nil {
			utility.Error("%s", err)
			return
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		if outputFormat == "json" || outputFormat == "custom" {
			ow.AppendDataWithLabel("instance_count_limit", strconv.Itoa(quota.InstanceCountLimit), "InstanceCountLimit")
			ow.AppendDataWithLabel("instance_count_usage", strconv.Itoa(quota.InstanceCountUsage), "InstanceCountUsage")
			ow.AppendDataWithLabel("cpu_core_limit", strconv.Itoa(quota.CPUCoreLimit), "InstanceCountUsage")
			ow.AppendDataWithLabel("cpu_core_usage", strconv.Itoa(quota.CPUCoreUsage), "InstanceCountUsage")
			ow.AppendDataWithLabel("ram_mb_limit", strconv.Itoa(quota.RAMMegabytesLimit), "InstanceCountUsage")
			ow.AppendDataWithLabel("ram_mb_usage", strconv.Itoa(quota.RAMMegabytesUsage), "InstanceCountUsage")
			ow.AppendDataWithLabel("disk_gb_limit", strconv.Itoa(quota.DiskGigabytesLimit), "InstanceCountUsage")
			ow.AppendDataWithLabel("disk_gb_usage", strconv.Itoa(quota.DiskGigabytesUsage), "InstanceCountUsage")
			ow.AppendDataWithLabel("disk_volume_count_limit", strconv.Itoa(quota.DiskVolumeCountLimit), "InstanceCountUsage")
			ow.AppendDataWithLabel("disk_volume_count_usage", strconv.Itoa(quota.DiskVolumeCountUsage), "InstanceCountUsage")
			ow.AppendDataWithLabel("disk_snapshot_count_limit", strconv.Itoa(quota.DiskSnapshotCountLimit), "InstanceCountUsage")
			ow.AppendDataWithLabel("disk_snapshot_count_usage", strconv.Itoa(quota.DiskSnapshotCountUsage), "InstanceCountUsage")
			ow.AppendDataWithLabel("public_ip_address_limit", strconv.Itoa(quota.PublicIPAddressLimit), "PublicIPAddressLimit")
			ow.AppendDataWithLabel("public_ip_address_usage", strconv.Itoa(quota.PublicIPAddressUsage), "PublicIPAddressUsage")
			ow.AppendDataWithLabel("subnet_count_limit", strconv.Itoa(quota.SubnetCountLimit), "SubnetCountLimit")
			ow.AppendDataWithLabel("subnet_count_usage", strconv.Itoa(quota.SubnetCountUsage), "SubnetCountUsage")
			ow.AppendDataWithLabel("network_count_limit", strconv.Itoa(quota.NetworkCountLimit), "NetworkCountLimit")
			ow.AppendDataWithLabel("network_count_usage", strconv.Itoa(quota.NetworkCountUsage), "NetworkCountUsage")
			ow.AppendDataWithLabel("security_group_limit", strconv.Itoa(quota.SecurityGroupLimit), "SecurityGroupLimit")
			ow.AppendDataWithLabel("security_group_usage", strconv.Itoa(quota.SecurityGroupUsage), "SecurityGroupUsage")
			ow.AppendDataWithLabel("security_group_rule_limit", strconv.Itoa(quota.SecurityGroupRuleLimit), "SecurityGroupRuleLimit")
			ow.AppendDataWithLabel("security_group_rule_usage", strconv.Itoa(quota.SecurityGroupRuleUsage), "SecurityGroupRuleUsage")
			ow.AppendDataWithLabel("loadbalancer_count_limit", strconv.Itoa(quota.LoadBalancerCountLimit), "LoadBalancerCountLimit")
			ow.AppendDataWithLabel("loadbalancer_count_usage", strconv.Itoa(quota.LoadBalancerCountUsage), "LoadBalancerCountUsage")
		} else {
			ow.AppendData("Instance Count", utility.CheckQuotaPercent(quota.InstanceCountLimit, quota.InstanceCountUsage))
			ow.AppendData("CPUCore", utility.CheckQuotaPercent(quota.CPUCoreLimit, quota.CPUCoreUsage))
			ow.AppendData("RAM Megabytes", utility.CheckQuotaPercent(quota.RAMMegabytesLimit, quota.RAMMegabytesUsage))
			ow.AppendData("Disk Gigabytes", utility.CheckQuotaPercent(quota.DiskGigabytesLimit, quota.DiskGigabytesUsage))
			ow.AppendData("Disk Volume Count", utility.CheckQuotaPercent(quota.DiskVolumeCountLimit, quota.DiskVolumeCountUsage))
			ow.AppendData("Disk Snapshot Count", utility.CheckQuotaPercent(quota.DiskSnapshotCountLimit, quota.DiskSnapshotCountUsage))
			ow.AppendData("Public IP Address", utility.CheckQuotaPercent(quota.PublicIPAddressLimit, quota.PublicIPAddressUsage))
			ow.AppendData("Subnet Count", utility.CheckQuotaPercent(quota.SubnetCountLimit, quota.SubnetCountUsage))
			ow.AppendData("Network Count", utility.CheckQuotaPercent(quota.NetworkCountLimit, quota.NetworkCountUsage))
			ow.AppendData("Security Group", utility.CheckQuotaPercent(quota.SecurityGroupLimit, quota.SecurityGroupUsage))
			ow.AppendData("Security Group Rule", utility.CheckQuotaPercent(quota.SecurityGroupRuleLimit, quota.SecurityGroupRuleUsage))
			ow.AppendData("LoadBalancer Count", utility.CheckQuotaPercent(quota.LoadBalancerCountLimit, quota.LoadBalancerCountUsage))
		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteKeyValues()
		}
	},
}

func init() {
	rootCmd.AddCommand(quotaCmd)
}
