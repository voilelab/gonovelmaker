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

type GenNextCmd struct {
	title         string
	chapterPrompt string
	apiKey        string
	baseURL       string
	model         string

	cmd *cobra.Command
}

func NewGenNextCmd() *GenNextCmd {
	g := &GenNextCmd{}
	g.cmd = &cobra.Command{
		Use:   "gen-next",
		Short: "Generate the next chapter using OpenAI API",
		Long: `Generates a new chapter based on existing project configuration, worldbook, 
and previous chapters using OpenAI API.`,
		RunE: g.run,
	}

	g.cmd.Flags().StringVarP(&g.title, "title", "t", "", "Title for the next chapter (required)")
	g.cmd.Flags().StringVarP(&g.chapterPrompt, "prompt", "p", "", "Additional prompt/instruction for chapter generation (optional)")
	g.cmd.MarkFlagRequired("title")

	// Allow overriding config values per-command
	g.cmd.Flags().StringVar(&g.apiKey, "api-key", "", "OpenAI API key to override config (optional)")
	g.cmd.Flags().StringVar(&g.baseURL, "base-url", "", "OpenAI base URL to override config (optional)")
	g.cmd.Flags().StringVar(&g.model, "model", "", "Model to use, overrides config (optional)")
	return g
}

func (g *GenNextCmd) run(cmd *cobra.Command, args []string) error {
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

	fmt.Println("Generating next chapter with OpenAI...")
	fmt.Printf("  Project: %s\n", project.Name)
	fmt.Printf("  Model: %s\n", effectiveModel)
	fmt.Printf("  Target: %s\n", g.title)
	fmt.Printf("  Context: %d worldbook entries, %d characters, %d previous chapters\n", len(worldbooks), len(characters), len(chapters))

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

	// Validate target title
	if g.title == "" {
		return fmt.Errorf("target title cannot be empty")
	}

	// Call OpenAI API
	content, err := renderer.RenderPrompt(
		project, worldbooks, characters, chapters, g.title, g.chapterPrompt)
	if err != nil {
		return fmt.Errorf("failed to generate chapter: %w", err)
	}

	// Determine next index
	nextIndex := 1
	if len(chapters) > 0 {
		maxIndex := 0
		for _, ch := range chapters {
			if ch.Index > maxIndex {
				maxIndex = ch.Index
			}
		}
		nextIndex = maxIndex + 1
	}

	// Generate ID based on number of chapters
	chapterID := fmt.Sprintf("ch%d", len(chapters)+1)

	// Create chapter struct and add via Vault helper (no --out-file behavior)
	ch := novelmaker.Chapter{
		ID:      chapterID,
		Index:   nextIndex,
		Title:   g.title,
		Content: content,
	}

	if err := vault.AddChapter(&ch); err != nil {
		return fmt.Errorf("failed to add chapter to vault: %w", err)
	}

	autoPath := filepath.Join(cwd, "Story", fmt.Sprintf("%03d_%s.md", nextIndex, chapterID))
	fmt.Printf("\n✓ Successfully generated chapter!\n")
	fmt.Printf("  File: %s\n", autoPath)
	fmt.Printf("  Index: %d\n", nextIndex)
	fmt.Printf("  Title: %s\n", g.title)

	return nil
}
