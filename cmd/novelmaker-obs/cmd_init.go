package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
	obsplugin "github.com/voilelab/gonovelmaker/obsidian-novelmaker"
)

type InitCmd struct {
	includePlugin bool
	noGit         bool

	cmd *cobra.Command
}

func NewInitCmd() *InitCmd {
	initCmd := &InitCmd{}
	initCmd.cmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a novel project structure in the current directory",
		Long: `Creates Config/, World/, and Story/ subdirectories in the current directory, 
along with sample files.`,
		RunE: initCmd.run,
	}

	initCmd.cmd.Flags().BoolVar(&initCmd.includePlugin, "include-plugin", true, "Include Obsidian plugin files in the initialization")
	initCmd.cmd.Flags().BoolVar(&initCmd.noGit, "no-git", false, "Skip git repository initialization")

	return initCmd
}

func (i *InitCmd) run(cmd *cobra.Command, args []string) error {
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

	// Initialize vault structure
	if err := vault.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize vault: %w", err)
	}

	fmt.Println("✓ Successfully initialized novel project structure!")

	if i.includePlugin {
		err = vault.UpdatePlugin(obsplugin.GetPluginFS())
		if err != nil {
			return fmt.Errorf("failed to copy Obsidian plugin files: %w", err)
		}
		fmt.Println("✓ Included Obsidian plugin files.")
	}

	if !i.noGit {
		if err := vault.InitGit(); err != nil {
			return fmt.Errorf("failed to initialize git repository: %w", err)
		}
		fmt.Println("✓ Initialized git repository.")
	}

	return nil
}
