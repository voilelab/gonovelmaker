package obsidian

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/voilelab/gonovelmaker/novelmaker"
)

//go:embed all:init_template
var initTemplateFolder embed.FS

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

const (
	configDirName = "Config"
	worldDirName  = "World"
	charDirName   = "Character"
	storyDirName  = "Story"
)

type Vault struct {
	root *os.Root
}

func NewVault(root string) (*Vault, error) {
	rt, err := os.OpenRoot(root)
	if err != nil {
		return nil, err
	}
	return &Vault{root: rt}, nil
}

func (v *Vault) Close() error {
	return v.root.Close()
}

func (v *Vault) Initialize() error {
	// Check if Config/ already exists
	if _, err := v.root.Stat(configDirName); err == nil {
		return fmt.Errorf("Config/ directory already exists")
	}

	return copyEmbedFS(initTemplateFolder, "init_template", v.root)
}

func (v *Vault) LoadProject() (*novelmaker.Project, error) {
	projectPath := filepath.Join(configDirName, "project.md")
	content, err := v.root.ReadFile(projectPath)
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

func (v *Vault) LoadWorldbooks() ([]novelmaker.Worldbook, error) {
	entries, err := fs.ReadDir(v.root.FS(), worldDirName)
	if err != nil {
		return nil, fmt.Errorf("failed to read worldbook directory %s: %w", worldDirName, err)
	}

	var worldbooks []novelmaker.Worldbook

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(worldDirName, entry.Name())
		wb, err := v.loadWorldbookFromRoot(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load worldbook %s: %v\n", filePath, err)
			continue
		}
		worldbooks = append(worldbooks, *wb)
	}

	return worldbooks, nil
}

func (v *Vault) LoadCharacters() ([]novelmaker.Character, error) {
	// Check if Character directory exists
	entries, err := fs.ReadDir(v.root.FS(), charDirName)
	if err != nil {
		if os.IsNotExist(err) {
			return []novelmaker.Character{}, nil
		}
		return nil, fmt.Errorf("failed to read character directory %s: %w", charDirName, err)
	}

	var characters []novelmaker.Character

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(charDirName, entry.Name())
		char, err := v.loadCharacterFromRoot(filePath)
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

func (v *Vault) AddCharacter(c *novelmaker.Character) error {
	// Ensure Character directory exists on disk
	if err := v.root.MkdirAll(charDirName, 0755); err != nil {
		return fmt.Errorf("failed to create Character directory: %w", err)
	}

	// Prepare frontmatter similar to CLI behavior
	frontmatter := fmt.Sprintf(`---
id: %s
name: %s
main: %t
---
%s`, c.ID, c.Name, c.Main, c.Profile)

	// Destination file
	filename := fmt.Sprintf("%s.md", c.ID)
	dstPath := filepath.Join(charDirName, filename)

	if err := v.root.WriteFile(dstPath, []byte(frontmatter), 0644); err != nil {
		return fmt.Errorf("failed to write character file %s: %w", dstPath, err)
	}

	return nil

}

func (v *Vault) LoadChapters() ([]novelmaker.Chapter, error) {
	entries, err := fs.ReadDir(v.root.FS(), storyDirName)
	if err != nil {
		return nil, fmt.Errorf("failed to read story directory %s: %w", storyDirName, err)
	}

	var chapters []novelmaker.Chapter

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(storyDirName, entry.Name())
		ch, err := v.loadChapterFromRoot(filePath)
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

func (v *Vault) AddChapter(c *novelmaker.Chapter) error {
	// Ensure Story directory exists on disk
	if err := v.root.MkdirAll(storyDirName, 0755); err != nil {
		return fmt.Errorf("failed to create Story directory: %w", err)
	}

	// Prepare frontmatter similar to loader expectations
	frontmatter := fmt.Sprintf(`---
id: %s
title: %s
index: %d
---
%s`, c.ID, c.Title, c.Index, c.Content)

	// Destination file: include index for readability
	filename := fmt.Sprintf("%03d_%s.md", c.Index, c.ID)
	dstPath := filepath.Join(storyDirName, filename)

	if err := v.root.WriteFile(dstPath, []byte(frontmatter), 0644); err != nil {
		return fmt.Errorf("failed to write chapter file %s: %w", dstPath, err)
	}

	return nil
}

func (v *Vault) loadWorldbookFromRoot(path string) (*novelmaker.Worldbook, error) {
	content, err := v.root.ReadFile(path)
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
	info, _ := v.root.Stat(path)
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

func (v *Vault) loadChapterFromRoot(path string) (*novelmaker.Chapter, error) {
	content, err := v.root.ReadFile(path)
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
	info, _ := v.root.Stat(path)
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

func (v *Vault) loadCharacterFromRoot(path string) (*novelmaker.Character, error) {
	content, err := v.root.ReadFile(path)
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
	info, _ := v.root.Stat(path)
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
