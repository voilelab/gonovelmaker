package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
)

type BackendAddCmd struct {
	backendType string
	apiKey      string
	baseURL     string
	model       string
	imageModel  string
	timeout     int

	cmd *cobra.Command
}

func NewBackendAddCmd() *BackendAddCmd {
	addCmd := &BackendAddCmd{}
	addCmd.cmd = &cobra.Command{
		Use:   "add <name>",
		Short: "Add or edit an LLM backend configuration",
		Long:  `Add a new LLM backend configuration or update an existing one.`,
		Args:  cobra.ExactArgs(1),
		RunE:  addCmd.run,
	}

	addCmd.cmd.Flags().StringVar(&addCmd.backendType, "type", "openai", "Backend type (openai, openrouter, etc.)")
	addCmd.cmd.Flags().StringVar(&addCmd.apiKey, "api_key", "", "API key for the backend")
	addCmd.cmd.Flags().StringVar(&addCmd.baseURL, "base_url", "", "Base URL for the backend API")
	addCmd.cmd.Flags().StringVar(&addCmd.model, "model", "", "Model name to use")
	addCmd.cmd.Flags().StringVar(&addCmd.imageModel, "image_model", "", "Image model name to use")
	addCmd.cmd.Flags().IntVar(&addCmd.timeout, "timeout", 0, "Request timeout in seconds")

	return addCmd
}

func (b *BackendAddCmd) run(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Load current config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if backend already exists
	existing := cfg.GetBackend(name)
	isUpdate := existing != nil

	// Create new backend config (preserve existing values if not specified)
	backendCfg := config.LLMBackendConfig{
		Type:       b.backendType,
		APIKey:     b.apiKey,
		BaseURL:    b.baseURL,
		Model:      b.model,
		ImageModel: b.imageModel,
		Timeout:    b.timeout,
	}

	// If updating, merge with existing values
	if isUpdate {
		if !cmd.Flags().Changed("type") {
			backendCfg.Type = existing.Type
		}
		if !cmd.Flags().Changed("api_key") {
			backendCfg.APIKey = existing.APIKey
		}
		if !cmd.Flags().Changed("base_url") {
			backendCfg.BaseURL = existing.BaseURL
		}
		if !cmd.Flags().Changed("model") {
			backendCfg.Model = existing.Model
		}
		if !cmd.Flags().Changed("image_model") {
			backendCfg.ImageModel = existing.ImageModel
		}
		if !cmd.Flags().Changed("timeout") {
			backendCfg.Timeout = existing.Timeout
		}
	}

	// Add or update the backend
	cfg.AddOrUpdateBackend(name, backendCfg)

	// Save the config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	if isUpdate {
		fmt.Printf("✓ Backend '%s' updated successfully\n", name)
	} else {
		fmt.Printf("✓ Backend '%s' added successfully\n", name)
	}

	return nil
}
