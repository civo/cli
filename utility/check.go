package utility

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"

	"github.com/civo/civogo"
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
		if strings.Contains(v.Name, "k3s") {
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

	return false, defaultRegion.Code, nil
}

// EnsureCurrentRegion there's a current region set, error out if not
func EnsureCurrentRegion() {
	if config.Current.Meta.DefaultRegion == "" {
		Error("No region set - list available regions with \"civo region ls\" and choose a default using \"civo region current REGION\", or specify one with every command using \"--region=REGION\"\n")
		os.Exit(1)
	}
}
