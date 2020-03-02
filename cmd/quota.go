package cmd

import (
	"fmt"
	"strconv"

	"github.com/civo/cli/config"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// quotaCmd represents the quota command
var quotaCmd = &cobra.Command{
	Use:     "quota",
	Aliases: []string{"quotas"},
	Short:   "Your account's current quota settings and usage",
}

// quotaShowCmd represents the command to list available API keys
var quotaShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get"},
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

Example: civo quota show -o custom -f "InstanceCountUsage/InstanceCountLimit"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			return
		}

		quota, err := client.GetQuota()
		if err != nil {
			fmt.Printf("Unable to get your quota: %s\n", aurora.Red(err))
			return
		}

		switch OutputFormat {
		case "json":
			fmt.Println(client.LastJSONResponse)
		case "custom":
			data := map[string]string{
				"InstanceCountLimit":     strconv.Itoa(quota.InstanceCountLimit),
				"InstanceCountUsage":     strconv.Itoa(quota.InstanceCountUsage),
				"CPUCoreLimit":           strconv.Itoa(quota.CPUCoreLimit),
				"CPUCoreUsage":           strconv.Itoa(quota.CPUCoreUsage),
				"RAMMegabytesLimit":      strconv.Itoa(quota.RAMMegabytesLimit),
				"RAMMegabytesUsage":      strconv.Itoa(quota.RAMMegabytesUsage),
				"DiskGigabytesLimit":     strconv.Itoa(quota.DiskGigabytesLimit),
				"DiskGigabytesUsage":     strconv.Itoa(quota.DiskGigabytesUsage),
				"DiskVolumeCountLimit":   strconv.Itoa(quota.DiskVolumeCountLimit),
				"DiskVolumeCountUsage":   strconv.Itoa(quota.DiskVolumeCountUsage),
				"DiskSnapshotCountLimit": strconv.Itoa(quota.DiskSnapshotCountLimit),
				"DiskSnapshotCountUsage": strconv.Itoa(quota.DiskSnapshotCountUsage),
				"PublicIPAddressLimit":   strconv.Itoa(quota.PublicIPAddressLimit),
				"PublicIPAddressUsage":   strconv.Itoa(quota.PublicIPAddressUsage),
				"SubnetCountLimit":       strconv.Itoa(quota.SubnetCountLimit),
				"SubnetCountUsage":       strconv.Itoa(quota.SubnetCountUsage),
				"NetworkCountLimit":      strconv.Itoa(quota.NetworkCountLimit),
				"NetworkCountUsage":      strconv.Itoa(quota.NetworkCountUsage),
				"SecurityGroupLimit":     strconv.Itoa(quota.SecurityGroupLimit),
				"SecurityGroupUsage":     strconv.Itoa(quota.SecurityGroupUsage),
				"SecurityGroupRuleLimit": strconv.Itoa(quota.SecurityGroupRuleLimit),
				"SecurityGroupRuleUsage": strconv.Itoa(quota.SecurityGroupRuleUsage),
			}
			outputKeyValue(data)
		default:
			data := make([][]string, 11)
			appendQuotaData(data, 0, "Instances", quota.InstanceCountUsage, quota.InstanceCountLimit)
			appendQuotaData(data, 1, "CPU cores", quota.CPUCoreUsage, quota.CPUCoreLimit)
			appendQuotaData(data, 1, "CPU cores", quota.CPUCoreUsage, quota.CPUCoreLimit)
			appendQuotaData(data, 2, "RAM MB", quota.RAMMegabytesUsage, quota.RAMMegabytesLimit)
			appendQuotaData(data, 3, "Disk GB", quota.DiskGigabytesUsage, quota.DiskGigabytesLimit)
			appendQuotaData(data, 4, "Volumes", quota.DiskVolumeCountUsage, quota.DiskVolumeCountLimit)
			appendQuotaData(data, 5, "Snapshots", quota.DiskSnapshotCountUsage, quota.DiskSnapshotCountLimit)
			appendQuotaData(data, 6, "Public IPs", quota.PublicIPAddressUsage, quota.PublicIPAddressLimit)
			appendQuotaData(data, 7, "Subnets", quota.SubnetCountUsage, quota.SubnetCountLimit)
			appendQuotaData(data, 8, "Private networks", quota.NetworkCountUsage, quota.NetworkCountLimit)
			appendQuotaData(data, 9, "Firewalls", quota.SecurityGroupUsage, quota.SecurityGroupLimit)
			appendQuotaData(data, 10, "Firewall rules", quota.SecurityGroupRuleUsage, quota.SecurityGroupRuleLimit)
			outputTable([]string{"Item", "Usage", "Limit"}, data)
		}
	},
}

func appendQuotaData(data [][]string, pos int, label string, usage int, limit int) {
	if float32(usage) >= (float32(limit) * 0.8) {
		data[pos] = []string{label, aurora.Red(strconv.Itoa(usage)).String(), strconv.Itoa(limit)}
	} else {
		data[pos] = []string{label, strconv.Itoa(usage), strconv.Itoa(limit)}
	}
}

func init() {
	rootCmd.AddCommand(quotaCmd)

	quotaCmd.AddCommand(quotaShowCmd)
}
