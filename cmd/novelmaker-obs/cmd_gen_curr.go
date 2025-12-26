package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
	"github.com/voilelab/gonovelmaker/internal/nmutil"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
	"github.com/voilelab/gonovelmaker/novelmaker"
)

type GenCurrCmd struct {
	json         bool
	filepath     string
	prevChapters int
	backend      string
	apiKey       string
	baseURL      string
	model        string
	timeout      int

	llmBackendMaker LLMBackendMaker

	cmd *cobra.Command
}

func NewGenCurrCmd(llmBackendMaker LLMBackendMaker) *GenCurrCmd {
	g := &GenCurrCmd{
		llmBackendMaker: llmBackendMaker,
	}
	g.cmd = &cobra.Command{
		Use:   "gen-curr",
		Short: "Regenerate an existing chapter using its chapter.prompt",
		Long: `Regenerates an existing chapter based on the prompt stored in its frontmatter.
The filepath should be relative to the vault root (e.g., "Story/001_ch1.md").`,
		RunE: g.run,
	}

	g.cmd.Flags().BoolVarP(&g.json, "json", "j", false, "Output in JSON format")

	g.cmd.Flags().StringVarP(&g.filepath, "filepath", "f", "", "Path to the chapter file relative to vault root (required)")
	g.cmd.Flags().IntVarP(&g.prevChapters, "prev-chapters", "c", 3, "Number of previous chapters to include as context (optional, default 3, max 10)")
	g.cmd.MarkFlagRequired("filepath")

	// Allow overriding config values per-command
	g.cmd.Flags().StringVar(&g.backend, "backend", "", "LLM backend to use (optional, uses default if not specified)")
	g.cmd.Flags().StringVar(&g.apiKey, "api-key", "", "OpenAI API key to override config (optional)")
	g.cmd.Flags().StringVar(&g.baseURL, "base-url", "", "OpenAI base URL to override config (optional)")
	g.cmd.Flags().StringVar(&g.model, "model", "", "Model to use, overrides config (optional)")
	g.cmd.Flags().IntVar(&g.timeout, "timeout", 0, "Timeout in seconds for the API request (optional)")
	return g
}

func (g *GenCurrCmd) run(cmd *cobra.Command, args []string) error {
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

	// Load all chapters
	chapters, err := vault.LoadChapters()
	if err != nil {
		return fmt.Errorf("failed to load chapters: %w", err)
	}

	// Load the target chapter from the specified filepath
	targetChapter, err := vault.LoadChapterByPath(g.filepath)
	if err != nil {
		return fmt.Errorf("failed to load target chapter from %s: %w", g.filepath, err)
	}

	// Get backend configuration
	backend := cfg.GetBackend(g.backend)
	if backend == nil {
		if g.backend != "" {
			return fmt.Errorf("backend '%s' not found in config", g.backend)
		}
		return fmt.Errorf("no LLM backend configured. Please configure at least one backend in ~/.novelmaker/config.toml")
	}

	// Determine effective API settings (flags override config)
	effectiveAPIKey := nmutil.FirstNonEmptyString(g.apiKey, backend.APIKey)
	effectiveModel := nmutil.FirstNonEmptyString(g.model, backend.Model)
	effectiveBaseURL := nmutil.FirstNonEmptyString(g.baseURL, backend.BaseURL)
	effectiveTimeout := time.Duration(nmutil.FirstNonZero(g.timeout, backend.Timeout)) * time.Second

	if !g.json {
		fmt.Println("Regenerating chapter with OpenAI...")
		fmt.Printf("  Project: %s\n", project.Name)
		fmt.Printf("  Model: %s\n", effectiveModel)
		fmt.Printf("  Target: %s (Index: %d)\n", targetChapter.Title, targetChapter.Index)
		fmt.Printf("  Prompt: %s\n", targetChapter.Prompt)
		fmt.Printf("  Context: %d worldbook entries, %d characters, %d chapters\n", len(worldbooks), len(characters), len(chapters))
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

	// Get previous chapters (those with index < target chapter index)
	var prevChapters []novelmaker.Chapter
	for _, ch := range chapters {
		if ch.Index < targetChapter.Index {
			prevChapters = append(prevChapters, ch)
		}
	}

	// Limit to the last N previous chapters
	prevK := min(g.prevChapters, len(prevChapters), maxPrevChapters)
	if len(prevChapters) > prevK {
		prevChapters = prevChapters[len(prevChapters)-prevK:]
	}

	// Call OpenAI API to regenerate content
	content, usage, err := renderer.RenderChapter(
		project, chapterPrompt, worldbooks, characters, prevChapters, targetChapter.Title, targetChapter.Prompt)
	if err != nil {
		return fmt.Errorf("failed to generate chapter: %w", err)
	}

	// Update the chapter with new content
	targetChapter.Content = content

	// Update the chapter file
	if err := vault.UpdateChapter(g.filepath, targetChapter); err != nil {
		return fmt.Errorf("failed to update chapter file %s: %w", g.filepath, err)
	}

	if !g.json {
		fmt.Printf("\n✓ Successfully regenerated chapter!\n")
		fmt.Printf("  File: %s\n", g.filepath)
		fmt.Printf("  Index: %d\n", targetChapter.Index)
		fmt.Printf("  Title: %s\n", targetChapter.Title)
		fmt.Printf("\nToken Usage:\n")
		fmt.Printf("  Input tokens:  %d\n", usage.InputTokens)
		fmt.Printf("  Output tokens: %d\n", usage.OutputTokens)
		fmt.Printf("  Total tokens:  %d\n", usage.TotalTokens)
	} else {
		output := map[string]any{
			"filepath":      g.filepath,
			"index":         targetChapter.Index,
			"title":         targetChapter.Title,
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
