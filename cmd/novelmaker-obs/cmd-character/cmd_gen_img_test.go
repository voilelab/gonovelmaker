package cmdcharacter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/voilelab/gonovelmaker/cmd/novelmaker-obs/testutil"
	"github.com/voilelab/gonovelmaker/internal/llmbackend"
	"github.com/voilelab/gonovelmaker/internal/obsidian"
)

// mockImageServer creates a test HTTP server that serves a fake image
func mockImageServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a simple 1x1 PNG image (smallest valid PNG)
		pngData := []byte{
			0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
			0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
			0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
			0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
			0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
			0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
			0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
			0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
			0x42, 0x60, 0x82,
		}
		w.Header().Set("Content-Type", "image/png")
		w.Write(pngData)
	}))
}

// testBackendMaker creates a backend that returns URLs from the test server
func testBackendMaker(imageURL string) llmbackend.LLMBackendMaker {
	return func(_, _, _ string) llmbackend.LLMBackend {
		return &testBackend{imageURL: imageURL}
	}
}

// testBackend is a mock backend for testing
type testBackend struct {
	imageURL string
}

func (t *testBackend) ChatCompletion(messages []llmbackend.Message, ctx context.Context) (string, llmbackend.UsageInfo, error) {
	return "Test response", llmbackend.UsageInfo{}, nil
}

func (t *testBackend) GenerateImage(prompt string, ctx context.Context) (string, error) {
	return t.imageURL, nil
}

func (t *testBackend) ListModels(ctx context.Context) ([]string, error) {
	return []string{"test-model-1", "test-model-2"}, nil
}

func TestGenCharImgCmd_Run_Success(t *testing.T) {
	t.Run("generate character image with test backend", func(t *testing.T) {
		server := mockImageServer()
		defer server.Close()

		tmpDir := testutil.SetupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharImgCmd := newGenImgCmd(testBackendMaker(server.URL))
		genCharImgCmd.name = "Alice"

		err := genCharImgCmd.run(genCharImgCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char-img command failed: %v", err)
		}

		// Verify the image file was created
		expectedPath := filepath.Join(tmpDir, "Character", "alice.png")
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Fatalf("expected image file to exist at %s", expectedPath)
		}

		// Verify file has content
		info, err := os.Stat(expectedPath)
		if err != nil {
			t.Fatalf("failed to stat image file: %v", err)
		}
		if info.Size() == 0 {
			t.Error("image file should not be empty")
		}
	})

	t.Run("generate character image with custom prompt", func(t *testing.T) {
		server := mockImageServer()
		defer server.Close()

		tmpDir := testutil.SetupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharImgCmd := newGenImgCmd(testBackendMaker(server.URL))
		genCharImgCmd.name = "Bob"
		genCharImgCmd.prompt = "A detailed portrait of a wise wizard with a long white beard"

		err := genCharImgCmd.run(genCharImgCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char-img command failed: %v", err)
		}

		// Verify the image file was created
		expectedPath := filepath.Join(tmpDir, "Character", "bob.png")
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Fatalf("expected image file to exist at %s", expectedPath)
		}
	})

	t.Run("generate image with custom output directory", func(t *testing.T) {
		server := mockImageServer()
		defer server.Close()

		tmpDir := testutil.SetupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		customOutputDir := filepath.Join(tmpDir, "CustomImages")
		genCharImgCmd := newGenImgCmd(testBackendMaker(server.URL))
		genCharImgCmd.name = "Alice"
		genCharImgCmd.outputDir = customOutputDir

		err := genCharImgCmd.run(genCharImgCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char-img command failed: %v", err)
		}

		// Verify the image file was created in custom directory
		expectedPath := filepath.Join(customOutputDir, "alice.png")
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Fatalf("expected image file to exist at %s", expectedPath)
		}
	})

	t.Run("generate image with special characters in character name", func(t *testing.T) {
		server := mockImageServer()
		defer server.Close()

		tmpDir := testutil.SetupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Add a character with special characters in name
		vault, err := obsidian.NewVault(tmpDir)
		if err != nil {
			t.Fatalf("failed to open vault: %v", err)
		}

		specialChar := `---
name: Sir John O'Brien III
main: false
prompt: A noble knight
---
Sir John O'Brien III is a noble knight with a mysterious past.`
		testutil.WriteTestFile(t, tmpDir, "Character/sir-john-obrien-iii.md", specialChar)
		vault.Close()

		genCharImgCmd := newGenImgCmd(testBackendMaker(server.URL))
		genCharImgCmd.name = "Sir John O'Brien III"

		err = genCharImgCmd.run(genCharImgCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char-img command failed: %v", err)
		}

		// Verify the image file was created with slugified name
		expectedPath := filepath.Join(tmpDir, "Character", "sir_john_obrien_iii.png")
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Fatalf("expected image file to exist at %s", expectedPath)
		}
	})
}

func TestGenCharImgCmd_Run_JSONOutput(t *testing.T) {
	t.Run("json output format", func(t *testing.T) {
		server := mockImageServer()
		defer server.Close()

		tmpDir := testutil.SetupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharImgCmd := newGenImgCmd(testBackendMaker(server.URL))
		genCharImgCmd.name = "Alice"
		genCharImgCmd.json = true

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := genCharImgCmd.run(genCharImgCmd.cmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		if err != nil {
			t.Fatalf("gen-char-img command failed: %v", err)
		}

		// Parse JSON output
		var buf []byte
		buf, _ = io.ReadAll(r)
		output := string(buf)

		var data map[string]any
		if err := json.Unmarshal([]byte(output), &data); err != nil {
			t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, output)
		}

		// Verify JSON contains required fields
		filepath, ok := data["filepath"].(string)
		if !ok {
			t.Fatal("JSON output should contain 'filepath' field")
		}

		if filepath == "" {
			t.Error("filepath should not be empty")
		}

		imageURL, ok := data["image_url"].(string)
		if !ok {
			t.Fatal("JSON output should contain 'image_url' field")
		}

		if imageURL == "" {
			t.Error("image_url should not be empty")
		}

		// Verify the file exists
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			t.Errorf("file should exist at path: %s", filepath)
		}
	})
}

func TestGenCharImgCmd_Run_ErrorCases(t *testing.T) {
	t.Run("error when character not found", func(t *testing.T) {
		tmpDir := testutil.SetupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharImgCmd := newGenImgCmd(llmbackend.MakeDummy)
		genCharImgCmd.name = "NonExistentCharacter"

		err := genCharImgCmd.run(genCharImgCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when character not found, got nil")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("error should mention 'not found', got: %v", err)
		}
	})

	t.Run("error when vault not properly initialized", func(t *testing.T) {
		server := mockImageServer()
		defer server.Close()

		tmpDir := testutil.CreateTestVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharImgCmd := newGenImgCmd(testBackendMaker(server.URL))
		genCharImgCmd.name = "TestChar"

		err := genCharImgCmd.run(genCharImgCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when vault not properly initialized, got nil")
		}

		// Should fail on either project load or character load
		if !strings.Contains(err.Error(), "failed to load") && !strings.Contains(err.Error(), "not found") {
			t.Errorf("error should mention 'failed to load' or 'not found', got: %v", err)
		}
	})

	t.Run("error when config cannot be loaded", func(t *testing.T) {
		server := mockImageServer()
		defer server.Close()

		tmpDir := testutil.CreateTestVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)

		// Change to a directory that doesn't have a config
		os.Chdir(tmpDir)

		genCharImgCmd := newGenImgCmd(testBackendMaker(server.URL))
		genCharImgCmd.name = "TestChar"

		err := genCharImgCmd.run(genCharImgCmd.cmd, []string{})
		if err == nil {
			t.Error("expected error when config cannot be loaded")
		}
	})
}

func TestGenCharImgCmd_Run_ConfigOverrides(t *testing.T) {
	t.Run("override image model via flag", func(t *testing.T) {
		server := mockImageServer()
		defer server.Close()

		tmpDir := testutil.SetupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharImgCmd := newGenImgCmd(testBackendMaker(server.URL))
		genCharImgCmd.name = "Alice"
		genCharImgCmd.imageModel = "dall-e-2"

		err := genCharImgCmd.run(genCharImgCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char-img command failed: %v", err)
		}

		// Verify the image file was created
		expectedPath := filepath.Join(tmpDir, "Character", "alice.png")
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Fatalf("expected image file to exist at %s", expectedPath)
		}
	})

	t.Run("override timeout via flag", func(t *testing.T) {
		server := mockImageServer()
		defer server.Close()

		tmpDir := testutil.SetupCompleteVault(t)
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		genCharImgCmd := newGenImgCmd(testBackendMaker(server.URL))
		genCharImgCmd.name = "Bob"
		genCharImgCmd.timeout = 120

		err := genCharImgCmd.run(genCharImgCmd.cmd, []string{})
		if err != nil {
			t.Fatalf("gen-char-img command failed: %v", err)
		}

		// Verify the image file was created
		expectedPath := filepath.Join(tmpDir, "Character", "bob.png")
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Fatalf("expected image file to exist at %s", expectedPath)
		}
	})
}
