package workspace

import (
	"encoding/json"
	"os"
	"testing"
)

func setupTestDirectory(t *testing.T) string {
	t.Helper()

	tempDir := t.TempDir()
	err := os.Setenv("XDG_CONFIG_HOME", tempDir)
	if err != nil {
		t.Fatalf("Failed to set XDG_CONFIG_HOME: %v", err)
	}
	return tempDir
}

func TestSaveAndLoadDefaultWorkspace(t *testing.T) {
	setupTestDirectory(t)

	gid := "12345"
	name := "Test Workspace"

	err := SaveDefaultWorkspace(gid, name)
	if err != nil {
		t.Fatalf("SaveDefaultWorkspace failed: %v", err)
	}

	configPath, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("getConfigFilePath failed: %v", err)
	}
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		t.Fatalf("Config file not created at %s", configPath)
	}

	fileContents, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	err = json.Unmarshal(fileContents, &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal config file: %v", err)
	}

	if config.WorkspaceGID != gid || config.WorkspaceName != name {
		t.Errorf("Config mismatch: got %+v, want %+v", config, Config{WorkspaceGID: gid, WorkspaceName: name})
	}

	loadedGID, loadedName, err := LoadDefaultWorkspace()
	if err != nil {
		t.Fatalf("LoadDefaultWorkspace failed: %v", err)
	}
	if loadedGID != gid || loadedName != name {
		t.Errorf("LoadDefaultWorkspace mismatch: got (%s, %s), want (%s, %s)", loadedGID, loadedName, gid, name)
	}
}

func TestLoadDefaultWorkspaceNoConfig(t *testing.T) {
	setupTestDirectory(t)

	_, _, err := LoadDefaultWorkspace()
	if err == nil || err.Error() != "no default workspace set" {
		t.Errorf("Expected 'no default workspace set' error, got: %v", err)
	}
}

func TestSaveDefaultWorkspaceFailure(t *testing.T) {
	setupTestDirectory(t)

	err := os.Setenv("XDG_CONFIG_HOME", "/nonexistent-path")
	if err != nil {
		t.Fatalf("Failed to set XDG_CONFIG_HOME: %v", err)
	}

	err = SaveDefaultWorkspace("12345", "Test Workspace")
	if err == nil {
		t.Fatal("Expected SaveDefaultWorkspace to fail, but it succeeded")
	}
}
