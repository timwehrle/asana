package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type WorkspaceConfig struct {
	DefaultWorkspace string `json:"default_workspace"`
}

func getConfigFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir: %w", err)
	}

	configPath := filepath.Join(configDir, "act")
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create config dir: %w", err)
	}

	return filepath.Join(configPath, "act_config.json"), nil
}

func SaveDefaultWorkspace(workspaceGID string) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	config := WorkspaceConfig{DefaultWorkspace: workspaceGID}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(config)
}

func LoadDefaultWorkspace() (string, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return "", err
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("no default workspace set")
		}
		return "", fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config WorkspaceConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return "", fmt.Errorf("failed to decode config file: %w", err)
	}

	return config.DefaultWorkspace, nil
}
