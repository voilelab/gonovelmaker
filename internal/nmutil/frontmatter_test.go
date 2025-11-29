package nmutil

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

		fm, body, err := ParseFrontmatter[TestFrontmatter]([]byte(content))
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

		_, _, err := ParseFrontmatter[TestFrontmatter]([]byte(content))
		if err == nil {
			t.Error("expected error for content without frontmatter, got nil")
		}
	})

	t.Run("unclosed frontmatter", func(t *testing.T) {
		content := `---
id: test-id
title: Test Title
This is missing the closing ---`

		_, _, err := ParseFrontmatter[TestFrontmatter]([]byte(content))
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

		_, _, err := ParseFrontmatter[TestFrontmatter]([]byte(content))
		if err == nil {
			t.Error("expected error for invalid YAML, got nil")
		}
	})

	t.Run("minimal frontmatter", func(t *testing.T) {
		content := `---
id: test
---
Body content only`

		fm, body, err := ParseFrontmatter[TestFrontmatter]([]byte(content))
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

		fm, body, err := ParseFrontmatter[TestFrontmatter]([]byte(content))
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

		_, body, err := ParseFrontmatter[TestFrontmatter]([]byte(content))
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

		_, body, err := ParseFrontmatter[TestFrontmatter]([]byte(content))
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

		fm, body, err := ParseFrontmatter[TestFrontmatter]([]byte(content))
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
}
