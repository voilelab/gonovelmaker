package obsidian

import (
	"testing"
)

type TestFrontmatter struct {
	ID    string   `yaml:"id"`
	Title string   `yaml:"title"`
	Tags  []string `yaml:"tags"`
	Count int      `yaml:"count"`
}

func TestParseFrontmatter(t *testing.T) {
	t.Run("valid frontmatter with body", func(t *testing.T) {
		content := `---
id: test-id
title: Test Title
tags:
  - tag1
  - tag2
count: 42
---
This is the body content.
It has multiple lines.`

		fm, body, err := parseFrontmatter[TestFrontmatter]([]byte(content))
		if err != nil {
			t.Fatalf("parseFrontmatter failed: %v", err)
		}

		if fm.ID != "test-id" {
			t.Errorf("fm.ID = %s, want test-id", fm.ID)
		}
		if fm.Title != "Test Title" {
			t.Errorf("fm.Title = %s, want Test Title", fm.Title)
		}
		if len(fm.Tags) != 2 {
			t.Errorf("len(fm.Tags) = %d, want 2", len(fm.Tags))
		}
		if fm.Count != 42 {
			t.Errorf("fm.Count = %d, want 42", fm.Count)
		}

		expectedBody := "This is the body content.\nIt has multiple lines."
		if body != expectedBody {
			t.Errorf("body = %q, want %q", body, expectedBody)
		}
	})

	t.Run("no frontmatter", func(t *testing.T) {
		content := `This is just regular content without frontmatter.`

		_, _, err := parseFrontmatter[TestFrontmatter]([]byte(content))
		if err == nil {
			t.Error("expected error for content without frontmatter, got nil")
		}
	})

	t.Run("unclosed frontmatter", func(t *testing.T) {
		content := `---
id: test-id
title: Test Title
This is missing the closing ---`

		_, _, err := parseFrontmatter[TestFrontmatter]([]byte(content))
		if err == nil {
			t.Error("expected error for unclosed frontmatter, got nil")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		content := `---
id: test-id
title: Test Title
tags: [unclosed array
---
Body content`

		_, _, err := parseFrontmatter[TestFrontmatter]([]byte(content))
		if err == nil {
			t.Error("expected error for invalid YAML, got nil")
		}
	})

	t.Run("minimal frontmatter", func(t *testing.T) {
		content := `---
id: test
---
Body content only`

		fm, body, err := parseFrontmatter[TestFrontmatter]([]byte(content))
		if err != nil {
			t.Fatalf("parseFrontmatter failed: %v", err)
		}

		if fm.ID != "test" {
			t.Errorf("fm.ID = %s, want test", fm.ID)
		}
		if fm.Title != "" {
			t.Errorf("fm.Title = %s, want empty string", fm.Title)
		}
		if body != "Body content only" {
			t.Errorf("body = %q, want 'Body content only'", body)
		}
	})

	t.Run("frontmatter with special characters", func(t *testing.T) {
		content := `---
id: test-123
title: "Title with: colons and \"quotes\""
---
Body: with special characters!`

		fm, body, err := parseFrontmatter[TestFrontmatter]([]byte(content))
		if err != nil {
			t.Fatalf("parseFrontmatter failed: %v", err)
		}

		if fm.ID != "test-123" {
			t.Errorf("fm.ID = %s, want test-123", fm.ID)
		}
		if fm.Title != "Title with: colons and \"quotes\"" {
			t.Errorf("fm.Title = %s, want 'Title with: colons and \"quotes\"'", fm.Title)
		}
		if body != "Body: with special characters!" {
			t.Errorf("body = %q, want 'Body: with special characters!'", body)
		}
	})

	t.Run("multiline body content", func(t *testing.T) {
		content := `---
id: test
---
Line 1

Line 2

Line 3`

		_, body, err := parseFrontmatter[TestFrontmatter]([]byte(content))
		if err != nil {
			t.Fatalf("parseFrontmatter failed: %v", err)
		}

		expected := "Line 1\n\nLine 2\n\nLine 3"
		if body != expected {
			t.Errorf("body = %q, want %q", body, expected)
		}
	})

	t.Run("body with --- in content", func(t *testing.T) {
		content := `---
id: test
---
This body has --- in it
And more content`

		_, body, err := parseFrontmatter[TestFrontmatter]([]byte(content))
		if err != nil {
			t.Fatalf("parseFrontmatter failed: %v", err)
		}

		expected := "This body has --- in it\nAnd more content"
		if body != expected {
			t.Errorf("body = %q, want %q", body, expected)
		}
	})

	t.Run("empty body", func(t *testing.T) {
		content := `---
id: test
title: Test
---
`

		fm, body, err := parseFrontmatter[TestFrontmatter]([]byte(content))
		if err != nil {
			t.Fatalf("parseFrontmatter failed: %v", err)
		}

		if fm.ID != "test" {
			t.Errorf("fm.ID = %s, want test", fm.ID)
		}
		if body != "" {
			t.Errorf("body = %q, want empty string", body)
		}
	})

	t.Run("WorldbookFrontmatter", func(t *testing.T) {
		content := `---
id: magic-system
tags:
  - magic
  - fantasy
---
This is worldbook content about the magic system.`

		fm, body, err := parseFrontmatter[WorldbookFrontmatter]([]byte(content))
		if err != nil {
			t.Fatalf("parseFrontmatter failed: %v", err)
		}

		if fm.ID != "magic-system" {
			t.Errorf("fm.ID = %s, want magic-system", fm.ID)
		}
		if len(fm.Tags) != 2 {
			t.Errorf("len(fm.Tags) = %d, want 2", len(fm.Tags))
		}
		if body != "This is worldbook content about the magic system." {
			t.Errorf("unexpected body content: %s", body)
		}
	})

	t.Run("ChapterFrontmatter", func(t *testing.T) {
		content := `---
id: chapter-one
title: The Beginning
index: 1
---
Chapter content here.`

		fm, body, err := parseFrontmatter[ChapterFrontmatter]([]byte(content))
		if err != nil {
			t.Fatalf("parseFrontmatter failed: %v", err)
		}

		if fm.ID != "chapter-one" {
			t.Errorf("fm.ID = %s, want chapter-one", fm.ID)
		}
		if fm.Title != "The Beginning" {
			t.Errorf("fm.Title = %s, want The Beginning", fm.Title)
		}
		if fm.Index != 1 {
			t.Errorf("fm.Index = %d, want 1", fm.Index)
		}
		if body != "Chapter content here." {
			t.Errorf("unexpected body content: %s", body)
		}
	})

	t.Run("CharacterFrontmatter", func(t *testing.T) {
		content := `---
id: alice
name: Alice Smith
main: true
---
Alice is the protagonist of the story.`

		fm, body, err := parseFrontmatter[CharacterFrontmatter]([]byte(content))
		if err != nil {
			t.Fatalf("parseFrontmatter failed: %v", err)
		}

		if fm.ID != "alice" {
			t.Errorf("fm.ID = %s, want alice", fm.ID)
		}
		if fm.Name != "Alice Smith" {
			t.Errorf("fm.Name = %s, want Alice Smith", fm.Name)
		}
		if !fm.Main {
			t.Error("fm.Main should be true")
		}
		if body != "Alice is the protagonist of the story." {
			t.Errorf("unexpected body content: %s", body)
		}
	})

	t.Run("ProjectFrontmatter", func(t *testing.T) {
		content := `---
id: my-novel
name: My Novel
system_prompt: You are a creative writer.
system_prompt_char: You create detailed characters.
---
This is the world description for my novel.`

		fm, body, err := parseFrontmatter[ProjectFrontmatter]([]byte(content))
		if err != nil {
			t.Fatalf("parseFrontmatter failed: %v", err)
		}

		if fm.ID != "my-novel" {
			t.Errorf("fm.ID = %s, want my-novel", fm.ID)
		}
		if fm.Name != "My Novel" {
			t.Errorf("fm.Name = %s, want My Novel", fm.Name)
		}
		if fm.SystemPrompt != "You are a creative writer." {
			t.Errorf("fm.SystemPrompt = %s, want 'You are a creative writer.'", fm.SystemPrompt)
		}
		if fm.SystemPromptChar != "You create detailed characters." {
			t.Errorf("fm.SystemPromptChar = %s, want 'You create detailed characters.'", fm.SystemPromptChar)
		}
		if body != "This is the world description for my novel." {
			t.Errorf("unexpected body content: %s", body)
		}
	})
}
