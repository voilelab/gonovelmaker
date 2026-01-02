package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cmdbackend "github.com/voilelab/gonovelmaker/cmd/novelmaker-obs/cmd-backend"
	"github.com/voilelab/gonovelmaker/internal/config"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "novelmaker-obs",
		Short: "A CLI tool for generating novels in Obsidian vaults",
		Long: `novelmaker-obs is a CLI tool that works with Obsidian vaults to manage 
novel projects, worldbooks, and chapters. It uses OpenAI API to generate new chapters 
based on your existing content.`,
	}

	initCmd := NewInitCmd()
	scanCmd := NewScanCmd()
	genNextCmd := NewGenNextCmd(openAIBackendMaker)
	genNextEmptyCmd := NewGenNextEmptyCmd()
	genCurrCmd := NewGenCurrCmd(openAIBackendMaker)
	genCharCmd := NewGenCharCmd(openAIBackendMaker)
	genCharCurrCmd := NewGenCharCurrCmd(openAIBackendMaker)
	genCharImgCmd := NewGenCharImgCmd(openAIBackendMaker)
	exportCmd := NewExportCmd()
	updatePluginCmd := NewUpdatePluginCmd()
	configTestCmd := NewConfigCheckCmd()

	backendCmd := cmdbackend.NewBackendCmd(openAIBackendMaker)

	rootCmd.AddCommand(initCmd.cmd)
	rootCmd.AddCommand(scanCmd.cmd)
	rootCmd.AddCommand(genNextCmd.cmd)
	rootCmd.AddCommand(genNextEmptyCmd.cmd)
	rootCmd.AddCommand(genCurrCmd.cmd)
	rootCmd.AddCommand(genCharCmd.cmd)
	rootCmd.AddCommand(genCharCurrCmd.cmd)
	rootCmd.AddCommand(genCharImgCmd.cmd)
	rootCmd.AddCommand(exportCmd.cmd)
	rootCmd.AddCommand(updatePluginCmd.cmd)
	rootCmd.AddCommand(configTestCmd.cmd)
	rootCmd.AddCommand(backendCmd.Command())

	// Initialize config (creates empty config file if it doesn't exist)
	config.Load()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
