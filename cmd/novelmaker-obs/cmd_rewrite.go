package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
	"github.com/voilelab/gonovelmaker/internal/llmbackend"
	"github.com/voilelab/gonovelmaker/internal/nmutil"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
	"github.com/voilelab/gonovelmaker/novelmaker"
)

type RewriteCmd struct {
	json        bool
	lineStart   int
	lineEnd     int
	contextPrev int
	contextNext int
	backend     string
	model       string
	timeout     int

	llmBackendMaker llmbackend.LLMBackendMaker

	cmd *cobra.Command
}

func NewRewriteCmd(llmBackendMaker llmbackend.LLMBackendMaker) *RewriteCmd {
	r := &RewriteCmd{
		llmBackendMaker: llmBackendMaker,
	}

	r.cmd = &cobra.Command{
		Use:   "rewrite [filepath]",
		Short: "Rewrite a specific range of lines in a chapter file",
		Long: `Rewrite a specific range of lines in a chapter file using AI.

This command reads the specified file, extracts the target lines and context,
sends them to the AI model for rewriting, and updates the file with the new content.

Example:
  novelmaker-obs rewrite Story/001_ch1.md --line-start 10 --line-end 15 --context-prev 5 --context-next 5
  novelmaker-obs rewrite Story/001_ch1.md -s 10 -e 15 -p 5 -n 5 --model gpt-4`,
		Args: cobra.ExactArgs(1),
		RunE: r.run,
	}

	r.cmd.Flags().BoolVarP(&r.json, "json", "j", false, "Output in JSON format")

	r.cmd.Flags().IntVarP(&r.lineStart, "line-start", "s", 0, "Starting line number to rewrite (required, 1-indexed)")
	r.cmd.Flags().IntVarP(&r.lineEnd, "line-end", "e", 0, "Ending line number to rewrite (required, 1-indexed, inclusive)")
	r.cmd.Flags().IntVarP(&r.contextPrev, "context-prev", "p", 3, "Number of lines before the target as context (default 3)")
	r.cmd.Flags().IntVarP(&r.contextNext, "context-next", "n", 3, "Number of lines after the target as context (default 3)")

	r.cmd.MarkFlagRequired("line-start")
	r.cmd.MarkFlagRequired("line-end")

	// Allow overriding config values per-command
	r.cmd.Flags().StringVar(&r.backend, "backend", "", "LLM backend to use (optional, uses default if not specified)")
	r.cmd.Flags().StringVar(&r.model, "model", "", "Model to use, overrides config (optional)")
	r.cmd.Flags().IntVar(&r.timeout, "timeout", 0, "Timeout in seconds for the API request (optional)")

	return r
}

func (r *RewriteCmd) run(cmd *cobra.Command, args []string) error {
	filepath := args[0]

	// Validate line numbers
	if r.lineStart < 1 {
		return fmt.Errorf("line-start must be >= 1, got %d", r.lineStart)
	}
	if r.lineEnd < r.lineStart {
		return fmt.Errorf("line-end (%d) must be >= line-start (%d)", r.lineEnd, r.lineStart)
	}
	if r.contextPrev < 0 {
		return fmt.Errorf("context-prev must be >= 0, got %d", r.contextPrev)
	}
	if r.contextNext < 0 {
		return fmt.Errorf("context-next must be >= 0, got %d", r.contextNext)
	}

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

	// Get backend configuration
	backend := cfg.GetBackend(r.backend)
	if backend == nil {
		if r.backend != "" {
			return fmt.Errorf("backend '%s' not found in config", r.backend)
		}
		return fmt.Errorf("no LLM backend configured. Please configure at least one backend in ~/.novelmaker/config.toml")
	}

	// Determine effective API settings (flags override config)
	effectiveModel := nmutil.FirstNonEmptyString(r.model, backend.Model)
	effectiveTimeout := time.Duration(nmutil.FirstNonZero(r.timeout, backend.Timeout)) * time.Second

	if !r.json {
		fmt.Println("Rewriting lines with AI...")
		fmt.Printf("  Project: %s\n", project.Name)
		fmt.Printf("  Model: %s\n", effectiveModel)
		fmt.Printf("  File: %s\n", filepath)
		fmt.Printf("  Lines: %d-%d\n", r.lineStart, r.lineEnd)
		fmt.Printf("  Context: %d lines before, %d lines after\n", r.contextPrev, r.contextNext)
	}

	// Read the file and extract lines
	fileLines, err := r.readFileLines(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	totalLines := len(fileLines)
	if r.lineStart > totalLines {
		return fmt.Errorf("line-start (%d) exceeds total lines in file (%d)", r.lineStart, totalLines)
	}
	if r.lineEnd > totalLines {
		return fmt.Errorf("line-end (%d) exceeds total lines in file (%d)", r.lineEnd, totalLines)
	}

	// Extract context and target (convert to 0-indexed)
	contextStartIdx := max(0, r.lineStart-1-r.contextPrev)
	targetStartIdx := r.lineStart - 1
	targetEndIdx := r.lineEnd // inclusive, so no -1
	contextEndIdx := min(totalLines, r.lineEnd+r.contextNext)

	contextBefore := strings.Join(fileLines[contextStartIdx:targetStartIdx], "\n")
	targetSentence := strings.Join(fileLines[targetStartIdx:targetEndIdx], "\n")
	contextAfter := strings.Join(fileLines[targetEndIdx:contextEndIdx], "\n")

	// Load rewrite prompt template from vault
	rewritePrompt, err := vault.LoadRewritePrompt()
	if err != nil {
		return fmt.Errorf("failed to load rewrite prompt from vault: %w", err)
	}

	llmBackend := r.llmBackendMaker(
		backend.APIKey,
		backend.BaseURL,
		effectiveModel,
	)

	renderer := novelmaker.NewRenderer(
		llmBackend,
		effectiveTimeout,
	)

	slog.Info("Rewriting text",
		"project", project.Name,
		"model", effectiveModel,
		"file", filepath,
		"line_range", fmt.Sprintf("%d-%d", r.lineStart, r.lineEnd),
	)

	// Call AI API
	rewrittenContent, usage, err := renderer.RenderRewrite(
		project, rewritePrompt, contextBefore, targetSentence, contextAfter)
	if err != nil {
		return fmt.Errorf("failed to rewrite text: %w", err)
	}

	// Replace the target lines with rewritten content
	rewrittenLines := strings.Split(strings.TrimSpace(rewrittenContent), "\n")
	newFileLines := make([]string, 0, len(fileLines))
	newFileLines = append(newFileLines, fileLines[:targetStartIdx]...)
	newFileLines = append(newFileLines, rewrittenLines...)
	newFileLines = append(newFileLines, fileLines[targetEndIdx:]...)

	// Write back to file
	newContent := strings.Join(newFileLines, "\n")
	if err := os.WriteFile(filepath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	if !r.json {
		fmt.Printf("\n✓ Successfully rewrote lines %d-%d!\n", r.lineStart, r.lineEnd)
		fmt.Printf("  Original lines: %d\n", r.lineEnd-r.lineStart+1)
		fmt.Printf("  New lines: %d\n", len(rewrittenLines))
		fmt.Printf("\nToken Usage:\n")
		fmt.Printf("  Input tokens:  %d\n", usage.InputTokens)
		fmt.Printf("  Output tokens: %d\n", usage.OutputTokens)
		fmt.Printf("  Total tokens:  %d\n", usage.TotalTokens)
		fmt.Printf("\nRewritten content:\n")
		fmt.Printf("  %s\n", strings.ReplaceAll(rewrittenContent, "\n", "\n  "))
	} else {
		output := map[string]any{
			"filepath":            filepath,
			"original_line_count": r.lineEnd - r.lineStart + 1,
			"new_line_count":      len(rewrittenLines),
			"input_tokens":        usage.InputTokens,
			"output_tokens":       usage.OutputTokens,
			"total_tokens":        usage.TotalTokens,
			"rewritten_content":   rewrittenContent,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			return fmt.Errorf("failed to encode JSON output: %w", err)
		}
	}

	return nil
}

func (r *RewriteCmd) readFileLines(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
