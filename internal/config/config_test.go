package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/timwehrle/asana-go"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestConfig(t *testing.T) {
	tmpDir := t.TempDir()
	originalXDGConfig := os.Getenv(xdgConfigHome)
	os.Setenv(xdgConfigHome, tmpDir)
	defer os.Setenv(xdgConfigHome, originalXDGConfig)

	t.Run("Save and Load Config", func(t *testing.T) {
		cfg := &Config{
			Username: "testuser",
			Workspace: &asana.Workspace{
				Name: "Test Workspace",
				ID:   "1234",
			},
		}

		// Test save
		err := cfg.Save()
		assert.NoError(t, err)

		// Verify file exists
		configPath := filepath.Join(tmpDir, "asana-cli", "config.yml")
		_, err = os.Stat(configPath)
		assert.NoError(t, err)

		newCfg := &Config{}
		err = newCfg.Load()
		assert.NoError(t, err)
		assert.Equal(t, cfg.Username, newCfg.Username)
		assert.Equal(t, cfg.Workspace.Name, newCfg.Workspace.Name)
		assert.Equal(t, cfg.Workspace.ID, newCfg.Workspace.ID)
	})

	t.Run("Load Non-existent Config", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "asana-cli", "config.yml")
		os.Remove(configPath)

		cfg := &Config{}
		err := cfg.Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "No configuration file found")
	})

	t.Run("Set Valid Field", func(t *testing.T) {
		cfg := &Config{}
		err := cfg.Set("Username", "newuser")
		assert.NoError(t, err)
		assert.Equal(t, "newuser", cfg.Username)
	})

	t.Run("Set Invalid Field", func(t *testing.T) {
		cfg := &Config{}
		err := cfg.Set("Nonexistentfield", "value")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist in config")
	})

	t.Run("Set Field Type Mismatch", func(t *testing.T) {
		cfg := &Config{}
		err := cfg.Set("Username", 123)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "value type mismatch")
	})

	t.Run("Concurrent Access", func(t *testing.T) {
		cfg := &Config{}
		done := make(chan bool)

		go func() {
			for i := 0; i < 100; i++ {
				cfg.Save()
				cfg.Load()
			}
			done <- true
		}()

		go func() {
			for i := 0; i < 100; i++ {
				cfg.Set("Username", "user"+string(rune(i)))
			}
			done <- true
		}()

		<-done
		<-done

		err := cfg.Load()
		assert.NoError(t, err)
	})
}

func TestConfigDir(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		onlyWindows bool
		env         map[string]string
		output      string
	}{
		{
			name: "HOME/USERPROFILE specified",
			env: map[string]string{
				xdgConfigHome: "",
				appData:       "",
				"USERPROFILE": tempDir,
				"HOME":        tempDir,
			},
			output: filepath.Join(tempDir, ".config", "asana-cli"),
		},
		{
			name: "XDG_CONFIG_HOME specified",
			env: map[string]string{
				xdgConfigHome: tempDir,
			},
			output: filepath.Join(tempDir, "asana-cli"),
		},
		{
			name:        "AppData specified",
			onlyWindows: true,
			env: map[string]string{
				appData: tempDir,
			},
			output: filepath.Join(tempDir, "Asana CLI"),
		},
		{
			name:        "XDG_CONFIG_HOME and AppData specified",
			onlyWindows: true,
			env: map[string]string{
				xdgConfigHome: tempDir,
				appData:       tempDir,
			},
			output: filepath.Join(tempDir, "asana-cli"),
		},
	}

	for _, tt := range tests {
		if tt.onlyWindows && runtime.GOOS != "windows" {
			continue
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.env != nil {
				for k, v := range tt.env {
					t.Setenv(k, v)
				}
			}
			assert.Equal(t, tt.output, configDir())
		})
	}
}
