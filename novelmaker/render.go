package novelmaker

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Renderer struct {
	APIKey            string
	Model             string
	BaseURL           string
	ChapterTemplate   *template.Template
	CharacterTemplate *template.Template
	timeout           time.Duration
}

func NewRenderer(apiKey, model, baseURL string, chapterTempl, characterTempl *template.Template, timeout time.Duration) *Renderer {
	return &Renderer{
		APIKey:            apiKey,
		Model:             model,
		BaseURL:           baseURL,
		ChapterTemplate:   chapterTempl,
		CharacterTemplate: characterTempl,
		timeout:           timeout,
	}
}

// ChapterPromptData holds the data for rendering chapter prompts
type ChapterPromptData struct {
	ProjectName string
	World       string
	Worldbook   []Worldbook
	Characters  []Character
	PreChapters []Chapter
	Title       string
	Prompt      string
}

// RenderPrompt generates a new chapter using the specified OpenAI model
func (r *Renderer) RenderPrompt(
	project *Project, worldbook []Worldbook, characters []Character, preChapters []Chapter, target string, prompt string) (string, error) {

	data := ChapterPromptData{
		ProjectName: project.Name,
		World:       project.World,
		Worldbook:   worldbook,
		Characters:  characters,
		PreChapters: preChapters,
		Title:       target,
		Prompt:      prompt,
	}

	var buf bytes.Buffer
	if err := r.ChapterTemplate.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute chapter template: %w", err)
	}
	promptContent := buf.String()

	opts := []option.RequestOption{option.WithAPIKey(r.APIKey)}
	if r.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(r.BaseURL))
	}
	client := openai.NewClient(opts...)

	systemPrompt := project.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = "You are a helpful assistant that writes novels."
	}

	msgs := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
		openai.UserMessage(promptContent),
	}

	ctx := context.Background()
	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	chatCompletion, err := client.Chat.Completions.New(ctx,
		openai.ChatCompletionNewParams{
			Messages: msgs,
			Model:    r.Model,
		})

	if err != nil {
		return "", err
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

// CharacterPromptData holds the data for rendering character prompts
type CharacterPromptData struct {
	ProjectName string
	World       string
	Worldbook   []Worldbook
	Characters  []Character
	Prompt      string
	Name        string
}

// RenderCharacter generates a new character profile using the specified OpenAI model
func (r *Renderer) RenderCharacter(
	project *Project, worldbook []Worldbook, characters []Character, prompt string, name string) (string, error) {

	data := CharacterPromptData{
		ProjectName: project.Name,
		World:       project.World,
		Worldbook:   worldbook,
		Characters:  characters,
		Prompt:      prompt,
		Name:        name,
	}

	var buf bytes.Buffer
	if err := r.CharacterTemplate.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute character template: %w", err)
	}
	promptContent := buf.String()

	opts := []option.RequestOption{option.WithAPIKey(r.APIKey)}
	if r.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(r.BaseURL))
	}
	client := openai.NewClient(opts...)

	systemPrompt := project.SystemPromptChar
	if systemPrompt == "" {
		systemPrompt = "You are a helpful assistant that creates detailed character profiles for novels."
	}

	msgs := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
		openai.UserMessage(promptContent),
	}

	ctx := context.Background()
	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	chatCompletion, err := client.Chat.Completions.New(ctx,
		openai.ChatCompletionNewParams{
			Messages: msgs,
			Model:    r.Model,
		})

	if err != nil {
		return "", err
	}

	return chatCompletion.Choices[0].Message.Content, nil
}
