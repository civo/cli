package utility

import (
	"testing"

	"github.com/civo/civogo"
)

// --- checkAppPlan tests ---

func TestCheckAppPlanFindsAppByName(t *testing.T) {
	appList := []civogo.KubernetesMarketplaceApplication{
		{Name: "redis"},
		{Name: "mysql"},
		{Name: "postgresql"},
	}

	idx, result, err := checkAppPlan(appList, "mysql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "mysql" {
		t.Errorf("expected 'mysql', got '%s'", result)
	}
	// Index must point to the matched entry
	if appList[idx].Name != "mysql" {
		t.Errorf("expected appList[idx].Name = 'mysql', got '%s'", appList[idx].Name)
	}
}

func TestCheckAppPlanFindsAppWithValidPlan(t *testing.T) {
	appList := []civogo.KubernetesMarketplaceApplication{
		{
			Name: "mysql",
			Plans: []civogo.KubernetesMarketplacePlan{
				{Label: "5GB"},
				{Label: "10GB"},
			},
		},
	}

	_, result, err := checkAppPlan(appList, "mysql:10GB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "mysql:10GB" {
		t.Errorf("expected 'mysql:10GB', got '%s'", result)
	}
}

func TestCheckAppPlanReturnsDefaultForInvalidPlan(t *testing.T) {
	appList := []civogo.KubernetesMarketplaceApplication{
		{
			Name: "mysql",
			Plans: []civogo.KubernetesMarketplacePlan{
				{Label: "5GB"},
				{Label: "10GB"},
			},
		},
	}

	_, result, err := checkAppPlan(appList, "mysql:999GB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Invalid plan falls back to default (first plan)
	if result != "mysql:5GB" {
		t.Errorf("expected 'mysql:5GB', got '%s'", result)
	}
}

func TestCheckAppPlanNoPlanSpecifiedWithPlansAvailable(t *testing.T) {
	appList := []civogo.KubernetesMarketplaceApplication{
		{
			Name: "mysql",
			Plans: []civogo.KubernetesMarketplacePlan{
				{Label: "5GB"},
			},
		},
	}

	_, result, err := checkAppPlan(appList, "mysql")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// No plan specified and plan is empty string — doesn't match any plan,
	// so it falls back to default
	if result != "mysql:5GB" {
		t.Errorf("expected 'mysql:5GB', got '%s'", result)
	}
}

func TestCheckAppPlanNoPlanSpecifiedNoPlansAvailable(t *testing.T) {
	appList := []civogo.KubernetesMarketplaceApplication{
		{Name: "traefik"},
	}

	_, result, err := checkAppPlan(appList, "traefik")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "traefik" {
		t.Errorf("expected 'traefik', got '%s'", result)
	}
}

func TestCheckAppPlanPartialNameMatch(t *testing.T) {
	appList := []civogo.KubernetesMarketplaceApplication{
		{Name: "metrics-server"},
		{Name: "mysql"},
	}

	idx, result, err := checkAppPlan(appList, "metrics")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// String returned is verbatim user input (no plans on this app)
	if result != "metrics" {
		t.Errorf("expected 'metrics', got '%s'", result)
	}
	// But the index points to the CANONICAL app — metrics-server
	if appList[idx].Name != "metrics-server" {
		t.Errorf("expected appList[idx].Name = 'metrics-server', got '%s'", appList[idx].Name)
	}
}

// --- RequestedSplit tests ---

func TestRequestedSplitSingleApp(t *testing.T) {
	appList := []civogo.KubernetesMarketplaceApplication{
		{Name: "redis"},
		{Name: "mysql"},
	}

	result := RequestedSplit(appList, "mysql", "k3s")
	if result != "mysql" {
		t.Errorf("expected 'mysql', got '%s'", result)
	}
}

func TestRequestedSplitMultipleApps(t *testing.T) {
	appList := []civogo.KubernetesMarketplaceApplication{
		{Name: "redis"},
		{Name: "mysql"},
		{Name: "postgresql"},
	}

	result := RequestedSplit(appList, "redis,mysql", "k3s")
	if result != "redis,mysql" {
		t.Errorf("expected 'redis,mysql', got '%s'", result)
	}
}

func TestRequestedSplitAppWithValidPlan(t *testing.T) {
	appList := []civogo.KubernetesMarketplaceApplication{
		{
			Name: "mysql",
			Plans: []civogo.KubernetesMarketplacePlan{
				{Label: "5GB"},
				{Label: "10GB"},
			},
		},
	}

	result := RequestedSplit(appList, "mysql:10GB", "k3s")
	if result != "mysql:10GB" {
		t.Errorf("expected 'mysql:10GB', got '%s'", result)
	}
}

func TestRequestedSplitAppWithInvalidPlanFallsBackToDefault(t *testing.T) {
	appList := []civogo.KubernetesMarketplaceApplication{
		{
			Name: "mysql",
			Plans: []civogo.KubernetesMarketplacePlan{
				{Label: "5GB"},
				{Label: "10GB"},
			},
		},
	}

	result := RequestedSplit(appList, "mysql:999GB", "k3s")
	expected := "mysql:5GB"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestRequestedSplitAppWithoutPlanAndPlansAvailable(t *testing.T) {
	appList := []civogo.KubernetesMarketplaceApplication{
		{
			Name: "mysql",
			Plans: []civogo.KubernetesMarketplacePlan{
				{Label: "5GB"},
			},
		},
	}

	// When no plan is specified and app has plans, find("") returns false,
	// so checkAppPlan falls back to the default (first) plan.
	result := RequestedSplit(appList, "mysql", "k3s")
	if result != "mysql:5GB" {
		t.Errorf("expected 'mysql:5GB', got '%s'", result)
	}
}

func TestRequestedSplitMultipleAppsWithPlans(t *testing.T) {
	appList := []civogo.KubernetesMarketplaceApplication{
		{Name: "redis"},
		{
			Name: "mysql",
			Plans: []civogo.KubernetesMarketplacePlan{
				{Label: "5GB"},
				{Label: "10GB"},
			},
		},
	}

	result := RequestedSplit(appList, "redis,mysql:10GB", "k3s")
	if result != "redis,mysql:10GB" {
		t.Errorf("expected 'redis,mysql:10GB', got '%s'", result)
	}
}

// NOTE: The Talos + metrics-server blocking path calls os.Exit(1) and
// cannot be unit tested directly. Manual testing is required for that path.
func TestRequestedSplitNonTalosAllowsMetricsServer(t *testing.T) {
	appList := []civogo.KubernetesMarketplaceApplication{
		{Name: "metrics-server"},
	}

	result := RequestedSplit(appList, "metrics-server", "k3s")
	if result != "metrics-server" {
		t.Errorf("expected 'metrics-server', got '%s'", result)
	}
}

// --- RemoveApplicationFromInstalledList tests ---

func TestRemoveApplicationFromInstalledListSimpleName(t *testing.T) {
	current := []civogo.KubernetesInstalledApplication{
		{
			Name: "mysql",
		},
	}
	uninstall := "mysql"

	ret := RemoveApplicationFromInstalledList(current, uninstall)

	if ret != "" {
		t.Errorf("expected '', got '%s'", ret)
	}
}

func TestRemoveApplicationFromInstalledListMissing(t *testing.T) {
	current := []civogo.KubernetesInstalledApplication{
		{
			Name: "mysql",
		},
	}
	uninstall := "postgresql"

	ret := RemoveApplicationFromInstalledList(current, uninstall)

	if ret != "mysql" {
		t.Errorf("expected 'mysql', got '%s'", ret)
	}
}

func TestRemoveApplicationFromInstalledListWithMultiple(t *testing.T) {
	current := []civogo.KubernetesInstalledApplication{
		{
			Name: "mysql",
		},
		{
			Name: "postgresql",
		},
		{
			Name: "redis",
		},
	}
	uninstall := "postgresql,mysql"

	ret := RemoveApplicationFromInstalledList(current, uninstall)

	if ret != "redis" {
		t.Errorf("expected 'redis', got '%s'", ret)
	}
}
