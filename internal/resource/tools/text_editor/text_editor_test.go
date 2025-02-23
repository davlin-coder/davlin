package texteditor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

func TestConcurrentEdits(t *testing.T) {
	te := newTextEditor()
	testPath := filepath.Join(t.TempDir(), "concurrent.txt")
	setupTestFile(t, testPath, "Initial content")

	const numGoroutines = 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			content := fmt.Sprintf("Content from goroutine %d", idx)
			_, err := te.write(nil, &Request{Path: testPath, FileText: content})
			if err != nil {
				t.Errorf("Concurrent write failed: %v", err)
			}
		}(i)
	}

	wg.Wait()

	// 验证历史记录是否正确保存
	if len(te.fileHistory[testPath]) != numGoroutines-1 {
		t.Errorf("Expected %d history entries, got %d", numGoroutines-1, len(te.fileHistory[testPath]))
	}
}

func TestWriteWithHistory(t *testing.T) {
	te := newTextEditor()
	testPath := filepath.Join(t.TempDir(), "test.txt")

	// 测试写入新文件（不应该有历史记录）
	t.Run("write new file without history", func(t *testing.T) {
		content := "Initial content"
		_, err := te.write(nil, &Request{Path: testPath, FileText: content})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(te.fileHistory[testPath]) != 0 {
			t.Error("Expected no history for new file")
		}
	})

	// 测试覆盖现有文件（应该有历史记录）
	t.Run("overwrite existing file with history", func(t *testing.T) {
		newContent := "Updated content"
		_, err := te.write(nil, &Request{Path: testPath, FileText: newContent})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(te.fileHistory[testPath]) != 1 {
			t.Error("Expected one history entry after overwrite")
		}
	})
}

func TestLargeFileHandling(t *testing.T) {
	te := newTextEditor()
	testPath := filepath.Join(t.TempDir(), "large.txt")

	// 创建一个超过大小限制的文件
	largeContent := strings.Repeat("a", 500*1024) // 500KB
	setupTestFile(t, testPath, largeContent)

	// 测试查看大文件
	t.Run("view large file", func(t *testing.T) {
		_, err := te.view(nil, &Request{Path: testPath})
		if err == nil {
			t.Error("Expected error for large file")
		}
	})
}
