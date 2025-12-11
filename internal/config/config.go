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

// LLMBackendConfig represents a single LLM backend configuration
type LLMBackendConfig struct {
	Type       string `toml:"type"` // "openai", "openrouter", etc.
	APIKey     string `toml:"api_key"`
	BaseURL    string `toml:"base_url"`
	Model      string `toml:"model"`
	ImageModel string `toml:"image_model"`
}

// Config represents the user's configuration settings
type Config struct {
	// Legacy fields for backward compatibility
	OpenAIKey  string `toml:"openai_api_key"`
	Model      string `toml:"model"`
	ImageModel string `toml:"image_model"`
	BaseURL    string `toml:"base_url"`
	Timeout    int    `toml:"timeout"`

	// New multi-backend support
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

	if cfg.ImageModel == "" {
		cfg.ImageModel = os.Getenv("OPENAI_IMAGE_MODEL")
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

// GetImageModelOrDefault returns the configured image model or the default "dall-e-3"
func (c *Config) GetImageModelOrDefault() string {
	if c.ImageModel == "" {
		return "dall-e-3"
	}
	return c.ImageModel
}

// GetBackend returns the configuration for a specific named backend
// If name is empty, it returns the default backend specified by user_llm_backend
// Falls back to legacy configuration if no named backends are configured
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

	// Fall back to legacy configuration
	return &LLMBackendConfig{
		Type:       "openai",
		APIKey:     c.OpenAIKey,
		BaseURL:    c.BaseURL,
		Model:      c.GetModelOrDefault(),
		ImageModel: c.GetImageModelOrDefault(),
	}
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
