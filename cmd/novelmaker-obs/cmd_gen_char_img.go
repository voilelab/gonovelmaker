package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/voilelab/gonovelmaker/internal/config"
	"github.com/voilelab/gonovelmaker/internal/nmutil"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
)

type GenCharImgCmd struct {
	json       bool
	prompt     string
	name       string
	apiKey     string
	baseURL    string
	imageModel string
	timeout    int
	outputDir  string

	llmBackendMaker LLMBackendMaker

	cmd *cobra.Command
}

func NewGenCharImgCmd(llmBackendMaker LLMBackendMaker) *GenCharImgCmd {
	g := &GenCharImgCmd{
		llmBackendMaker: llmBackendMaker,
	}
	g.cmd = &cobra.Command{
		Use:   "gen-char-img",
		Short: "Generate an image for a character using OpenAI DALL-E API",
		Long: `Generates an image for a character based on their profile using OpenAI's DALL-E API.
The image will be downloaded and saved to the Character directory.`,
		RunE: g.run,
	}

	g.cmd.Flags().BoolVarP(&g.json, "json", "j", false, "Output in JSON format")

	g.cmd.Flags().StringVarP(&g.name, "name", "n", "", "Name of the character (required)")
	g.cmd.Flags().StringVarP(&g.prompt, "prompt", "p", "", "Custom prompt for image generation (optional, uses character profile if not provided)")
	g.cmd.MarkFlagRequired("name")

	// Allow overriding config values per-command
	g.cmd.Flags().StringVar(&g.apiKey, "api-key", "", "OpenAI API key to override config (optional)")
	g.cmd.Flags().StringVar(&g.baseURL, "base-url", "", "OpenAI base URL to override config (optional)")
	g.cmd.Flags().StringVar(&g.imageModel, "image-model", "", "Image model to use (e.g., dall-e-3, dall-e-2), overrides config (optional)")
	g.cmd.Flags().IntVar(&g.timeout, "timeout", 60, "Timeout in seconds for the API request (optional)")
	g.cmd.Flags().StringVar(&g.outputDir, "output-dir", "", "Output directory for the image (defaults to Character/)")

	return g
}

func (g *GenCharImgCmd) run(cmd *cobra.Command, args []string) error {
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

	// Load existing characters to find the target character
	characters, err := vault.LoadCharacters()
	if err != nil {
		return fmt.Errorf("failed to load characters: %w", err)
	}

	// Find the character by name
	var targetCharacter *string
	for _, char := range characters {
		if char.Name == g.name {
			targetCharacter = &char.Profile
			break
		}
	}

	if targetCharacter == nil {
		return fmt.Errorf("character with name '%s' not found", g.name)
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
	effectiveImageModel := cfg.GetImageModelOrDefault()
	if g.imageModel != "" {
		effectiveImageModel = g.imageModel
	}
	effectiveBaseURL := backend.BaseURL
	if g.baseURL != "" {
		effectiveBaseURL = g.baseURL
	}
	effectiveTimeout := time.Duration(g.timeout) * time.Second

	// Determine image generation prompt
	imagePrompt := g.prompt
	if imagePrompt == "" {
		// Use character profile as prompt
		imagePrompt = fmt.Sprintf("Character portrait: %s", *targetCharacter)
	}

	if !g.json {
		fmt.Println("Generating character image with OpenAI DALL-E...")
		fmt.Printf("  Character: %s\n", g.name)
		fmt.Printf("  Model: %s\n", effectiveImageModel)
		fmt.Printf("  Prompt: %s\n", imagePrompt)
	}

	// Create LLM backend for image generation
	llmBackend := g.llmBackendMaker(
		effectiveAPIKey,
		effectiveBaseURL,
		effectiveImageModel,
	)

	// Generate image using OpenAI API
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if effectiveTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, effectiveTimeout)
		defer cancel()
	}

	imageURL, err := llmBackend.GenerateImage(imagePrompt, ctx)
	if err != nil {
		return fmt.Errorf("failed to generate image: %w", err)
	}

	// Download the image
	if !g.json {
		fmt.Printf("\n✓ Image generated successfully!\n")
		fmt.Printf("  URL: %s\n", imageURL)
		fmt.Println("  Downloading image...")
	}

	resp, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read image data: %w", err)
	}

	// Determine output directory
	outputDir := g.outputDir
	if outputDir == "" {
		outputDir = filepath.Join(cwd, "Character")
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save image to file
	filename := fmt.Sprintf("%s.png", nmutil.Slugify(g.name))
	outputPath := filepath.Join(outputDir, filename)

	if err := os.WriteFile(outputPath, imageData, 0644); err != nil {
		return fmt.Errorf("failed to write image file: %w", err)
	}

	if !g.json {
		fmt.Printf("\n✓ Image saved successfully!\n")
		fmt.Printf("  File: %s\n", outputPath)
	} else {
		output := map[string]any{
			"filepath":  outputPath,
			"image_url": imageURL,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(output); err != nil {
			return fmt.Errorf("failed to encode JSON output: %w", err)
		}
	}

	return nil
}
