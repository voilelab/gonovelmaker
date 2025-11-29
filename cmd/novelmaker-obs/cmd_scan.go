package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
)

type ScanCmd struct {
	json bool

	cmd *cobra.Command
}

func NewScanCmd() *ScanCmd {
	scanCmd := &ScanCmd{}
	scanCmd.cmd = &cobra.Command{
		Use:   "scan",
		Short: "Scan and display the current novel project structure",
		Long: `Scans the project directory structure (Config/, World/, Story/) and displays 
the project configuration, worldbook entries, and chapters.`,
		RunE: scanCmd.run,
	}
	scanCmd.cmd.Flags().BoolVarP(&scanCmd.json, "json", "j", false, "Output in JSON format")
	return scanCmd
}

func (s *ScanCmd) run(cmd *cobra.Command, args []string) error {
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

	// Load project
	project, err := vault.LoadProject()
	if err != nil {
		return fmt.Errorf("failed to load project: %w", err)
	}

	// Load worldbooks
	worldbooks, err := vault.LoadWorldbooks()
	if err != nil {
		return fmt.Errorf("failed to load worldbooks: %w", err)
	}

	// Load characters
	characters, err := vault.LoadCharacters()
	if err != nil {
		return fmt.Errorf("failed to load characters: %w", err)
	}

	// Load chapters
	chapters, err := vault.LoadChapters()
	if err != nil {
		return fmt.Errorf("failed to load chapters: %w", err)
	}

	if s.json {
		data := map[string]any{
			"project":    project,
			"worldbooks": worldbooks,
			"characters": characters,
			"chapters":   chapters,
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")

		return encoder.Encode(data)
	}

	// Pretty print
	fmt.Println("=== PROJECT ===")
	fmt.Printf("Name: %s\n", project.Name)
	fmt.Printf("World: %s\n", project.World)
	fmt.Printf("Created: %s\n", project.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n\n", project.UpdatedAt.Format("2006-01-02 15:04:05"))

	fmt.Printf("=== WORLDBOOK (%d entries) ===\n", len(worldbooks))
	for i, wb := range worldbooks {
		fmt.Printf("[%d] ID: %s | Tags: %s\n", i+1, wb.ID, strings.Join(wb.Tags, ","))
		content := wb.Content
		if len(content) > 80 {
			content = content[:77] + "..."
		}
		fmt.Printf("    %s\n\n", strings.ReplaceAll(content, "\n", " "))
	}

	fmt.Printf("=== CHARACTERS (%d total) ===\n", len(characters))
	for _, ch := range characters {
		fmt.Printf("Name: %s \n", ch.Name)
		content := ch.Profile
		if len(content) > 100 {
			content = content[:97] + "..."
		}
		fmt.Printf("    %s\n\n", strings.ReplaceAll(content, "\n", " "))
	}

	fmt.Printf("=== CHAPTERS (%d total) ===\n", len(chapters))
	for _, ch := range chapters {
		fmt.Printf("[%d] %s\n", ch.Index, ch.Title)
		content := ch.Content
		if len(content) > 100 {
			content = content[:97] + "..."
		}
		fmt.Printf("    %s\n\n", strings.ReplaceAll(content, "\n", " "))
	}

	return nil
}
