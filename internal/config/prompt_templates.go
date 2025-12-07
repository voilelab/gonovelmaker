package config

import (
	"strings"
	"text/template"
)

// PromptTemplates holds the templates for generating prompts
// Note: Character and chapter prompt templates are now loaded from the vault's Config/ folder
type PromptTemplates struct {
	// Deprecated: Character prompts are now loaded from vault's Config/character_prompt.md
	CharacterTemplate *template.Template
}

// templateFuncMap provides custom functions for templates
var templateFuncMap = template.FuncMap{
	"join": strings.Join,
}

// LoadPromptTemplates is deprecated and kept for backward compatibility
// Character and chapter prompts are now loaded directly from the vault's Config/ folder
// This function returns an empty PromptTemplates struct
func LoadPromptTemplates() (*PromptTemplates, error) {
	return &PromptTemplates{}, nil
}
