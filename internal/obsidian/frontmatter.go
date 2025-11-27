package obsidian

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// parseFrontmatter extracts YAML frontmatter from markdown content
// Returns the parsed frontmatter as type T, the body content, and any error
func parseFrontmatter[T any](content []byte) (*T, string, error) {
	str := string(content)

	// Check if file starts with ---
	if !strings.HasPrefix(str, "---\n") {
		return nil, "", fmt.Errorf("no frontmatter found (must start with ---)")
	}

	// Find the closing ---
	rest := str[4:] // Skip the first "---\n"
	endIdx := strings.Index(rest, "\n---\n")
	if endIdx == -1 {
		return nil, "", fmt.Errorf("frontmatter not properly closed (missing closing ---)")
	}

	// Extract frontmatter and body
	frontmatterStr := rest[:endIdx]
	body := strings.TrimSpace(rest[endIdx+5:]) // Skip "\n---\n"

	// Parse YAML
	var fm T
	decoder := yaml.NewDecoder(bytes.NewReader([]byte(frontmatterStr)))
	if err := decoder.Decode(&fm); err != nil {
		return nil, "", fmt.Errorf("failed to parse YAML frontmatter: %w", err)
	}

	return &fm, body, nil
}
