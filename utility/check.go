package utility

import (
	"fmt"
	"math"

	"os"
	"runtime"
	"strings"

	"github.com/civo/civogo"
	"golang.org/x/crypto/ssh"

	"github.com/civo/cli/config"
)

// CheckOS is a function to check the OS of the user
func CheckOS() string {
	os := runtime.GOOS
	var returnValue string
	switch os {
	case "windows":
		returnValue = "windows"
	case "darwin":
		returnValue = "darwin"
	case "linux":
		returnValue = "linux"
	default:
		fmt.Printf("%s.\n", os)
	}

	return returnValue
}

// CheckQuotaPercent function to check the percent of the quota
func CheckQuotaPercent(limit int, usage int) string {
	var returnText string

	calculation := float64(usage) / float64(limit) * 100
	percent := math.Round(calculation)

	switch {
	case percent >= 80 && percent < 100:
		returnText = Orange(fmt.Sprintf("%d/%d", usage, limit))
	case percent == 100:
		returnText = Red(fmt.Sprintf("%d/%d", usage, limit))
	default:
		returnText = Green(fmt.Sprintf("%d/%d", usage, limit))
	}

	return returnText
}

// GetK3sSize is a functon to return only the k3s size
func GetK3sSize() ([]string, error) {
	client, err := config.CivoAPIClient()
	if err != nil {
		return nil, err
	}

	k8sSize := []string{}
	allSize, err := client.ListInstanceSizes()
	if err != nil {
		return nil, err
	}

	for _, v := range allSize {
		if strings.Contains(v.Name, ".kube.") {
			k8sSize = append(k8sSize, v.Name)
		}
	}

	return k8sSize, nil
}

// CheckAPPName is a function to check if the app name is valid
func CheckAPPName(appName string) bool {
	client, err := config.CivoAPIClient()
	if err != nil {
		return false
	}

	allAPP, err := client.ListKubernetesMarketplaceApplications()
	if err != nil {
		return false
	}

	for _, v := range allAPP {
		if strings.Contains(appName, v.Name) {
			return true
		}
	}

	return false
}

// ListDefaultApps is a function to list the default apps in the marketplace
func ListDefaultApps() ([]string, error) {
	client, err := config.CivoAPIClient()
	if err != nil {
		return nil, err
	}

	allApps, err := client.ListKubernetesMarketplaceApplications()
	if err != nil {
		return nil, err
	}

	var defaultApps []string
	for _, v := range allApps {
		if v.Default {
			defaultApps = append(defaultApps, v.Name)
		}
	}
	return defaultApps, nil
}

// CheckAvailability is a function to check if the user can
// create Iaas and k8s cluster base on the result of region
func CheckAvailability(resource string, regionSet string) (bool, string, error) {
	var defaultRegion *civogo.Region
	client, err := config.CivoAPIClient()
	if err != nil {
		return false, "", err
	}

	if regionSet != "" {
		client.Region = regionSet
	}

	switch {
	case config.Current.Meta.DefaultRegion == "" && regionSet != "" || config.Current.Meta.DefaultRegion != "" && regionSet != "":
		defaultRegion, err = client.FindRegion(regionSet)
		if err != nil {
			return false, "", err
		}
	case config.Current.Meta.DefaultRegion != "" && regionSet == "":
		defaultRegion, err = client.FindRegion(config.Current.Meta.DefaultRegion)
		if err != nil {
			return false, "", err
		}
	default:
		defaultRegion, err = client.GetDefaultRegion()
		if err != nil {
			return false, "", err
		}
	}

	if resource == "kubernetes" {
		if defaultRegion.Features.Kubernetes && !defaultRegion.OutOfCapacity {
			return true, "", nil
		}
	}
	if resource == "instance" {
		if defaultRegion.Features.Iaas && !defaultRegion.OutOfCapacity {
			return true, "", nil
		}
	}

	if resource == "object_store" {
		if defaultRegion.Features.ObjectStore && !defaultRegion.OutOfCapacity {
			return true, "", nil
		}
	}

	if resource == "loadbalancer" {
		if defaultRegion.Features.LoadBalancer && !defaultRegion.OutOfCapacity {
			return true, "", nil
		}
	}

	if resource == "dbaas" {
		if defaultRegion.Features.DBaaS && !defaultRegion.OutOfCapacity {
			return true, "", nil
		}
	}

	if resource == "volume" {
		if defaultRegion.Features.Volume && !defaultRegion.OutOfCapacity {
			return true, "", nil
		}
	}

	if resource == "kfaas" {
		if defaultRegion.Features.KFaaS && !defaultRegion.OutOfCapacity {
			return true, "", nil
		}
	}

	return false, defaultRegion.Code, nil
}

// EnsureCurrentRegion there's a current region set, error out if not
func EnsureCurrentRegion() {
	if config.Current.Meta.DefaultRegion == "" {
		Error("No region set - list available regions with \"civo region ls\" and choose a default using \"civo region current REGION\", or specify one with every command using \"--region=REGION\"\n")
		os.Exit(1)
	}
}

// ValidNameLength is a function to check is the name is valid
func ValidNameLength(name string) bool {
	return len(name) > 63
}

// CanManageVolume is a function to check if a cluster can manage the volume
func CanManageVolume(volume *civogo.Volume) bool {
	return volume.ClusterID == ""
}

// ValidateSSHKey is a function to check if the public key is valid
func ValidateSSHKey(key []byte) error {
	_, _, _, _, err := ssh.ParseAuthorizedKey(key)
	if err != nil {
		return err
	}

	return nil
}
