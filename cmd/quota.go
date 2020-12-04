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
	* SecurityGroupRuleUsage`,
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
		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
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
