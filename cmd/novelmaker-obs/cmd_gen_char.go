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

type GenCharCmd struct {
	json    bool
	prompt  string
	name    string
	apiKey  string
	baseURL string
	model   string
	timeout int

	llmBackendMaker LLMBackendMaker

	cmd *cobra.Command
}

func NewGenCharCmd(llmBackendMaker LLMBackendMaker) *GenCharCmd {
	g := &GenCharCmd{
		llmBackendMaker: llmBackendMaker,
	}
	g.cmd = &cobra.Command{
		Use:   "gen-char",
		Short: "Generate a new character using OpenAI API",
		Long: `Generates a new character profile based on the project context and a user-provided prompt 
using OpenAI API.`,
		RunE: g.run,
	}

	g.cmd.Flags().BoolVarP(&g.json, "json", "j", false, "Output in JSON format")

	g.cmd.Flags().StringVarP(&g.name, "name", "n", "", "Name for the character (required)")
	g.cmd.Flags().StringVarP(&g.prompt, "prompt", "p", "", "Description/prompt for the character to generate")
	g.cmd.MarkFlagRequired("name")

	// Allow overriding config values per-command
	g.cmd.Flags().StringVar(&g.apiKey, "api-key", "", "OpenAI API key to override config (optional)")
	g.cmd.Flags().StringVar(&g.baseURL, "base-url", "", "OpenAI base URL to override config (optional)")
	g.cmd.Flags().StringVar(&g.model, "model", "", "Model to use, overrides config (optional)")
	g.cmd.Flags().IntVar(&g.timeout, "timeout", 0, "Timeout in seconds for the API request (optional)")
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

	// Get backend configuration
	backend := cfg.GetBackend("")
	if backend == nil {
		return fmt.Errorf("no LLM backend configured. Please configure at least one backend in ~/.novelmaker/config.toml")
	}

	// Determine effective API settings (flags override config)
	effectiveAPIKey := backend.APIKey
	if g.apiKey != "" {
		effectiveAPIKey = g.apiKey
	}
	effectiveModel := backend.Model
	if g.model != "" {
		effectiveModel = g.model
	}
	effectiveBaseURL := backend.BaseURL
	if g.baseURL != "" {
		effectiveBaseURL = g.baseURL
	}
	effectiveTimeout := time.Duration(backend.Timeout) * time.Second
	if g.timeout > 0 {
		effectiveTimeout = time.Duration(g.timeout) * time.Second
	}

	if !g.json {
		fmt.Println("Generating character with OpenAI...")
		fmt.Printf("  Project: %s\n", project.Name)
		fmt.Printf("  Model: %s\n", effectiveModel)
		fmt.Printf("  Prompt: %s\n", g.prompt)
		fmt.Printf("  Context: %d worldbook entries, %d existing characters\n", len(worldbooks), len(characters))
	}

	// Load character prompt template from vault
	characterPrompt, err := vault.LoadCharacterPrompt()
	if err != nil {
		return fmt.Errorf("failed to load character prompt: %w", err)
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

	// Call OpenAI API
	profile, usage, err := renderer.RenderCharacter(
		project, characterPrompt, worldbooks, characters, g.prompt, g.name)
	if err != nil {
		return fmt.Errorf("failed to generate character: %w", err)
	}

	// Create character struct
	ch := novelmaker.Character{
		Name:    g.name,
		Main:    false,
		Prompt:  g.prompt,
		Profile: profile,
	}

	filePath, err := vault.AddCharacter(&ch)
	if err != nil {
		return fmt.Errorf("failed to add character to vault: %w", err)
	}

	if !g.json {
		fmt.Printf("\n✓ Successfully generated character!\n")
		fmt.Printf("  File: %s\n", filePath)
		fmt.Printf("  Name: %s\n", g.name)
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
