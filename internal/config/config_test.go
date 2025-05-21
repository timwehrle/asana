package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/timwehrle/asana/internal/api/asana"
)

const windows = "windows"

func setupTestConfig(t *testing.T) *Config {
	t.Helper()

	cfg := &Config{}
	err := cfg.Save()
	require.NoError(t, err)
	return cfg
}

func TestConfig(t *testing.T) {
	tmpDir := t.TempDir()
	originalXDGConfig := os.Getenv(xdgConfigHome)
	t.Setenv(xdgConfigHome, tmpDir)
	defer t.Setenv(xdgConfigHome, originalXDGConfig)

	t.Run("Save and load config", func(t *testing.T) {
		cfg := &Config{
			CreatedAt: time.Now(),
			Username:  "testuser",
			Workspace: &asana.Workspace{
				ID:   "123",
				Name: "TestWorkspace",
			},
		}

		err := cfg.Save()
		require.NoError(t, err)

		newCfg := &Config{}
		err = newCfg.Load()
		require.NoError(t, err)

		assert.Equal(t, cfg.Username, newCfg.Username)
		assert.Equal(t, cfg.Workspace.ID, newCfg.Workspace.ID)
		assert.Equal(t, cfg.Workspace.Name, newCfg.Workspace.Name)
	})

	t.Run("Set Valid Field", func(t *testing.T) {
		cfg := setupTestConfig(t)
		err := cfg.Set("Username", "newuser")
		require.NoError(t, err)

		newCfg := &Config{}
		err = newCfg.Load()
		require.NoError(t, err)
		assert.Equal(t, "newuser", newCfg.Username)
	})

	t.Run("Set Invalid Field", func(t *testing.T) {
		cfg := setupTestConfig(t)
		err := cfg.Set("nonexistentfield", "value")
		require.NoError(t, err)

		newCfg := &Config{}
		err = newCfg.Load()
		require.NoError(t, err)
		assert.Empty(t, newCfg.Username)
	})

	t.Run("Concurrent Access", func(t *testing.T) {
		cfg := &Config{}
		done := make(chan bool)

		go func() {
			for i := 0; i < 100; i++ {
				if err := cfg.Save(); err != nil {
					t.Error("Save failed:", err)
				}
				if err := cfg.Load(); err != nil {
					t.Error("Load failed:", err)
				}
			}
			done <- true
		}()

		go func() {
			for i := 0; i < 100; i++ {
				if err := cfg.Set("Username", "user"+string(rune(i))); err != nil {
					t.Error("Set failed:", err)
				}
			}
			done <- true
		}()

		<-done
		<-done

		err := cfg.Load()
		require.NoError(t, err)
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
		if tt.onlyWindows && runtime.GOOS != windows {
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

func TestConfigErrors(t *testing.T) {
	t.Run("Invalid Config Format", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv(xdgConfigHome, tmpDir)

		configPath := filepath.Join(configDir(), "config.yml")
		err := os.MkdirAll(filepath.Dir(configPath), 0755)
		require.NoError(t, err)

		err = os.WriteFile(configPath, []byte("invalid: yaml: content: {["), 0600)
		require.NoError(t, err)

		cfg := &Config{}
		err = cfg.Load()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read config")
	})

	t.Run("Permission Denied", func(t *testing.T) {
		if runtime.GOOS == windows {
			t.Skip("Skipping on Windows")
		}

		tmpDir := t.TempDir()
		err := os.MkdirAll(tmpDir, 0555)
		require.NoError(t, err)

		t.Setenv(xdgConfigHome, tmpDir)
		cfg := &Config{}
		err = cfg.Save()
		require.Error(t, err)
	})
}
