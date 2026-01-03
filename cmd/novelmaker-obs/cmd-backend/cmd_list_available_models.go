package cmdbackend

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
	"github.com/voilelab/gonovelmaker/internal/llmbackend"
)

const (
	modelsCacheTTL = 86400 * time.Second // 24 hours
)

type modelsCache struct {
	Backend   string    `json:"backend"`
	Timestamp time.Time `json:"timestamp"`
	Models    []string  `json:"models"`
}

type listAvailableModelsCmd struct {
	timeout      int
	jsonOutput   bool
	noCache      bool
	refreshCache bool

	llmbackendMaker llmbackend.LLMBackendMaker

	cmd *cobra.Command
}

func newListAvailableModelsCmd(llmbackendMaker llmbackend.LLMBackendMaker) *listAvailableModelsCmd {
	listModelsCmd := &listAvailableModelsCmd{
		llmbackendMaker: llmbackendMaker,
	}

	listModelsCmd.cmd = &cobra.Command{
		Use:   "list-available-models <backend>",
		Short: "List available models for a backend",
		Long:  `Query the backend API to list all available models.`,
		Args:  cobra.ExactArgs(1),
		RunE:  listModelsCmd.run,
	}

	listModelsCmd.cmd.Flags().IntVar(&listModelsCmd.timeout, "timeout", 30, "Request timeout in seconds")
	listModelsCmd.cmd.Flags().BoolVar(&listModelsCmd.jsonOutput, "json", false, "Output in JSON format")
	listModelsCmd.cmd.Flags().BoolVar(&listModelsCmd.noCache, "no-cache", false, "Skip cache and fetch fresh data")
	listModelsCmd.cmd.Flags().BoolVar(&listModelsCmd.refreshCache, "refresh", false, "Refresh the cache with fresh data")

	return listModelsCmd
}

func (b *listAvailableModelsCmd) run(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Load current config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get the backend
	backend := cfg.GetBackend(name)
	if backend == nil {
		return fmt.Errorf("backend '%s' not found", name)
	}

	if !b.jsonOutput {
		fmt.Printf("Fetching available models from backend '%s'...\n", name)
		fmt.Printf("  Type:     %s\n", backend.Type)
		if backend.BaseURL != "" {
			fmt.Printf("  Base URL: %s\n", backend.BaseURL)
		}
		fmt.Println()
	}

	// Create the backend client
	var llmBackend llmbackend.LLMBackend
	switch backend.Type {
	case "openai", "openrouter":
		llmBackend = b.llmbackendMaker(backend.APIKey, backend.BaseURL, backend.Model)
	default:
		return fmt.Errorf("unsupported backend type: %s", backend.Type)
	}

	// Set timeout
	timeout := time.Duration(b.timeout) * time.Second
	if backend.Timeout > 0 {
		timeout = time.Duration(backend.Timeout) * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var models []string
	var fromCache bool

	// Try to load from cache if not explicitly disabled or refreshing
	if !b.noCache && !b.refreshCache {
		cachedModels, err := b.loadCache(name)
		if err == nil && cachedModels != nil {
			models = cachedModels
			fromCache = true
			if !b.jsonOutput {
				fmt.Println("✅ (from cache)")
			}
		}
	}

	// Fetch models if not from cache
	if models == nil {
		if !b.jsonOutput {
			fmt.Print("Fetching models... ")
		}

		models, err = llmBackend.ListModels(ctx)
		if err != nil {
			if b.jsonOutput {
				jsonData, _ := json.MarshalIndent(map[string]any{
					"success": false,
					"error":   err.Error(),
					"backend": name,
				}, "", "  ")
				fmt.Println(string(jsonData))
			} else {
				fmt.Println("❌ FAILED")
			}
			return fmt.Errorf("failed to list models: %w", err)
		}

		if !b.jsonOutput {
			fmt.Println("✅ SUCCESS")
		}

		// Save to cache if not using --no-cache
		if !b.noCache {
			if err := b.saveCache(name, models); err != nil {
				// Log but don't fail on cache save errors
				if !b.jsonOutput {
					fmt.Printf("Warning: Failed to save cache: %v\n", err)
				}
			}
		}
		fromCache = false
	}

	// Sort models for consistent output
	sort.Strings(models)

	if b.jsonOutput {
		jsonData, err := json.MarshalIndent(map[string]any{
			"success":    true,
			"backend":    name,
			"count":      len(models),
			"models":     models,
			"from_cache": fromCache,
		}, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(jsonData))
	} else {
		if !fromCache {
			fmt.Println()
		}
		cacheNote := ""
		if fromCache {
			cacheNote = " (cached)"
		}
		fmt.Printf("Available models%s (%d):\n", cacheNote, len(models))
		for _, model := range models {
			fmt.Printf("  • %s\n", model)
		}
	}

	return nil
}

func (b *listAvailableModelsCmd) getCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	cacheDir := filepath.Join(homeDir, ".novelmaker", "cache")
	return cacheDir, nil
}

func (b *listAvailableModelsCmd) getCachePath(backendName string) (string, error) {
	cacheDir, err := b.getCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, fmt.Sprintf("models-%s.json", backendName)), nil
}

func (b *listAvailableModelsCmd) loadCache(backendName string) ([]string, error) {
	cachePath, err := b.getCachePath(backendName)
	if err != nil {
		return nil, err
	}

	// Check if cache file exists
	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No cache, not an error
		}
		return nil, fmt.Errorf("failed to read cache: %w", err)
	}

	var cache modelsCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to parse cache: %w", err)
	}

	// Check if cache is still valid (within TTL)
	if time.Since(cache.Timestamp) > modelsCacheTTL {
		return nil, nil // Cache expired, not an error
	}

	return cache.Models, nil
}

func (b *listAvailableModelsCmd) saveCache(backendName string, models []string) error {
	cachePath, err := b.getCachePath(backendName)
	if err != nil {
		return err
	}

	// Ensure cache directory exists
	cacheDir := filepath.Dir(cachePath)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	cache := modelsCache{
		Backend:   backendName,
		Timestamp: time.Now(),
		Models:    models,
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}

	return nil
}
