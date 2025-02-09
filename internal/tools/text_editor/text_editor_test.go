package texteditor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	t.Cleanup(func() { os.Remove(path) })
}

func TestViewCommand(t *testing.T) {
	te := newTextEditor()
	testPath := filepath.Join(t.TempDir(), "test.txt")
	setupTestFile(t, testPath, "Hello World")

	t.Run("view existing file", func(t *testing.T) {
		resp, err := te.view(nil, &Request{Path: testPath})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !strings.Contains(resp.FileText, "Hello World") {
			t.Errorf("Expected file content in response, got: %s", resp.FileText)
		}
	})

	t.Run("view non-existent file", func(t *testing.T) {
		_, err := te.view(nil, &Request{Path: "nonexistent.txt"})
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})
}

func TestWriteCommand(t *testing.T) {
	te := newTextEditor()
	testPath := filepath.Join(t.TempDir(), "test.txt")

	t.Run("write new file", func(t *testing.T) {
		content := "New content"
		resp, err := te.write(nil, &Request{Path: testPath, FileText: content})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if resp.FileText != content {
			t.Errorf("Expected content %q, got %q", content, resp.FileText)
		}

		// Verify file content
		data, err := os.ReadFile(testPath)
		if err != nil {
			t.Errorf("Failed to read created file: %v", err)
		}
		if string(data) != content {
			t.Errorf("File content mismatch, expected %q, got %q", content, string(data))
		}
	})
}

func TestStrReplaceCommand(t *testing.T) {
	te := newTextEditor()
	testPath := filepath.Join(t.TempDir(), "test.txt")
	original := "Hello Old World\nSecond Line\nThird Line"
	setupTestFile(t, testPath, original)

	t.Run("valid replacement", func(t *testing.T) {
		req := &Request{
			Path:   testPath,
			OldStr: "Old",
			NewStr: "New",
		}

		resp, err := te.strReplace(nil, req)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !strings.Contains(resp.Message, "has been edited") {
			t.Errorf("Unexpected success message: %s", resp.Message)
		}

		// Verify file content
		data, err := os.ReadFile(testPath)
		if err != nil {
			t.Errorf("Failed to read modified file: %v", err)
		}
		if !strings.Contains(string(data), "Hello New World") {
			t.Errorf("Replacement not successful, content: %s", string(data))
		}
	})

	t.Run("multiple matches", func(t *testing.T) {
		setupTestFile(t, testPath, "Old Old Old")
		_, err := te.strReplace(nil, &Request{Path: testPath, OldStr: "Old"})
		if err == nil {
			t.Error("Expected error for multiple matches")
		}
	})
}

func TestUndoEditCommand(t *testing.T) {
	te := newTextEditor()
	testPath := filepath.Join(t.TempDir(), "test.txt")
	original := "Original content"
	setupTestFile(t, testPath, original)

	// First make an edit
	_, err := te.strReplace(nil, &Request{
		Path:   testPath,
		OldStr: "Original",
		NewStr: "Modified",
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	t.Run("successful undo", func(t *testing.T) {
		resp, err := te.undoEdit(nil, &Request{Path: testPath})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !strings.Contains(resp.Message, "undone") {
			t.Errorf("Unexpected undo message: %s", resp.Message)
		}

		data, err := os.ReadFile(testPath)
		if err != nil {
			t.Errorf("Failed to read undone file: %v", err)
		}
		if string(data) != original {
			t.Errorf("Undo failed, expected %q, got %q", original, string(data))
		}
	})

	t.Run("undo without history", func(t *testing.T) {
		_, err := te.undoEdit(nil, &Request{Path: "nonexistent.txt"})
		if err == nil {
			t.Error("Expected error for undo without history")
		}
	})
}

func TestInvalidCommand(t *testing.T) {
	te := newTextEditor()
	_, err := te.editText(nil, &Request{Command: "invalid"})
	if err == nil {
		t.Error("Expected error for invalid command")
	}
}
