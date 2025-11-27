package novelmaker

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Project struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	World            string    `json:"world"`
	SystemPrompt     string    `json:"system_prompt"`
	SystemPromptChar string    `json:"system_prompt_char"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Worldbook struct {
	ID        string    `json:"id"`
	Tags      []string  `json:"tags"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Character struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Main      bool      `json:"main"`
	Profile   string    `json:"profile"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Chapter struct {
	ID        string    `json:"id"`
	Index     int       `json:"index"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Renderer struct {
	APIKey            string
	Model             string
	BaseURL           string
	ChapterTemplate   *template.Template
	CharacterTemplate *template.Template
}

func NewRenderer(apiKey, model, baseURL string, chapterTempl, characterTempl *template.Template) *Renderer {
	return &Renderer{
		APIKey:            apiKey,
		Model:             model,
		BaseURL:           baseURL,
		ChapterTemplate:   chapterTempl,
		CharacterTemplate: characterTempl,
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

	chatCompletion, err := client.Chat.Completions.New(context.TODO(),
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
	project *Project, worldbook []Worldbook, characters []Character, prompt string, name string) (profile string, extractedName string, err error) {

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
		return "", "", fmt.Errorf("failed to execute character template: %w", err)
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

	chatCompletion, err := client.Chat.Completions.New(context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: msgs,
			Model:    r.Model,
		})

	if err != nil {
		return "", "", err
	}

	profile = chatCompletion.Choices[0].Message.Content

	// Try to extract name from the response if not provided
	if name == "" {
		lines := strings.Split(profile, "\n")
		for i, line := range lines {
			line = strings.TrimSpace(line)
			// Look for "Name: XXX" pattern
			if strings.HasPrefix(strings.ToLower(line), "name:") {
				extractedName = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
				extractedName = strings.TrimSpace(strings.TrimPrefix(extractedName, "name:"))
				// Remove this line from the profile
				profile = strings.Join(append(lines[:i], lines[i+1:]...), "\n")
				break
			}
			// Look for "# Name" or "## Name" pattern in first few lines
			if i < 3 && (strings.HasPrefix(line, "# ") || strings.HasPrefix(line, "## ")) {
				extractedName = strings.TrimPrefix(line, "##")
				extractedName = strings.TrimPrefix(extractedName, "#")
				extractedName = strings.TrimSpace(extractedName)
				if extractedName != "" {
					break
				}
			}
		}
	}

	return profile, extractedName, nil
}
