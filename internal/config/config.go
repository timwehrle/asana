package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/viper"
	"github.com/timwehrle/asana-api"
)

type Config struct {
	Username  string           `mapstructure:"username"`
	Workspace *asana.Workspace `mapstructure:"workspace"`
	mu        sync.RWMutex
}

const (
	appData       = "AppData"
	xdgConfigHome = "XDG_CONFIG_HOME"
)

type Error struct {
	Message string
}

func (e Error) Error() string {
	return e.Message
}

var errConfigFileNotFound viper.ConfigFileNotFoundError

// configDir determines the directory for storing configuration files.
func configDir() string {
	var path string

	if a := os.Getenv(xdgConfigHome); a != "" {
		path = filepath.Join(a, "asana-cli")
	} else if b := os.Getenv(appData); runtime.GOOS == "windows" && b != "" {
		path = filepath.Join(b, "Asana CLI")
	} else {
		d, _ := os.UserHomeDir()
		path = filepath.Join(d, ".config", "asana-cli")
	}
	return path
}

func initViper() error {
	configPath := configDir()
	if err := os.MkdirAll(configPath, 0750); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	return nil
}

func (c *Config) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := initViper(); err != nil {
		return err
	}

	viper.Set("username", c.Username)
	viper.Set("workspace", c.Workspace)

	if err := viper.WriteConfig(); err != nil {
		if errors.As(err, &errConfigFileNotFound) {
			return viper.SafeWriteConfig()
		}
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (c *Config) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := initViper(); err != nil {
		return err
	}

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &errConfigFileNotFound) {
			return Error{Message: heredoc.Docf(`
                No configuration file found. Please run %[1]sasana auth login%[1]s to authenticate.
            `, "`")}
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to decode confid: %w", err)
	}

	return nil
}

func (c *Config) Set(field string, value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := initViper(); err != nil {
		return err
	}

	viper.Set(field, value)

	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to update config struct: %w", err)
	}

	return viper.WriteConfig()
}
