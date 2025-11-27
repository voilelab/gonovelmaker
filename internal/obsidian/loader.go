package obsidian

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/voilelab/gonovelmaker/novelmaker"
)

// WorldbookFrontmatter represents the YAML frontmatter for worldbook entries
type WorldbookFrontmatter struct {
	ID   string   `yaml:"id"`
	Tags []string `yaml:"tags"`
}

// ChapterFrontmatter represents the YAML frontmatter for chapters
type ChapterFrontmatter struct {
	ID    string `yaml:"id"`
	Title string `yaml:"title"`
	Index int    `yaml:"index"`
}

// CharacterFrontmatter represents the YAML frontmatter for characters
type CharacterFrontmatter struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
	Main bool   `yaml:"main"`
}

// ProjectFrontmatter represents the YAML frontmatter for project config
type ProjectFrontmatter struct {
	ID               string `yaml:"id"`
	Name             string `yaml:"name"`
	SystemPrompt     string `yaml:"system_prompt"`
	SystemPromptChar string `yaml:"system_prompt_char"`
}

// loadProjectFromRoot loads the project configuration from Config/project.md using os.Root
func loadProjectFromRoot(root *os.Root) (*novelmaker.Project, error) {
	projectPath := filepath.Join("Config", "project.md")
	content, err := root.ReadFile(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read project file %s: %w", projectPath, err)
	}

	fm, body, err := parseFrontmatter[ProjectFrontmatter](content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project frontmatter: %w", err)
	}

	if fm.ID == "" {
		return nil, fmt.Errorf("project config missing 'id' field")
	}

	if fm.Name == "" {
		return nil, fmt.Errorf("project config missing 'name' field")
	}

	now := time.Now()
	return &novelmaker.Project{
		ID:               fm.ID,
		Name:             fm.Name,
		World:            body,
		SystemPrompt:     fm.SystemPrompt,
		SystemPromptChar: fm.SystemPromptChar,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// loadWorldbooksFromRoot scans World/*.md and loads all worldbook entries using os.Root
func loadWorldbooksFromRoot(root *os.Root) ([]novelmaker.Worldbook, error) {
	worldDir := "World"

	entries, err := fs.ReadDir(root.FS(), worldDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read worldbook directory %s: %w", worldDir, err)
	}

	var worldbooks []novelmaker.Worldbook

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(worldDir, entry.Name())
		wb, err := loadWorldbookFromRoot(root, filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load worldbook %s: %v\n", filePath, err)
			continue
		}
		worldbooks = append(worldbooks, *wb)
	}

	return worldbooks, nil
}

// loadChaptersFromRoot scans Story/*.md and loads all chapters using os.Root
func loadChaptersFromRoot(root *os.Root) ([]novelmaker.Chapter, error) {
	storyDir := "Story"

	entries, err := fs.ReadDir(root.FS(), storyDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read story directory %s: %w", storyDir, err)
	}

	var chapters []novelmaker.Chapter

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(storyDir, entry.Name())
		ch, err := loadChapterFromRoot(root, filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load chapter %s: %v\n", filePath, err)
			continue
		}
		chapters = append(chapters, *ch)
	}

	// Sort by index
	sort.Slice(chapters, func(i, j int) bool {
		return chapters[i].Index < chapters[j].Index
	})

	return chapters, nil
}

// loadCharactersFromRoot scans Character/*.md and loads all characters using os.Root
func loadCharactersFromRoot(root *os.Root) ([]novelmaker.Character, error) {
	charDir := "Character"

	// Check if Character directory exists
	entries, err := fs.ReadDir(root.FS(), charDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []novelmaker.Character{}, nil
		}
		return nil, fmt.Errorf("failed to read character directory %s: %w", charDir, err)
	}

	var characters []novelmaker.Character

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(charDir, entry.Name())
		char, err := loadCharacterFromRoot(root, filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load character %s: %v\n", filePath, err)
			continue
		}
		characters = append(characters, *char)
	}

	// Sort characters: main characters first, then by name
	sort.Slice(characters, func(i, j int) bool {
		if characters[i].Main != characters[j].Main {
			return characters[i].Main // main characters first
		}
		return characters[i].Name < characters[j].Name
	})

	return characters, nil
}

func loadWorldbookFromRoot(root *os.Root, path string) (*novelmaker.Worldbook, error) {
	content, err := root.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	fm, body, err := parseFrontmatter[WorldbookFrontmatter](content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter in %s: %w", path, err)
	}

	if fm.ID == "" {
		return nil, fmt.Errorf("missing or invalid 'id' in frontmatter of %s", path)
	}

	// Get file modification time
	info, _ := root.Stat(path)
	updatedAt := time.Now()
	if info != nil {
		updatedAt = info.ModTime()
	}

	return &novelmaker.Worldbook{
		ID:        fm.ID,
		Tags:      fm.Tags,
		Content:   body,
		UpdatedAt: updatedAt,
	}, nil
}

func loadChapterFromRoot(root *os.Root, path string) (*novelmaker.Chapter, error) {
	content, err := root.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	fm, body, err := parseFrontmatter[ChapterFrontmatter](content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter in %s: %w", path, err)
	}

	if fm.ID == "" {
		return nil, fmt.Errorf("missing or invalid 'id' in frontmatter of %s", path)
	}

	if fm.Title == "" {
		return nil, fmt.Errorf("missing or invalid 'title' in frontmatter of %s", path)
	}

	if fm.Index == 0 {
		return nil, fmt.Errorf("missing or invalid 'index' in frontmatter of %s", path)
	}

	// Get file modification time
	info, _ := root.Stat(path)
	updatedAt := time.Now()
	if info != nil {
		updatedAt = info.ModTime()
	}

	return &novelmaker.Chapter{
		ID:        fm.ID,
		Index:     fm.Index,
		Title:     fm.Title,
		Content:   body,
		UpdatedAt: updatedAt,
	}, nil
}

func loadCharacterFromRoot(root *os.Root, path string) (*novelmaker.Character, error) {
	content, err := root.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	fm, body, err := parseFrontmatter[CharacterFrontmatter](content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter in %s: %w", path, err)
	}

	if fm.ID == "" {
		return nil, fmt.Errorf("missing or invalid 'id' in frontmatter of %s", path)
	}

	if fm.Name == "" {
		return nil, fmt.Errorf("missing or invalid 'name' in frontmatter of %s", path)
	}

	// Get file modification time
	info, _ := root.Stat(path)
	updatedAt := time.Now()
	if info != nil {
		updatedAt = info.ModTime()
	}

	return &novelmaker.Character{
		ID:        fm.ID,
		Name:      fm.Name,
		Main:      fm.Main,
		Profile:   body,
		UpdatedAt: updatedAt,
	}, nil
}

// PrintJSON outputs project, worldbooks, and chapters as JSON
func PrintJSON(project *novelmaker.Project, worldbooks []novelmaker.Worldbook, chapters []novelmaker.Chapter) error {
	data := map[string]any{
		"project":    project,
		"worldbooks": worldbooks,
		"chapters":   chapters,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
