package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
)

type UpdatePluginCmd struct {
	cmd *cobra.Command
}

func NewUpdatePluginCmd() *UpdatePluginCmd {
	updatePluginCmd := &UpdatePluginCmd{}
	updatePluginCmd.cmd = &cobra.Command{
		Use:   "update-plugin",
		Short: "Update the Obsidian Novel Maker plugin in the current vault",
		Long: `Copies the latest Novel Maker plugin files (main.js, manifest.json, styles.css) 
from the embedded plugin directory to the vault's .obsidian/plugins/obsidian-novelmaker/ directory.`,
		RunE: updatePluginCmd.run,
	}

	return updatePluginCmd
}

func (u *UpdatePluginCmd) run(cmd *cobra.Command, args []string) error {
	fmt.Println("Updating Obsidian Novel Maker plugin...")

	// Get current working directory (should be the vault root)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create vault
	vault, err := obsidian.NewVault(cwd)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}
	defer vault.Close()

	err = vault.UpdatePlugin(obsidianNovelmaker, "obsidian-novelmaker")
	if err != nil {
		return fmt.Errorf("failed to update Obsidian plugin files: %w", err)
	}

	fmt.Println("✓ Successfully updated Obsidian Novel Maker plugin!")

	return nil
}
