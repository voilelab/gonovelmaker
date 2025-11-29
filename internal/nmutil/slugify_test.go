package nmutil

import "testing"

func TestSlugify(t *testing.T) {
	t.Run("basic lowercase conversion", func(t *testing.T) {
		result := Slugify("Hello World")
		expected := "hello_world"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "Hello World", result, expected)
		}
	})

	t.Run("removes punctuation", func(t *testing.T) {
		result := Slugify("Hello: World? Yes!")
		expected := "hello_world_yes"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "Hello: World? Yes!", result, expected)
		}
	})

	t.Run("removes commas and periods", func(t *testing.T) {
		result := Slugify("Hello, World. How are you?")
		expected := "hello_world_how_are_you"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "Hello, World. How are you?", result, expected)
		}
	})

	t.Run("preserves numbers", func(t *testing.T) {
		result := Slugify("Chapter 123")
		expected := "chapter_123"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "Chapter 123", result, expected)
		}
	})

	t.Run("preserves underscores", func(t *testing.T) {
		result := Slugify("hello_world_test")
		expected := "hello_world_test"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "hello_world_test", result, expected)
		}
	})

	t.Run("removes special characters", func(t *testing.T) {
		result := Slugify("Hello@#$%^&*()World")
		expected := "helloworld"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "Hello@#$%^&*()World", result, expected)
		}
	})

	t.Run("preserves unicode characters (Chinese)", func(t *testing.T) {
		result := Slugify("你好 World")
		expected := "你好_world"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "你好 World", result, expected)
		}
	})

	t.Run("preserves unicode characters (Japanese)", func(t *testing.T) {
		result := Slugify("こんにちは World")
		expected := "こんにちは_world"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "こんにちは World", result, expected)
		}
	})

	t.Run("mixed unicode and ascii with punctuation", func(t *testing.T) {
		result := Slugify("Chapter 1: 序章")
		expected := "chapter_1_序章"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "Chapter 1: 序章", result, expected)
		}
	})

	t.Run("empty string", func(t *testing.T) {
		result := Slugify("")
		expected := ""
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "", result, expected)
		}
	})

	t.Run("only punctuation", func(t *testing.T) {
		result := Slugify("!?:,.")
		expected := ""
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "!?:,.", result, expected)
		}
	})

	t.Run("multiple spaces", func(t *testing.T) {
		result := Slugify("hello    world")
		expected := "hello____world"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "hello    world", result, expected)
		}
	})

	t.Run("leading and trailing spaces", func(t *testing.T) {
		result := Slugify("  hello world  ")
		expected := "__hello_world__"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "  hello world  ", result, expected)
		}
	})

	t.Run("already lowercase", func(t *testing.T) {
		result := Slugify("hello_world")
		expected := "hello_world"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "hello_world", result, expected)
		}
	})

	t.Run("mixed case with numbers", func(t *testing.T) {
		result := Slugify("Chapter1Section2Part3")
		expected := "chapter1section2part3"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "Chapter1Section2Part3", result, expected)
		}
	})

	t.Run("unicode emoji preserved", func(t *testing.T) {
		result := Slugify("Hello 😀 World")
		expected := "hello_😀_world"
		if result != expected {
			t.Errorf("Slugify(%q) = %q, want %q", "Hello 😀 World", result, expected)
		}
	})
}
