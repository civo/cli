package cmd

import (
	"strconv"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var quotaShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get"},
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
	* SecurityGroupRuleUsage`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			return
		}

		quota, err := client.GetQuota()
		if err != nil {
			utility.Error("Unable to get your quota %s", err)
			return
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		if outputFormat == "json" || outputFormat == "custom" {
			ow.AppendData("InstanceCountLimit", strconv.Itoa(quota.InstanceCountLimit))
			ow.AppendData("InstanceCountUsage", strconv.Itoa(quota.InstanceCountUsage))
			ow.AppendData("CPUCoreLimit", strconv.Itoa(quota.CPUCoreLimit))
			ow.AppendData("CPUCoreUsage", strconv.Itoa(quota.CPUCoreUsage))
			ow.AppendData("RAMMegabytesLimit", strconv.Itoa(quota.RAMMegabytesLimit))
			ow.AppendData("RAMMegabytesUsage", strconv.Itoa(quota.RAMMegabytesUsage))
			ow.AppendData("DiskGigabytesLimit", strconv.Itoa(quota.DiskGigabytesLimit))
			ow.AppendData("DiskGigabytesUsage", strconv.Itoa(quota.DiskGigabytesUsage))
			ow.AppendData("DiskVolumeCountLimit", strconv.Itoa(quota.DiskVolumeCountLimit))
			ow.AppendData("DiskVolumeCountUsage", strconv.Itoa(quota.DiskVolumeCountUsage))
			ow.AppendData("DiskSnapshotCountLimit", strconv.Itoa(quota.DiskSnapshotCountLimit))
			ow.AppendData("DiskSnapshotCountUsage", strconv.Itoa(quota.DiskSnapshotCountUsage))
			ow.AppendData("PublicIPAddressLimit", strconv.Itoa(quota.PublicIPAddressLimit))
			ow.AppendData("PublicIPAddressUsage", strconv.Itoa(quota.PublicIPAddressUsage))
			ow.AppendData("SubnetCountLimit", strconv.Itoa(quota.SubnetCountLimit))
			ow.AppendData("SubnetCountUsage", strconv.Itoa(quota.SubnetCountUsage))
			ow.AppendData("NetworkCountLimit", strconv.Itoa(quota.NetworkCountLimit))
			ow.AppendData("NetworkCountUsage", strconv.Itoa(quota.NetworkCountUsage))
			ow.AppendData("SecurityGroupLimit", strconv.Itoa(quota.SecurityGroupLimit))
			ow.AppendData("SecurityGroupUsage", strconv.Itoa(quota.SecurityGroupUsage))
			ow.AppendData("SecurityGroupRuleLimit", strconv.Itoa(quota.SecurityGroupRuleLimit))
			ow.AppendData("SecurityGroupRuleUsage", strconv.Itoa(quota.SecurityGroupRuleUsage))
		} else {
			ow.AppendData("Instance Count Limit", strconv.Itoa(quota.InstanceCountLimit))
			ow.AppendData("Instance Count Usage", strconv.Itoa(quota.InstanceCountUsage))
			ow.AppendData("CPUCore Limit", strconv.Itoa(quota.CPUCoreLimit))
			ow.AppendData("CPUCore Usage", strconv.Itoa(quota.CPUCoreUsage))
			ow.AppendData("RAM Megabytes Limit", strconv.Itoa(quota.RAMMegabytesLimit))
			ow.AppendData("RAM Megabytes Usage", strconv.Itoa(quota.RAMMegabytesUsage))
			ow.AppendData("Disk Gigabytes Limit", strconv.Itoa(quota.DiskGigabytesLimit))
			ow.AppendData("Disk Gigabytes Usage", strconv.Itoa(quota.DiskGigabytesUsage))
			ow.AppendData("Disk Volume Count Limit", strconv.Itoa(quota.DiskVolumeCountLimit))
			ow.AppendData("Disk Volume Count Usage", strconv.Itoa(quota.DiskVolumeCountUsage))
			ow.AppendData("Disk Snapshot Count Limit", strconv.Itoa(quota.DiskSnapshotCountLimit))
			ow.AppendData("Disk Snapshot Count Usage", strconv.Itoa(quota.DiskSnapshotCountUsage))
			ow.AppendData("Public IP Address Limit", strconv.Itoa(quota.PublicIPAddressLimit))
			ow.AppendData("Public IP Address Usage", strconv.Itoa(quota.PublicIPAddressUsage))
			ow.AppendData("Subnet Count Limit", strconv.Itoa(quota.SubnetCountLimit))
			ow.AppendData("Subnet Count Usage", strconv.Itoa(quota.SubnetCountUsage))
			ow.AppendData("Network Count Limit", strconv.Itoa(quota.NetworkCountLimit))
			ow.AppendData("Network Count Usage", strconv.Itoa(quota.NetworkCountUsage))
			ow.AppendData("Security Group Limit", strconv.Itoa(quota.SecurityGroupLimit))
			ow.AppendData("Security Group Usage", strconv.Itoa(quota.SecurityGroupUsage))
			ow.AppendData("Security Group Rule Limit", strconv.Itoa(quota.SecurityGroupRuleLimit))
			ow.AppendData("Security Group Rule Usage", strconv.Itoa(quota.SecurityGroupRuleUsage))
		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
