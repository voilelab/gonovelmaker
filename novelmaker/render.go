package novelmaker

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/voilelab/gonovelmaker/internal/llmbackend"
)

type Renderer struct {
	llmBackend llmbackend.LLMBackend

	timeout time.Duration
}

func NewRenderer(llmBackend llmbackend.LLMBackend, timeout time.Duration) *Renderer {
	return &Renderer{
		llmBackend: llmBackend,
		timeout:    timeout,
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

// RenderChapter generates a new chapter using the specified OpenAI model
func (r *Renderer) RenderChapter(
	project *Project, chapterPrompt *ChapterPrompt, worldbook []Worldbook, characters []Character, preChapters []Chapter, target string, prompt string) (string, llmbackend.UsageInfo, error) {

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
	if err := chapterPrompt.AssistantTemplate.Execute(&buf, data); err != nil {
		return "", llmbackend.UsageInfo{}, fmt.Errorf("failed to execute chapter template: %w", err)
	}
	promptContent := buf.String()

	systemPrompt := chapterPrompt.System
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
	project *Project, characterPrompt *CharacterPrompt, worldbook []Worldbook, characters []Character, prompt string, name string) (string, llmbackend.UsageInfo, error) {

	data := CharacterPromptData{
		ProjectName: project.Name,
		World:       project.World,
		Worldbook:   worldbook,
		Characters:  characters,
		Prompt:      prompt,
		Name:        name,
	}

	var buf bytes.Buffer
	if err := characterPrompt.AssistantTemplate.Execute(&buf, data); err != nil {
		return "", llmbackend.UsageInfo{}, fmt.Errorf("failed to execute character template: %w", err)
	}
	promptContent := buf.String()

	systemPrompt := characterPrompt.System
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

// RewritePromptData holds the data for rendering rewrite prompts
type RewritePromptData struct {
	ContextBefore  string
	TargetSentence string
	ContextAfter   string
}

// RenderRewrite generates a rewritten text using the specified OpenAI model
func (r *Renderer) RenderRewrite(
	project *Project, rewritePrompt *RewritePrompt, contextBefore string, targetSentence string, contextAfter string) (string, llmbackend.UsageInfo, error) {

	data := RewritePromptData{
		ContextBefore:  contextBefore,
		TargetSentence: targetSentence,
		ContextAfter:   contextAfter,
	}

	var buf bytes.Buffer
	if err := rewritePrompt.AssistantTemplate.Execute(&buf, data); err != nil {
		return "", llmbackend.UsageInfo{}, fmt.Errorf("failed to execute rewrite template: %w", err)
	}
	promptContent := buf.String()

	systemPrompt := rewritePrompt.System
	if systemPrompt == "" {
		systemPrompt = "You are a helpful assistant that rewrites and improves text for novels."
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
