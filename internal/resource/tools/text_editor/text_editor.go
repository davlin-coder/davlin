package texteditor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

const description = `Perform text editing operations on files.

The ` + "`" + `command` + "`" + ` parameter specifies the operation to perform. Allowed options are:
    - ` + "`" + `view` + "`" + `: View the content of a file.
    - ` + "`" + `write` + "`" + `: Create or overwrite a file with the given content
    - ` + "`" + `str_replace` + "`" + `: Replace a string in a file with a new string.
    - ` + "`" + `undo_edit` + "`" + `: Undo the last edit made to a file.

To use the write command, you must specify ` + "`" + `file_text` + "`" + ` which will become the new content of the file. Be careful with
existing files! This is a full overwrite, so you must include everything - not just sections you are modifying.

To use the str_replace command, you must specify both ` + "`" + `old_str` + "`" + ` and ` + "`" + `new_str` + "`" + ` - the ` + "`" + `old_str` + "`" + ` needs to exactly match one
unique section of the original file, including any whitespace. Make sure to include enough context that the match is not
ambiguous. The entire original string will be replaced with ` + "`" + `new_str` + "`" + `.`

func New() (tool.InvokableTool, error) {
	te := newTextEditor()
	return utils.InferTool("text_editor", description, te.editText)
}

type Request struct { // enum
	Command  string `json:"command" jsonschema:"required,enum=view,enum=write,enum=str_replace,enum=undo_edit,description=The operation to perform. Allowed options are: 'view', 'write', 'str_replace', 'undo_edit'."`
	Path     string `json:"request.Path" jsonschema:"required,description=Absolute request.Path to file or directory, e.g. 'repo/file.py' or 'repo'."`
	OldStr   string `json:"old_str,omitempty" jsonschema:"description=The string to be replaced."`
	NewStr   string `json:"new_str,omitempty" jsonschema:"description=The string to replace the old string."`
	FileText string `json:"file_text,omitempty" jsonschema:"description=The new content of the file."`
}

type Response struct {
	FileText string `json:"file_text,omitempty" jsonschema:"description=The new content of the file."`
	Message  string `json:"message,omitempty" jsonschema:"description=Success message."`
}

type TextEditor struct {
	fileHistory map[string][]string
	mutex      sync.RWMutex
}

func newTextEditor() *TextEditor {
	return &TextEditor{
		fileHistory: make(map[string][]string),
	}
}

func (f *TextEditor) editText(ctx context.Context, request *Request) (*Response, error) {
	switch request.Command {
	case "view":
		return f.view(ctx, request)
	case "write":
		return f.write(ctx, request)
	case "str_replace":
		return f.strReplace(ctx, request)
	case "undo_edit":
		return f.undoEdit(ctx, request)
	default:
		return nil, fmt.Errorf("invalid command: %s", request.Command)
	}
}

func (f *TextEditor) view(_ context.Context, request *Request) (*Response, error) {
	info, err := os.Stat(request.Path)
	if err != nil {
		return nil, fmt.Errorf("The request.Path '%s' does not exist or is not a file.", request.Path)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("The request.Path '%s' is not a file.", request.Path)
	}

	const MAX_FILE_SIZE = 400 * 1024 // 400KB
	const MAX_CHAR_COUNT = 400000    // 400,000 characters

	// Check file size
	if info.Size() > MAX_FILE_SIZE {
		return nil, fmt.Errorf("File '%s' is too large (%.2fKB). Maximum size is 400KB to prevent memory issues.",
			request.Path, float64(info.Size())/1024.0)
	}

	// Read file content
	data, err := os.ReadFile(request.Path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %v", err)
	}
	content := string(data)

	// count the file content length
	charCount := len([]rune(content))
	if charCount > MAX_CHAR_COUNT {
		return nil, fmt.Errorf("File '%s' has too many characters (%d). Maximum character count is %d.",
			request.Path, charCount, MAX_CHAR_COUNT)
	}

	ext := filepath.Ext(request.Path)
	language := strings.TrimPrefix(ext, ".")

	return &Response{
		FileText: fmt.Sprintf("```%s\n%s\n```", language, content),
	}, nil
}

func (f *TextEditor) write(_ context.Context, request *Request) (*Response, error) {
	// 检查文件是否存在
	var oldContent []byte
	if info, err := os.Stat(request.Path); err == nil && !info.IsDir() {
		// 如果文件存在，读取旧内容用于历史记录
		oldContent, err = os.ReadFile(request.Path)
		if err != nil {
			return nil, fmt.Errorf("Failed to read existing file: %v", err)
		}
	}

	// 写入新内容
	if err := os.WriteFile(request.Path, []byte(request.FileText), 0644); err != nil {
		return nil, fmt.Errorf("Failed to write file: %v", err)
	}

	// 如果有旧内容，添加到历史记录
	if len(oldContent) > 0 {
		f.mutex.Lock()
		f.fileHistory[request.Path] = append(f.fileHistory[request.Path], string(oldContent))
		f.mutex.Unlock()
	}

	return &Response{
		FileText: request.FileText,
		Message: fmt.Sprintf("File '%s' has been written successfully.", request.Path),
	}, nil
}
func (f *TextEditor) strReplace(_ context.Context, request *Request) (*Response, error) {
	if _, err := os.Stat(request.Path); err != nil {
		return nil, fmt.Errorf("File '%s' does not exist, you can write a new file with the `write` command", request.Path)
	}

	// Read file content
	data, err := os.ReadFile(request.Path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %v", err)
	}
	content := string(data)

	// Check occurrence count of oldStr
	count := strings.Count(content, request.OldStr)
	if count > 1 {
		return nil, fmt.Errorf("'old_str' must appear exactly once in the file, but it appears multiple times")
	}
	if count == 0 {
		return nil, fmt.Errorf("'old_str' must appear exactly once in the file, but it does not appear in the file. Make sure the string exactly matches existing file content, including whitespace!")
	}

	f.mutex.Lock()
	f.fileHistory[request.Path] = append(f.fileHistory[request.Path], content)
	f.mutex.Unlock()

	// Replace content and write back to file (only first occurrence)
	newContent := strings.Replace(content, request.OldStr, request.NewStr, 1)
	if err := os.WriteFile(request.Path, []byte(newContent), 0644); err != nil {
		return nil, fmt.Errorf("Failed to write file: %v", err)
	}

	// Detect language based on file extension
	ext := filepath.Ext(request.Path)
	language := strings.TrimPrefix(ext, ".")

	// Show context snippet of modified section
	const SNIPPET_LINES = 4

	// Calculate line number before replacement
	idx := strings.Index(content, request.OldStr)
	if idx == -1 {
		// This should never happen
		return nil, fmt.Errorf("Unexpected error: oldStr not found")
	}
	prefix := content[:idx]
	replacementLine := strings.Count(prefix, "\n")

	// Calculate start and end line numbers for snippet
	startLine := replacementLine - SNIPPET_LINES
	if startLine < 0 {
		startLine = 0
	}
	newStrLineCount := strings.Count(request.NewStr, "\n")
	endLine := replacementLine + SNIPPET_LINES + newStrLineCount

	// Split all lines and extract the relevant ones
	lines := strings.Split(newContent, "\n")
	if endLine >= len(lines) {
		endLine = len(lines) - 1
	}
	snippetLines := lines[startLine : endLine+1]
	snippet := strings.Join(snippetLines, "\n")
	output := fmt.Sprintf("```%s\n%s\n```", language, snippet)

	successMessage := fmt.Sprintf(`The file %s has been edited, and the section now reads:
%s
Review the changes above for errors. Undo and edit the file again if necessary!`, request.Path, output)

	return &Response{
		Message: successMessage,
	}, nil
}

func (f *TextEditor) undoEdit(_ context.Context, request *Request) (*Response, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	history, ok := f.fileHistory[request.Path]
	if !ok || len(history) == 0 {
		return nil, fmt.Errorf("No edit history found for file '%s'", request.Path)
	}

	lastEdit := history[len(history)-1]
	if err := os.WriteFile(request.Path, []byte(lastEdit), 0644); err != nil {
		return nil, fmt.Errorf("Failed to write file: %v", err)
	}

	f.fileHistory[request.Path] = history[:len(history)-1]

	return &Response{
		Message: fmt.Sprintf("The last edit to file '%s' has been undone.", request.Path),
	}, nil
}
