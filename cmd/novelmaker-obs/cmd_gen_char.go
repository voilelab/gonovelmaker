package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
	"github.com/voilelab/gonovelmaker/novelmaker"
)

type GenCharCmd struct {
	prompt  string
	name    string
	apiKey  string
	baseURL string
	model   string

	cmd *cobra.Command
}

func NewGenCharCmd() *GenCharCmd {
	g := &GenCharCmd{}
	g.cmd = &cobra.Command{
		Use:   "gen-char",
		Short: "Generate a new character using OpenAI API",
		Long: `Generates a new character profile based on the project context and a user-provided prompt 
using OpenAI API.`,
		RunE: g.run,
	}

	g.cmd.Flags().StringVarP(&g.prompt, "prompt", "p", "", "Description/prompt for the character to generate")
	g.cmd.Flags().StringVarP(&g.name, "name", "n", "", "Name for the character (optional, will be extracted from AI response if not provided)")
	g.cmd.MarkFlagRequired("prompt")

	// Allow overriding config values per-command
	g.cmd.Flags().StringVar(&g.apiKey, "api-key", "", "OpenAI API key to override config (optional)")
	g.cmd.Flags().StringVar(&g.baseURL, "base-url", "", "OpenAI base URL to override config (optional)")
	g.cmd.Flags().StringVar(&g.model, "model", "", "Model to use, overrides config (optional)")

	return g
}

func (g *GenCharCmd) run(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

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

	// Load existing characters for context
	characters, err := vault.LoadCharacters()
	if err != nil {
		return fmt.Errorf("failed to load characters: %w", err)
	}

	// Determine effective API settings (flags override config)
	effectiveAPIKey := cfg.OpenAIKey
	if g.apiKey != "" {
		effectiveAPIKey = g.apiKey
	}
	effectiveModel := cfg.GetModelOrDefault()
	if g.model != "" {
		effectiveModel = g.model
	}
	effectiveBaseURL := cfg.BaseURL
	if g.baseURL != "" {
		effectiveBaseURL = g.baseURL
	}

	fmt.Println("Generating character with OpenAI...")
	fmt.Printf("  Project: %s\n", project.Name)
	fmt.Printf("  Model: %s\n", effectiveModel)
	fmt.Printf("  Prompt: %s\n", g.prompt)
	fmt.Printf("  Context: %d worldbook entries, %d existing characters\n", len(worldbooks), len(characters))

	// Load prompt templates
	promptTemplates, err := config.LoadPromptTemplates()
	if err != nil {
		return fmt.Errorf("failed to load prompt templates: %w", err)
	}

	renderer := novelmaker.NewRenderer(
		effectiveAPIKey,
		effectiveModel,
		effectiveBaseURL,
		promptTemplates.ChapterTemplate,
		promptTemplates.CharacterTemplate)

	// Call OpenAI API
	profile, extractedName, err := renderer.RenderCharacter(
		project, worldbooks, characters, g.prompt, g.name)
	if err != nil {
		return fmt.Errorf("failed to generate character: %w", err)
	}

	// Use provided name or extracted name
	finalName := g.name
	if finalName == "" {
		finalName = extractedName
	}

	if finalName == "" {
		return fmt.Errorf("could not determine character name. Please provide --name flag")
	}

	// Generate character ID
	charID := slugify(finalName)

	// Create character struct
	ch := novelmaker.Character{
		ID:      charID,
		Name:    finalName,
		Main:    false,
		Profile: profile,
	}

	// Always use Vault helper to add character to vault (do not support --out-file here)
	if err := vault.AddCharacter(&ch); err != nil {
		return fmt.Errorf("failed to add character to vault: %w", err)
	}

	autoPath := filepath.Join(cwd, "Character", fmt.Sprintf("%s.md", charID))
	fmt.Printf("\n✓ Successfully generated character!\n")
	fmt.Printf("  File: %s\n", autoPath)
	fmt.Printf("  Name: %s\n", finalName)
	fmt.Printf("  ID: %s\n", charID)
	return nil
}
