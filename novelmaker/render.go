package novelmaker

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/voilelab/gonovelmaker/internal/llmbackend"
)

type Renderer struct {
	llmBackend llmbackend.LLMBackend

	ChapterTemplate   *template.Template
	CharacterTemplate *template.Template
	timeout           time.Duration
}

func NewRenderer(llmBackend llmbackend.LLMBackend, chapterTempl, characterTempl *template.Template, timeout time.Duration) *Renderer {
	return &Renderer{
		llmBackend:        llmBackend,
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
	project *Project, worldbook []Worldbook, characters []Character, preChapters []Chapter, target string, prompt string) (string, llmbackend.UsageInfo, error) {

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
		return "", llmbackend.UsageInfo{}, fmt.Errorf("failed to execute chapter template: %w", err)
	}
	promptContent := buf.String()

	systemPrompt := project.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = "You are a helpful assistant that writes novels."
	}

	msgs := []llmbackend.Message{
		{Role: llmbackend.RoleSystem, Content: systemPrompt},
		{Role: llmbackend.RoleUser, Content: promptContent},
	}

	ctx := context.Background()
	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	chatCompletion, usage, err := r.llmBackend.ChatCompletion(msgs, ctx)

	if err != nil {
		return "", llmbackend.UsageInfo{}, err
	}

	return chatCompletion, usage, nil
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
	project *Project, worldbook []Worldbook, characters []Character, prompt string, name string) (string, llmbackend.UsageInfo, error) {

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
		return "", llmbackend.UsageInfo{}, fmt.Errorf("failed to execute character template: %w", err)
	}
	promptContent := buf.String()

	systemPrompt := project.SystemPromptChar
	if systemPrompt == "" {
		systemPrompt = "You are a helpful assistant that creates detailed character profiles for novels."
	}

	msgs := []llmbackend.Message{
		{Role: llmbackend.RoleSystem, Content: systemPrompt},
		{Role: llmbackend.RoleUser, Content: promptContent},
	}

	ctx := context.Background()
	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	chatCompletion, usage, err := r.llmBackend.ChatCompletion(msgs, ctx)

	if err != nil {
		return "", llmbackend.UsageInfo{}, err
	}

	return chatCompletion, usage, nil
}
