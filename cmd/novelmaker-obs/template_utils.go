package main

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/voilelab/gonovelmaker/internal/nmutil"
)

// ChapterPromptFrontmatter represents the YAML frontmatter for chapter prompts
type ChapterPromptFrontmatter struct {
	System string `yaml:"system"`
}

// templateFuncMap provides custom functions for templates
var templateFuncMap = template.FuncMap{
	"join": strings.Join,
}

// parseChapterTemplate parses the chapter_prompt.md content (with frontmatter) into a template
func parseChapterTemplate(content string) (*template.Template, error) {
	fm, body, err := nmutil.ParseFrontmatter[ChapterPromptFrontmatter]([]byte(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse chapter prompt frontmatter: %w", err)
	}

	// The body contains the actual template content
	tmpl, err := template.New("chapter").Funcs(templateFuncMap).Parse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chapter template: %w", err)
	}

	// Note: The system prompt from frontmatter can be accessed via fm.System if needed
	// For now, we're using the project's system prompt from project.md
	_ = fm

	return tmpl, nil
}
