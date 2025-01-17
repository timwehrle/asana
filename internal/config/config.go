package config

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/timwehrle/asana-api"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
)

type Config struct {
	Username  string           `yaml:"username"`
	Workspace *asana.Workspace `yaml:"workspace"`
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

// configDir determines the directory for storing configuration files
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

func getConfigFilePath() (string, error) {
	configPath := configDir()

	if err := os.MkdirAll(configPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create config dir: %w", err)
	}

	return filepath.Join(configPath, "config.yml"), nil
}

func (c *Config) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	err = encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

func (c *Config) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Error{Message: heredoc.Docf(`
				No configuration file found.
				Please run %[1]sasana auth login%[1]s to authenticate.
			`, "`")}
		}
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(c); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}

	return nil
}

func (c *Config) Set(field string, value interface{}) error {
	rv := reflect.ValueOf(c).Elem()
	f := rv.FieldByName(field)

	if !f.IsValid() {
		return fmt.Errorf("field '%s' does not exist in config", field)
	}

	if !f.CanSet() {
		return fmt.Errorf("field '%s' cannot be set", field)
	}

	fv := reflect.ValueOf(value)
	if f.Kind() != fv.Kind() {
		return fmt.Errorf("value type mismatch for field '%s'", field)
	}

	f.Set(fv)
	return c.Save()
}
