package config

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestMultiBackendConfig(t *testing.T) {
	tomlData := `
user_llm_backend = "openai"

[llm_backend.openai]
type = "openai"
api_key = "sk-test-openai"
base_url = "https://api.openai.com/v1"
model = "gpt-4o"
image_model = "dall-e-3"

[llm_backend.openrouter]
type = "openai"
api_key = "sk-test-openrouter"
base_url = "https://openrouter.ai/api/v1"
model = "anthropic/claude-3.5-sonnet"
image_model = "openai/dall-e-3"
`

	var cfg Config
	err := toml.Unmarshal([]byte(tomlData), &cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Test user_llm_backend
	if cfg.UserLLMBackend != "openai" {
		t.Errorf("Expected user_llm_backend to be 'openai', got '%s'", cfg.UserLLMBackend)
	}

	// Test backend retrieval by name
	openaiBackend := cfg.GetBackend("openai")
	if openaiBackend == nil {
		t.Fatal("Expected openai backend to exist")
	}
	if openaiBackend.APIKey != "sk-test-openai" {
		t.Errorf("Expected openai api_key to be 'sk-test-openai', got '%s'", openaiBackend.APIKey)
	}
	if openaiBackend.Model != "gpt-4o" {
		t.Errorf("Expected openai model to be 'gpt-4o', got '%s'", openaiBackend.Model)
	}

	openrouterBackend := cfg.GetBackend("openrouter")
	if openrouterBackend == nil {
		t.Fatal("Expected openrouter backend to exist")
	}
	if openrouterBackend.APIKey != "sk-test-openrouter" {
		t.Errorf("Expected openrouter api_key to be 'sk-test-openrouter', got '%s'", openrouterBackend.APIKey)
	}
	if openrouterBackend.Model != "anthropic/claude-3.5-sonnet" {
		t.Errorf("Expected openrouter model to be 'anthropic/claude-3.5-sonnet', got '%s'", openrouterBackend.Model)
	}

	// Test default backend retrieval
	defaultBackend := cfg.GetBackend("")
	if defaultBackend.APIKey != "sk-test-openai" {
		t.Errorf("Expected default backend to be openai with api_key 'sk-test-openai', got '%s'", defaultBackend.APIKey)
	}

	// Test GetBackendNames
	names := cfg.GetBackendNames()
	if len(names) != 2 {
		t.Errorf("Expected 2 backend names, got %d", len(names))
	}
}

func TestLegacyConfig(t *testing.T) {
	tomlData := `
openai_api_key = "sk-legacy"
model = "gpt-4"
image_model = "dall-e-2"
base_url = "https://api.example.com/v1"
timeout = 120
`

	var cfg Config
	err := toml.Unmarshal([]byte(tomlData), &cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Test legacy fields
	if cfg.OpenAIKey != "sk-legacy" {
		t.Errorf("Expected openai_api_key to be 'sk-legacy', got '%s'", cfg.OpenAIKey)
	}
	if cfg.Model != "gpt-4" {
		t.Errorf("Expected model to be 'gpt-4', got '%s'", cfg.Model)
	}

	// Test that GetBackend falls back to legacy config
	backend := cfg.GetBackend("")
	if backend.APIKey != "sk-legacy" {
		t.Errorf("Expected backend api_key to be 'sk-legacy', got '%s'", backend.APIKey)
	}
	if backend.Model != "gpt-4" {
		t.Errorf("Expected backend model to be 'gpt-4', got '%s'", backend.Model)
	}
	if backend.BaseURL != "https://api.example.com/v1" {
		t.Errorf("Expected backend base_url to be 'https://api.example.com/v1', got '%s'", backend.BaseURL)
	}
}

func TestMixedConfig(t *testing.T) {
	// Test that new config takes precedence over legacy
	tomlData := `
openai_api_key = "sk-legacy"
model = "gpt-4"

user_llm_backend = "custom"

[llm_backend.custom]
type = "openai"
api_key = "sk-custom"
model = "gpt-4o-mini"
`

	var cfg Config
	err := toml.Unmarshal([]byte(tomlData), &cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Test that GetBackend returns the new config when available
	backend := cfg.GetBackend("")
	if backend.APIKey != "sk-custom" {
		t.Errorf("Expected backend api_key to be 'sk-custom', got '%s'", backend.APIKey)
	}
	if backend.Model != "gpt-4o-mini" {
		t.Errorf("Expected backend model to be 'gpt-4o-mini', got '%s'", backend.Model)
	}
}

func TestGetModelOrDefaultForBackend(t *testing.T) {
	tomlData := `
user_llm_backend = "openai"

[llm_backend.openai]
type = "openai"
model = "gpt-4o"

[llm_backend.nomodel]
type = "openai"
`

	var cfg Config
	err := toml.Unmarshal([]byte(tomlData), &cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Test with explicit model
	model := cfg.GetModelOrDefaultForBackend("openai")
	if model != "gpt-4o" {
		t.Errorf("Expected model 'gpt-4o', got '%s'", model)
	}

	// Test with default fallback
	defaultModel := cfg.GetModelOrDefaultForBackend("nomodel")
	if defaultModel != "gpt-4o" {
		t.Errorf("Expected default model 'gpt-4o', got '%s'", defaultModel)
	}
}
