package utility

import (
	"testing"

	"github.com/civo/civogo"
)

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

// func TestRemoveApplicationFromInstalledListWithPlan(t *testing.T) {
// 	current := []civogo.KubernetesInstalledApplication{
// 		{
// 			Name: "mysql",
// 		},
// 	}
// 	uninstall := "mysql"

// 	ret := RemoveApplicationFromInstalledList(current, uninstall)

// 	if ret != "mysql" {
// 		t.Errorf("expected 'mysql', got '%s'", ret)
// 	}
// }
