package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
	"github.com/voilelab/gonovelmaker/novelmaker"
)

var rootCmd = &cobra.Command{
	Use:   "novelmaker-obs",
	Short: "A CLI tool for generating novels in Obsidian vaults",
	Long: `novelmaker-obs is a CLI tool that works with Obsidian vaults to manage 
novel projects, worldbooks, and chapters. It uses OpenAI API to generate new chapters 
based on your existing content.

Configuration:
  Config file: ~/.novelmaker/config.toml
  Example:
    openai_api_key = "sk-xxx"
    model = "gpt-4o"

  You can also set OPENAI_API_KEY as an environment variable.`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a novel project structure in the current directory",
	Long: `Creates Config/, World/, and Story/ subdirectories in the current directory, 
along with sample files.`,
	RunE: runInit,
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan and display the current novel project structure",
	Long: `Scans the project directory structure (Config/, World/, Story/) and displays 
the project configuration, worldbook entries, and chapters.`,
	RunE: runScan,
}

var genNextCmd = &cobra.Command{
	Use:   "gen-next",
	Short: "Generate the next chapter using OpenAI API",
	Long: `Generates a new chapter based on existing project configuration, worldbook, 
and previous chapters using OpenAI API.`,
	RunE: runGenNext,
}

var genCharCmd = &cobra.Command{
	Use:   "gen-char",
	Short: "Generate a new character using OpenAI API",
	Long: `Generates a new character profile based on the project context and a user-provided prompt 
using OpenAI API.`,
	RunE: runGenChar,
}

var (
	jsonOutput    bool
	title         string
	charPrompt    string
	charName      string
	chapterPrompt string
	apiKey        string
	baseURL       string
	model         string
)

func init() {
	scanCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")

	genNextCmd.Flags().StringVarP(&title, "title", "t", "", "Title for the next chapter (required)")
	genNextCmd.Flags().StringVarP(&chapterPrompt, "prompt", "p", "", "Additional prompt/instruction for chapter generation (optional)")
	genNextCmd.MarkFlagRequired("title")

	// Allow overriding config values per-command
	genNextCmd.Flags().StringVar(&apiKey, "api-key", "", "OpenAI API key to override config (optional)")
	genNextCmd.Flags().StringVar(&baseURL, "base-url", "", "OpenAI base URL to override config (optional)")
	genNextCmd.Flags().StringVar(&model, "model", "", "Model to use, overrides config (optional)")

	genCharCmd.Flags().StringVarP(&charPrompt, "prompt", "p", "", "Description/prompt for the character to generate")
	genCharCmd.Flags().StringVarP(&charName, "name", "n", "", "Name for the character (optional, will be extracted from AI response if not provided)")
	genCharCmd.MarkFlagRequired("name")

	// Allow overriding config values per-command
	genCharCmd.Flags().StringVar(&apiKey, "api-key", "", "OpenAI API key to override config (optional)")
	genCharCmd.Flags().StringVar(&baseURL, "base-url", "", "OpenAI base URL to override config (optional)")
	genCharCmd.Flags().StringVar(&model, "model", "", "Model to use, overrides config (optional)")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(genNextCmd)
	rootCmd.AddCommand(genCharCmd)
}

func main() {
	// Initialize config (creates empty config file if it doesn't exist)
	config.Load()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runInit(cmd *cobra.Command, args []string) error {
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
	fmt.Printf("  Created in: %s\n", cwd)
	fmt.Println("  - Config/project.md")
	fmt.Println("  - World/001_world_sample.md")
	fmt.Println("  - Character/001_character_sample.md")
	fmt.Println("  - Story/001_prologue.md")
	fmt.Println("  - README.md")
	fmt.Println("  - .obsidian/app.json")

	return nil
}

func runScan(cmd *cobra.Command, args []string) error {
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

	// Load chapters
	chapters, err := vault.LoadChapters()
	if err != nil {
		return fmt.Errorf("failed to load chapters: %w", err)
	}

	if jsonOutput {
		data := map[string]any{
			"project":    project,
			"worldbooks": worldbooks,
			"chapters":   chapters,
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")

		return encoder.Encode(data)
	}

	// Pretty print
	fmt.Println("=== PROJECT ===")
	fmt.Printf("ID: %s\n", project.ID)
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

	fmt.Printf("=== CHAPTERS (%d total) ===\n", len(chapters))
	for _, ch := range chapters {
		fmt.Printf("[%d] %s (ID: %s)\n", ch.Index, ch.Title, ch.ID)
		content := ch.Content
		if len(content) > 100 {
			content = content[:97] + "..."
		}
		fmt.Printf("    %s\n\n", strings.ReplaceAll(content, "\n", " "))
	}

	return nil
}

func runGenNext(cmd *cobra.Command, args []string) error {
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
	if apiKey != "" {
		effectiveAPIKey = apiKey
	}
	effectiveModel := cfg.GetModelOrDefault()
	if model != "" {
		effectiveModel = model
	}
	effectiveBaseURL := cfg.BaseURL
	if baseURL != "" {
		effectiveBaseURL = baseURL
	}

	fmt.Println("Generating next chapter with OpenAI...")
	fmt.Printf("  Project: %s\n", project.Name)
	fmt.Printf("  Model: %s\n", effectiveModel)
	fmt.Printf("  Target: %s\n", title)
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
	if title == "" {
		return fmt.Errorf("target title cannot be empty")
	}

	// Call OpenAI API
	content, err := renderer.RenderPrompt(
		project, worldbooks, characters, chapters, title, chapterPrompt)
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
		Title:   title,
		Content: content,
	}

	if err := vault.AddChapter(&ch); err != nil {
		return fmt.Errorf("failed to add chapter to vault: %w", err)
	}

	autoPath := filepath.Join(cwd, "Story", fmt.Sprintf("%03d_%s.md", nextIndex, chapterID))
	fmt.Printf("\n✓ Successfully generated chapter!\n")
	fmt.Printf("  File: %s\n", autoPath)
	fmt.Printf("  Index: %d\n", nextIndex)
	fmt.Printf("  Title: %s\n", title)

	return nil
}

func runGenChar(cmd *cobra.Command, args []string) error {
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
	if apiKey != "" {
		effectiveAPIKey = apiKey
	}
	effectiveModel := cfg.GetModelOrDefault()
	if model != "" {
		effectiveModel = model
	}
	effectiveBaseURL := cfg.BaseURL
	if baseURL != "" {
		effectiveBaseURL = baseURL
	}

	fmt.Println("Generating character with OpenAI...")
	fmt.Printf("  Project: %s\n", project.Name)
	fmt.Printf("  Model: %s\n", effectiveModel)
	fmt.Printf("  Prompt: %s\n", charPrompt)
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
		project, worldbooks, characters, charPrompt, charName)
	if err != nil {
		return fmt.Errorf("failed to generate character: %w", err)
	}

	// Use provided name or extracted name
	finalName := charName
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

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, ":", "")
	s = strings.ReplaceAll(s, "?", "")
	s = strings.ReplaceAll(s, "!", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, ".", "")
	// Keep alphanumeric characters, underscore, and Unicode letters (including Chinese, Japanese, etc.)
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r > 127 {
			result.WriteRune(r)
		}
	}
	return result.String()
}
