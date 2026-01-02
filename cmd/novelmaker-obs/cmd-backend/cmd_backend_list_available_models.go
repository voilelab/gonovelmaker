package cmdbackend

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
	"github.com/voilelab/gonovelmaker/internal/llmbackend"
)

type BackendListAvailableModelsCmd struct {
	timeout    int
	jsonOutput bool

	llmbackendMaker llmbackend.LLMBackendMaker

	cmd *cobra.Command
}

func NewBackendListAvailableModelsCmd(llmbackendMaker llmbackend.LLMBackendMaker) *BackendListAvailableModelsCmd {
	listModelsCmd := &BackendListAvailableModelsCmd{
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

	return listModelsCmd
}

func (b *BackendListAvailableModelsCmd) run(cmd *cobra.Command, args []string) error {
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

	// Fetch models
	if !b.jsonOutput {
		fmt.Print("Fetching models... ")
	}

	models, err := llmBackend.ListModels(ctx)
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

	// Sort models for consistent output
	sort.Strings(models)

	if b.jsonOutput {
		jsonData, err := json.MarshalIndent(map[string]any{
			"success": true,
			"backend": name,
			"count":   len(models),
			"models":  models,
		}, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(jsonData))
	} else {
		fmt.Println("✅ SUCCESS")
		fmt.Println()
		fmt.Printf("Available models (%d):\n", len(models))
		for _, model := range models {
			fmt.Printf("  • %s\n", model)
		}
	}

	return nil
}
