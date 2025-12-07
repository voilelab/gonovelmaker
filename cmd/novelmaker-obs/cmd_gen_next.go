package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
	"github.com/voilelab/gonovelmaker/novelmaker"
)

const maxPrevChapters = 10

type GenNextCmd struct {
	json         bool
	title        string
	prompt       string
	prevChapters int
	apiKey       string
	baseURL      string
	model        string
	timeout      int

	llmBackendMaker LLMBackendMaker

	cmd *cobra.Command
}

func NewGenNextCmd(llmBackendMaker LLMBackendMaker) *GenNextCmd {
	g := &GenNextCmd{
		llmBackendMaker: llmBackendMaker,
	}
	g.cmd = &cobra.Command{
		Use:   "gen-next",
		Short: "Generate the next chapter using OpenAI API",
		Long: `Generates a new chapter based on existing project configuration, worldbook, 
and previous chapters using OpenAI API.`,
		RunE: g.run,
	}

	g.cmd.Flags().BoolVarP(&g.json, "json", "j", false, "Output in JSON format")

	g.cmd.Flags().StringVarP(&g.title, "title", "t", "", "Title for the next chapter (required)")
	g.cmd.Flags().StringVarP(&g.prompt, "prompt", "p", "", "Additional prompt/instruction for chapter generation (optional)")
	g.cmd.Flags().IntVarP(&g.prevChapters, "prev-chapters", "c", 3, "Number of previous chapters to include as context (optional, default 3, max 10)")
	g.cmd.MarkFlagRequired("title")

	// Allow overriding config values per-command
	g.cmd.Flags().StringVar(&g.apiKey, "api-key", "", "OpenAI API key to override config (optional)")
	g.cmd.Flags().StringVar(&g.baseURL, "base-url", "", "OpenAI base URL to override config (optional)")
	g.cmd.Flags().StringVar(&g.model, "model", "", "Model to use, overrides config (optional)")
	g.cmd.Flags().IntVar(&g.timeout, "timeout", 0, "Timeout in seconds for the API request (optional)")
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
	effectiveTimeout := time.Duration(cfg.Timeout) * time.Second
	if g.timeout > 0 {
		effectiveTimeout = time.Duration(g.timeout) * time.Second
	}

	if !g.json {
		fmt.Println("Generating next chapter with OpenAI...")
		fmt.Printf("  Project: %s\n", project.Name)
		fmt.Printf("  Model: %s\n", effectiveModel)
		fmt.Printf("  Target: %s\n", g.title)
		fmt.Printf("  Context: %d worldbook entries, %d characters, %d previous chapters\n", len(worldbooks), len(characters), len(chapters))
	}

	// Load chapter prompt template from vault
	chapterPrompt, err := vault.LoadChapterPrompt()
	if err != nil {
		return fmt.Errorf("failed to load chapter prompt from vault: %w", err)
	}

	llmBackend := g.llmBackendMaker(
		effectiveAPIKey,
		effectiveBaseURL,
		effectiveModel,
	)

	renderer := novelmaker.NewRenderer(
		llmBackend,
		effectiveTimeout,
	)

	// Validate target title
	if g.title == "" {
		return fmt.Errorf("target title cannot be empty")
	}

	prevK := min(g.prevChapters, len(chapters), maxPrevChapters)
	prevChapters := chapters[len(chapters)-prevK:]

	// Call OpenAI API
	content, usage, err := renderer.RenderChapter(
		project, chapterPrompt, worldbooks, characters, prevChapters, g.title, g.prompt)
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

	ch := novelmaker.Chapter{
		Index:   nextIndex,
		Title:   g.title,
		Prompt:  g.prompt,
		Content: content,
	}

	filePath, err := vault.AddChapter(&ch)
	if err != nil {
		return fmt.Errorf("failed to add chapter to vault: %w", err)
	}

	if !g.json {
		fmt.Printf("\n✓ Successfully generated chapter!\n")
		fmt.Printf("  File: %s\n", filePath)
		fmt.Printf("  Index: %d\n", nextIndex)
		fmt.Printf("  Title: %s\n", g.title)
		fmt.Printf("\nToken Usage:\n")
		fmt.Printf("  Input tokens:  %d\n", usage.InputTokens)
		fmt.Printf("  Output tokens: %d\n", usage.OutputTokens)
		fmt.Printf("  Total tokens:  %d\n", usage.TotalTokens)
	} else {
		output := map[string]any{
			"filepath":      filePath,
			"input_tokens":  usage.InputTokens,
			"output_tokens": usage.OutputTokens,
			"total_tokens":  usage.TotalTokens,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			return fmt.Errorf("failed to encode JSON output: %w", err)
		}
	}

	return nil
}
