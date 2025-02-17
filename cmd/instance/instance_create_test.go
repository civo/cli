package instance_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCreateInstanceBasic(t *testing.T) {
	instanceName := "test-basic-instance"

	defer func() { // Cleanup
		deleteCmd := exec.Command("civo", "instance", "delete", instanceName, "--region=lon1", "--yes")
		deleteCmd.CombinedOutput()
	}()

	// Create a basic instance
	cmd := exec.Command("civo", "instance", "create", instanceName, "--size=g3.xsmall", "--region=lon1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to create instance: %v\nOutput: %s", err, string(output))
	}

	// Verify instance creation
	if !strings.Contains(string(output), instanceName) {
		t.Errorf("Expected instance name '%s' in output, got: %s", instanceName, string(output))
	}

	// Verify `civo instance list` shows the instance
	listCmd := exec.Command("civo", "instance", "list", "--region=lon1")
	listOutput, listErr := listCmd.CombinedOutput()
	if listErr != nil {
		t.Fatalf("Failed to list instances: %v\nOutput: %s", listErr, string(listOutput))
	}
	if !strings.Contains(string(listOutput), instanceName) {
		t.Errorf("Instance '%s' not found in list output: %s", instanceName, string(listOutput))
	}
}

func TestCreateInstanceInvalidSize(t *testing.T) {
	instanceName := "test-invalid-size"

	cmd := exec.Command("civo", "instance", "create", instanceName, "--size=invalid-size", "--region=lon1")
	output, err := cmd.CombinedOutput()

	// Ensure command fails
	if err == nil {
		t.Fatalf("Expected failure for invalid size, but command succeeded")
	}

	// Verify output mentions invalid size
	if !strings.Contains(string(output), "The provided size is not valid") {
		t.Errorf("Expected error message about invalid size, got: %s", string(output))
	}
}

func TestCreateInstanceWithScript(t *testing.T) {
	instanceName := "test-instance-script"
	scriptFile := "./test_script.sh"

	// Create a script file for testing
	scriptContent := "#!/bin/bash\necho Hello, World!"
	if err := os.WriteFile(scriptFile, []byte(scriptContent), 0644); err != nil {
		t.Fatalf("Failed to create test script file: %v", err)
	}
	defer os.Remove(scriptFile) // Cleanup script file

	defer func() { // Cleanup instance
		deleteCmd := exec.Command("civo", "instance", "delete", instanceName, "--region=lon1", "--yes")
		deleteCmd.CombinedOutput()
	}()

	// Create an instance with a script
	cmd := exec.Command("civo", "instance", "create", instanceName, "--size=g3.xsmall", "--region=lon1", "--script="+scriptFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to create instance with script: %v\nOutput: %s", err, string(output))
	}

	// Verify instance creation
	if !strings.Contains(string(output), instanceName) {
		t.Errorf("Expected instance name '%s' in output, got: %s", instanceName, string(output))
	}
}

func TestCreateInstanceWithJSONOutput(t *testing.T) {
	instanceName := "test-json-instance"

	defer func() { // Cleanup
		deleteCmd := exec.Command("civo", "instance", "delete", instanceName, "--region=lon1", "--yes")
		deleteCmd.CombinedOutput()
	}()

	// Create an instance with JSON output
	cmd := exec.Command("civo", "instance", "create", instanceName, "--size=g3.xsmall", "--region=lon1", "--output=json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to create instance with JSON output: %v\nOutput: %s", err, string(output))
	}

	// Parse and verify JSON output
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Check the instance name in JSON
	if result["hostname"] != instanceName {
		t.Errorf("Expected hostname '%s' in JSON output, got: %s", instanceName, result["hostname"])
	}
}
