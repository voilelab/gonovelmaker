package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
)

type BackendRemoveCmd struct {
	cmd *cobra.Command
}

func NewBackendRemoveCmd() *BackendRemoveCmd {
	removeCmd := &BackendRemoveCmd{}
	removeCmd.cmd = &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove an LLM backend configuration",
		Long:  `Remove an existing LLM backend configuration by name.`,
		Args:  cobra.ExactArgs(1),
		RunE:  removeCmd.run,
	}

	return removeCmd
}

func (b *BackendRemoveCmd) run(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Load current config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Remove the backend
	if !cfg.RemoveBackend(name) {
		return fmt.Errorf("backend '%s' not found", name)
	}

	// Save the config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Backend '%s' removed successfully\n", name)
	if cfg.UserLLMBackend == "" && len(cfg.LLMBackend) > 0 {
		fmt.Println("Note: No default backend set. Use 'backend use <name>' to set one.")
	}

	return nil
}
