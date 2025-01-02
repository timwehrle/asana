package config

import (
	"bitbucket.org/mikehouston/asana-go"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Username  string           `yaml:"username"`
	Workspace *asana.Workspace `yaml:"workspace"`
}

type DefaultWorkspace struct {
	ID   string `yaml:"gid"`
	Name string `yaml:"name"`
}

type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func getConfigFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir: %w", err)
	}

	configPath := filepath.Join(configDir, "asana")
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create config dir: %w", err)
	}

	return filepath.Join(configPath, "config.yml"), nil
}

func SaveConfig(config Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	err = encoder.Encode(&config)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfig() (config Config, err error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, Error{Message: heredoc.Docf(`No configuration file found. Please run %[1]sasana auth login%[1]s to authenticate.`, "`")}
		}
		return Config{}, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func UpdateDefaultWorkspace(gid, name string) error {
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	config.Workspace = &asana.Workspace{
		ID:   gid,
		Name: name,
	}

	if err := SaveConfig(config); err != nil {
		return fmt.Errorf("failed to save updated configuration: %w", err)
	}

	return nil
}

func GetDefaultWorkspace() (*asana.Workspace, error) {
	config, err := LoadConfig()
	if err != nil {
		return &asana.Workspace{}, err
	}

	return config.Workspace, nil
}
