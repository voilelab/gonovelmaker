package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
)

type BackendListCmd struct {
	jsonOutput bool

	cmd *cobra.Command
}

func NewBackendListCmd() *BackendListCmd {
	listCmd := &BackendListCmd{}
	listCmd.cmd = &cobra.Command{
		Use:   "list",
		Short: "List all configured LLM backends",
		Long:  `List all configured LLM backends with their settings.`,
		RunE:  listCmd.run,
	}

	listCmd.cmd.Flags().BoolVar(&listCmd.jsonOutput, "json", false, "Output in JSON format")

	return listCmd
}

func (b *BackendListCmd) run(cmd *cobra.Command, args []string) error {
	// Load current config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.LLMBackend) == 0 {
		if !b.jsonOutput {
			fmt.Println("No backends configured.")
		} else {
			fmt.Println("[]")
		}
		return nil
	}

	// Sort backend names for consistent output
	names := cfg.GetBackendNames()
	sort.Strings(names)

	if b.jsonOutput {
		// JSON output
		output := make([]map[string]interface{}, 0, len(names))
		for _, name := range names {
			backend := cfg.GetBackend(name)
			output = append(output, map[string]interface{}{
				"name":        name,
				"type":        backend.Type,
				"api_key":     maskAPIKey(backend.APIKey),
				"base_url":    backend.BaseURL,
				"model":       backend.Model,
				"image_model": backend.ImageModel,
				"timeout":     backend.Timeout,
				"is_default":  name == cfg.UserLLMBackend,
			})
		}
		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(jsonData))
	} else {
		// Human-readable output
		fmt.Println("Configured backends:")
		fmt.Println()
		for _, name := range names {
			backend := cfg.GetBackend(name)
			defaultMarker := ""
			if name == cfg.UserLLMBackend {
				defaultMarker = " (default)"
			}
			fmt.Printf("• %s%s\n", name, defaultMarker)
			fmt.Printf("  Type:        %s\n", backend.Type)
			if backend.APIKey != "" {
				fmt.Printf("  API Key:     %s\n", maskAPIKey(backend.APIKey))
			}
			if backend.BaseURL != "" {
				fmt.Printf("  Base URL:    %s\n", backend.BaseURL)
			}
			if backend.Model != "" {
				fmt.Printf("  Model:       %s\n", backend.Model)
			}
			if backend.ImageModel != "" {
				fmt.Printf("  Image Model: %s\n", backend.ImageModel)
			}
			if backend.Timeout > 0 {
				fmt.Printf("  Timeout:     %ds\n", backend.Timeout)
			}
			fmt.Println()
		}
	}

	return nil
}

// maskAPIKey masks all but the last 4 characters of an API key
func maskAPIKey(key string) string {
	if key == "" {
		return ""
	}
	if len(key) <= 4 {
		return "***"
	}
	return "..." + key[len(key)-4:]
}
