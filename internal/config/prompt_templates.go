package config

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/character_prompt.tmpl
var defaultCharacterPrompt string

// PromptTemplates holds the templates for generating prompts
type PromptTemplates struct {
	CharacterTemplate *template.Template
}

// templateFuncMap provides custom functions for templates
var templateFuncMap = template.FuncMap{
	"join": strings.Join,
}

// LoadPromptTemplates loads prompt templates from config directory or uses defaults
// If custom templates don't exist, it automatically creates them from defaults
// Note: Chapter prompt template is now loaded from the vault's Config/chapter_prompt.md
func LoadPromptTemplates() (*PromptTemplates, error) {
	pt := &PromptTemplates{}

	// Get home directory for custom templates
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Can't access home directory, use defaults only
		pt.CharacterTemplate, err = template.New("character").Funcs(templateFuncMap).Parse(defaultCharacterPrompt)
		if err != nil {
			return nil, err
		}
		return pt, nil
	}

	customTemplatesDir := filepath.Join(homeDir, ".novelmaker", "templates")

	// Load character template
	characterTmplPath := filepath.Join(customTemplatesDir, "character_prompt.tmpl")
	if data, err := os.ReadFile(characterTmplPath); err == nil {
		// Custom template exists, use it
		pt.CharacterTemplate, err = template.New("character").Funcs(templateFuncMap).Parse(string(data))
		if err != nil {
			return nil, err
		}
	} else {
		// Create templates directory if it doesn't exist
		os.MkdirAll(customTemplatesDir, 0755)

		// Write default template
		os.WriteFile(characterTmplPath, []byte(defaultCharacterPrompt), 0644)

		// Use default template
		pt.CharacterTemplate, err = template.New("character").Funcs(templateFuncMap).Parse(defaultCharacterPrompt)
		if err != nil {
			return nil, err
		}
	}

	return pt, nil
}

// GetDefaultCharacterPrompt returns the default character prompt template content
func GetDefaultCharacterPrompt() string {
	return defaultCharacterPrompt
}
