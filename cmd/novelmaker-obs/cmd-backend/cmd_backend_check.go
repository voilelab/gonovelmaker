package cmdbackend

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
	"github.com/voilelab/gonovelmaker/internal/llmbackend"
)

type BackendCheckCmd struct {
	timeout    int
	jsonOutput bool

	llmbackendMaker llmbackend.LLMBackendMaker

	cmd *cobra.Command
}

func NewBackendCheckCmd(llmbackendMaker llmbackend.LLMBackendMaker) *BackendCheckCmd {
	checkCmd := &BackendCheckCmd{
		llmbackendMaker: llmbackendMaker,
	}

	checkCmd.cmd = &cobra.Command{
		Use:   "check <name>",
		Short: "Check an LLM backend connection",
		Long:  `Check an LLM backend by sending a simple ping request to verify the configuration is working.`,
		Args:  cobra.ExactArgs(1),
		RunE:  checkCmd.run,
	}

	checkCmd.cmd.Flags().IntVar(&checkCmd.timeout, "timeout", 30, "Request timeout in seconds")
	checkCmd.cmd.Flags().BoolVar(&checkCmd.jsonOutput, "json", false, "Output in JSON format")

	return checkCmd
}

func (b *BackendCheckCmd) run(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Load current config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get the backend
	backend := cfg.GetBackend(name)
	if backend == nil {
		return fmt.Errorf("backend '%s' not found", name)
	}

	if !b.jsonOutput {
		fmt.Printf("Checking backend '%s'...\n", name)
		fmt.Printf("  Type:     %s\n", backend.Type)
		fmt.Printf("  Model:    %s\n", backend.Model)
		if backend.BaseURL != "" {
			fmt.Printf("  Base URL: %s\n", backend.BaseURL)
		}
		fmt.Println()
	}

	// Create the backend client
	var llmBackend llmbackend.LLMBackend
	switch backend.Type {
	case "openai", "openrouter":
		llmBackend = b.llmbackendMaker(backend.APIKey, backend.BaseURL, backend.Model)
	default:
		return fmt.Errorf("unsupported backend type: %s", backend.Type)
	}

	// Set timeout
	timeout := time.Duration(b.timeout) * time.Second
	if backend.Timeout > 0 {
		timeout = time.Duration(backend.Timeout) * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Send a simple test message
	startTime := time.Now()
	messages := []llmbackend.Message{
		{Role: llmbackend.RoleUser, Content: "Hello! Please respond with 'OK' if you can read this."},
	}

	if !b.jsonOutput {
		fmt.Print("Sending test request... ")
	}
	response, usage, err := llmBackend.ChatCompletion(messages, ctx)
	elapsed := time.Since(startTime)

	if err != nil {
		if b.jsonOutput {
			jsonData, _ := json.MarshalIndent(map[string]interface{}{
				"success":       false,
				"error":         err.Error(),
				"backend":       name,
				"backend_type":  backend.Type,
				"model":         backend.Model,
				"response_time": elapsed.Milliseconds(),
			}, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			fmt.Println("❌ FAILED")
		}
		return fmt.Errorf("check failed: %w", err)
	}

	if b.jsonOutput {
		jsonData, err := json.MarshalIndent(map[string]interface{}{
			"success":       true,
			"backend":       name,
			"backend_type":  backend.Type,
			"model":         backend.Model,
			"base_url":      backend.BaseURL,
			"response_time": elapsed.Milliseconds(),
			"response":      response,
			"token_usage": map[string]interface{}{
				"input_tokens":  usage.InputTokens,
				"output_tokens": usage.OutputTokens,
				"total_tokens":  usage.TotalTokens,
			},
		}, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(jsonData))
	} else {
		fmt.Println("✅ SUCCESS")
		fmt.Println()
		fmt.Printf("Response time: %v\n", elapsed.Round(time.Millisecond))
		fmt.Printf("Response: %s\n", response)
		fmt.Println()
		fmt.Printf("Token usage:\n")
		fmt.Printf("  Input tokens:  %d\n", usage.InputTokens)
		fmt.Printf("  Output tokens: %d\n", usage.OutputTokens)
		fmt.Printf("  Total tokens:  %d\n", usage.TotalTokens)
	}

	return nil
}
