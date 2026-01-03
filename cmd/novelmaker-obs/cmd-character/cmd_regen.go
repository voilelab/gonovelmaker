package cmdcharacter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
	"github.com/voilelab/gonovelmaker/internal/llmbackend"
	"github.com/voilelab/gonovelmaker/internal/nmutil"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
	"github.com/voilelab/gonovelmaker/novelmaker"
)

type regenCmd struct {
	json     bool
	filepath string
	backend  string
	model    string
	timeout  int

	llmBackendMaker llmbackend.LLMBackendMaker

	cmd *cobra.Command
}

func newRegenCmd(llmBackendMaker llmbackend.LLMBackendMaker) *regenCmd {
	g := &regenCmd{
		llmBackendMaker: llmBackendMaker,
	}
	g.cmd = &cobra.Command{
		Use:   "regen",
		Short: "Regenerate an existing character using its stored prompt",
		Long: `Regenerates an existing character profile based on the prompt stored in its frontmatter.
The filepath should be relative to the vault root (e.g., "Character/alice.md").`,
		RunE: g.run,
	}

	g.cmd.Flags().BoolVarP(&g.json, "json", "j", false, "Output in JSON format")

	g.cmd.Flags().StringVarP(&g.filepath, "filepath", "f", "", "Path to the character file relative to vault root (required)")
	g.cmd.MarkFlagRequired("filepath")

	// Allow overriding config values per-command
	g.cmd.Flags().StringVar(&g.backend, "backend", "", "LLM backend to use (optional, uses default if not specified)")
	g.cmd.Flags().StringVar(&g.model, "model", "", "Model to use, overrides config (optional)")
	g.cmd.Flags().IntVar(&g.timeout, "timeout", 0, "Timeout in seconds for the API request (optional)")
	return g
}

func (g *regenCmd) run(cmd *cobra.Command, args []string) error {
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

	// Load the target character from the specified filepath
	targetCharacter, err := vault.LoadCharacterByPath(g.filepath)
	if err != nil {
		return fmt.Errorf("failed to load target character from %s: %w", g.filepath, err)
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
	effectiveModel := nmutil.FirstNonEmptyString(g.model, backend.Model)
	effectiveTimeout := time.Duration(nmutil.FirstNonZero(g.timeout, backend.Timeout)) * time.Second

	if !g.json {
		fmt.Println("Regenerating character with OpenAI...")
		fmt.Printf("  Project: %s\n", project.Name)
		fmt.Printf("  Model: %s\n", effectiveModel)
		fmt.Printf("  Character: %s\n", targetCharacter.Name)
		fmt.Printf("  Prompt: %s\n", targetCharacter.Prompt)
		fmt.Printf("  Context: %d worldbook entries, %d existing characters\n", len(worldbooks), len(characters))
	}

	// Load character prompt template from vault
	characterPrompt, err := vault.LoadCharacterPrompt()
	if err != nil {
		return fmt.Errorf("failed to load character prompt: %w", err)
	}

	llmBackend := g.llmBackendMaker(
		backend.APIKey,
		backend.BaseURL,
		effectiveModel,
	)

	renderer := novelmaker.NewRenderer(
		llmBackend,
		effectiveTimeout,
	)

	// Call OpenAI API to regenerate character profile
	profile, usage, err := renderer.RenderCharacter(
		project, characterPrompt, worldbooks, characters, targetCharacter.Prompt, targetCharacter.Name)
	if err != nil {
		return fmt.Errorf("failed to generate character: %w", err)
	}

	// Update the character with new profile
	targetCharacter.Profile = profile

	// Update the character file
	if err := vault.UpdateCharacter(g.filepath, targetCharacter); err != nil {
		return fmt.Errorf("failed to update character file %s: %w", g.filepath, err)
	}

	if !g.json {
		fmt.Printf("\n✓ Successfully regenerated character!\n")
		fmt.Printf("  File: %s\n", g.filepath)
		fmt.Printf("  Name: %s\n", targetCharacter.Name)
		fmt.Printf("\nToken Usage:\n")
		fmt.Printf("  Input tokens:  %d\n", usage.InputTokens)
		fmt.Printf("  Output tokens: %d\n", usage.OutputTokens)
		fmt.Printf("  Total tokens:  %d\n", usage.TotalTokens)
	} else {
		output := map[string]any{
			"filepath":      g.filepath,
			"name":          targetCharacter.Name,
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
