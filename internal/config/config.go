package config

import (
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

//go:embed example_config.toml
var exampleConfig string

// LLMBackendConfig represents a single LLM backend configuration
type LLMBackendConfig struct {
	Type       string `toml:"type"` // "openai", "openrouter", etc.
	APIKey     string `toml:"api_key"`
	BaseURL    string `toml:"base_url"`
	Model      string `toml:"model"`
	ImageModel string `toml:"image_model"`
	Timeout    int    `toml:"timeout"`
}

// Config represents the user's configuration settings
type Config struct {
	// Multi-backend support
	UserLLMBackend string                      `toml:"user_llm_backend"` // Default backend to use
	LLMBackend     map[string]LLMBackendConfig `toml:"llm_backend"`      // Named backends
}

// Load reads the configuration from ~/.novelmaker/config.toml
// If the file doesn't exist, it creates an empty config file
// It falls back to environment variables if config values are empty
func Load() (*Config, error) {
	cfg := &Config{}

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Warn("unable to get home directory, using environment variables only", "error", err)
		return cfg, nil // Return empty config if we can't get home dir
	}

	// Construct config file path
	configDir := filepath.Join(homeDir, ".novelmaker")
	configPath := filepath.Join(configDir, "config.toml")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err == nil {
			os.WriteFile(configPath, []byte(exampleConfig), 0644)
		}
		return cfg, nil
	}

	// Read and parse the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetBackend returns the configuration for a specific named backend
// If name is empty, it returns the default backend specified by user_llm_backend
// Returns nil if no backend is found
func (c *Config) GetBackend(name string) *LLMBackendConfig {
	// If no name provided, use the default
	if name == "" {
		name = c.UserLLMBackend
	}

	// If we have named backends configured, try to get the requested one
	if len(c.LLMBackend) > 0 {
		if backend, ok := c.LLMBackend[name]; ok {
			return &backend
		}
		// If name is empty and we have named backends, return the first one
		if name == "" {
			for _, backend := range c.LLMBackend {
				return &backend
			}
		}
	}

	// No backend found
	return nil
}

// GetBackendNames returns a list of all configured backend names
func (c *Config) GetBackendNames() []string {
	names := make([]string, 0, len(c.LLMBackend))
	for name := range c.LLMBackend {
		names = append(names, name)
	}
	return names
}

// GetModelOrDefaultForBackend returns the model for a specific backend
func (c *Config) GetModelOrDefaultForBackend(backendName string) string {
	backend := c.GetBackend(backendName)
	if backend.Model == "" {
		return "gpt-4o"
	}
	return backend.Model
}

// GetImageModelOrDefaultForBackend returns the image model for a specific backend
func (c *Config) GetImageModelOrDefaultForBackend(backendName string) string {
	backend := c.GetBackend(backendName)
	if backend.ImageModel == "" {
		return "dall-e-3"
	}
	return backend.ImageModel
}

// Save writes the configuration back to ~/.novelmaker/config.toml
func (c *Config) Save() error {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Construct config file path
	configDir := filepath.Join(homeDir, ".novelmaker")
	configPath := filepath.Join(configDir, "config.toml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Marshal config to TOML
	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(configPath, data, 0644)
}

// AddOrUpdateBackend adds a new backend or updates an existing one
func (c *Config) AddOrUpdateBackend(name string, backend LLMBackendConfig) {
	if c.LLMBackend == nil {
		c.LLMBackend = make(map[string]LLMBackendConfig)
	}
	c.LLMBackend[name] = backend
}

// RemoveBackend removes a backend by name
func (c *Config) RemoveBackend(name string) bool {
	if c.LLMBackend == nil {
		return false
	}
	_, exists := c.LLMBackend[name]
	if exists {
		delete(c.LLMBackend, name)
		// If we removed the default backend, clear the default
		if c.UserLLMBackend == name {
			c.UserLLMBackend = ""
		}
	}
	return exists
}

// SetDefaultBackend sets the default backend to use
func (c *Config) SetDefaultBackend(name string) error {
	// Check if backend exists
	if _, ok := c.LLMBackend[name]; !ok {
		return fmt.Errorf("backend '%s' does not exist", name)
	}
	c.UserLLMBackend = name
	return nil
}
