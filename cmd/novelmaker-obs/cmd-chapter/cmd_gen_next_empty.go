package cmdchapter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
	"github.com/voilelab/gonovelmaker/novelmaker"
)

type GenNextEmptyCmd struct {
	json   bool
	title  string
	prompt string

	cmd *cobra.Command
}

func NewGenNextEmptyCmd() *GenNextEmptyCmd {
	g := &GenNextEmptyCmd{}
	g.cmd = &cobra.Command{
		Use:   "gen-next-empty",
		Short: "Generate an empty next chapter without AI content",
		Long: `Creates a new empty chapter file with frontmatter but no content. 
This is useful for manually writing chapters or creating placeholders.`,
		RunE: g.run,
	}

	g.cmd.Flags().BoolVarP(&g.json, "json", "j", false, "Output in JSON format")
	g.cmd.Flags().StringVarP(&g.title, "title", "t", "", "Title for the next chapter (required)")
	g.cmd.Flags().StringVarP(&g.prompt, "prompt", "p", "", "Additional prompt/instruction for chapter (optional)")
	g.cmd.MarkFlagRequired("title")

	return g
}

func (g *GenNextEmptyCmd) run(cmd *cobra.Command, args []string) error {
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

	// Load chapters to determine next index
	chapters, err := vault.LoadChapters()
	if err != nil {
		return fmt.Errorf("failed to load chapters: %w", err)
	}

	// Validate title
	if g.title == "" {
		return fmt.Errorf("title cannot be empty")
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

	if !g.json {
		fmt.Println("Creating empty chapter...")
		fmt.Printf("  Index: %d\n", nextIndex)
		fmt.Printf("  Title: %s\n", g.title)
	}

	// Create empty chapter
	ch := novelmaker.Chapter{
		Index:   nextIndex,
		Title:   g.title,
		Prompt:  g.prompt,
		Content: "", // Empty content
	}

	filePath, err := vault.AddChapter(&ch)
	if err != nil {
		return fmt.Errorf("failed to add chapter to vault: %w", err)
	}

	if !g.json {
		fmt.Printf("\n✓ Successfully created empty chapter!\n")
		fmt.Printf("  File: %s\n", filePath)
		fmt.Printf("  Index: %d\n", nextIndex)
		fmt.Printf("  Title: %s\n", g.title)
	} else {
		output := map[string]any{
			"filepath": filePath,
			"index":    nextIndex,
			"title":    g.title,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			return fmt.Errorf("failed to encode JSON output: %w", err)
		}
	}

	return nil
}
