package config

import (
	_ "embed"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

//go:embed example_config.toml
var exampleConfig string

// Config represents the user's configuration settings
type Config struct {
	OpenAIKey string `toml:"openai_api_key"`
	Model     string `toml:"model"`
	BaseURL   string `toml:"base_url"`
	Timeout   int    `toml:"timeout"`
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

		cfg.OpenAIKey = os.Getenv("OPENAI_API_KEY")
		cfg.BaseURL = os.Getenv("OPENAI_BASE_URL")
		cfg.Model = os.Getenv("OPENAI_MODEL")
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

	// Fallback to environment variables if values are empty
	if cfg.OpenAIKey == "" {
		cfg.OpenAIKey = os.Getenv("OPENAI_API_KEY")
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = os.Getenv("OPENAI_BASE_URL")
	}

	if cfg.Model == "" {
		cfg.Model = os.Getenv("OPENAI_MODEL")
	}

	return cfg, nil
}

// GetModelOrDefault returns the configured model or the default "gpt-4o"
func (c *Config) GetModelOrDefault() string {
	if c.Model == "" {
		return "gpt-4o"
	}
	return c.Model
}
