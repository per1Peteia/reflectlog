package agent

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

const absWorkingDir = "/Users/peripeteia/workspace/github.com/per1Peteia/code-editing-agent/"

type ToolDefinition struct {
	Name        string                                      `json:"name"`
	Description string                                      `json:"description"`  // tells the model what the tool does, when to (not) use it
	InputSchema anthropic.ToolInputSchemaParam              `json:"input_schema"` // describes as json schema what inputs this tool expects and which form they take
	Function    func(input json.RawMessage) (string, error) // executes tool with input the model sends the programmer, returns the result
}

// generate any JSON schema to send to the tool
func GenerateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	var v T

	schema := reflector.Reflect(v)

	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties,
	}
}

type PathNotPermittedError struct {
	Path string
}

func (p PathNotPermittedError) Error() string {
	return fmt.Sprintf("Error: cannot read %s as it is outside the permitted working directory", p.Path)
}

//
// READING FILES
//

var ReadFileDefinition = ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: ReadFileInputSchema,
	Function:    ReadFile,
}

type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

var ReadFileInputSchema = GenerateSchema[ReadFileInput]()

func ReadFile(input json.RawMessage) (string, error) {
	readFileInput := ReadFileInput{}
	err := json.Unmarshal(input, &readFileInput)
	if err != nil {
		panic(err)
	}

	// TODO make the working directory an environment variable dependent on what the cwd is when agent is executed, hardcoded for now
	absTargetPath, err := filepath.Abs(filepath.Join(absWorkingDir, readFileInput.Path))
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(absTargetPath, absWorkingDir) {
		return "", PathNotPermittedError{readFileInput.Path}
	}

	content, err := os.ReadFile(absTargetPath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

//
// LISTING FILES
//

var ListFilesDefinition = ToolDefinition{
	Name:        "list_files",
	Description: "List files and directories at a given path. If no path is provided, lists files in the current directory.",
	InputSchema: ListFilesInputSchema,
	Function:    ListFiles,
}

type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

var ListFilesInputSchema = GenerateSchema[ListFilesInput]()

// TODO this is very very basic traversal and will run into problems very quick
// ListFiles returns the list of files and directories in the current working directory
func ListFiles(input json.RawMessage) (string, error) {
	listFilesInput := ListFilesInput{}
	err := json.Unmarshal(input, &listFilesInput)
	if err != nil {
		panic(err)
	}

	dir := absWorkingDir // TODO i need to make this a configurable or automatically set variable that depends on the current os wd
	if listFilesInput.Path != "" {
		dir, err = filepath.Abs(filepath.Join(absWorkingDir, listFilesInput.Path))
		if err != nil {
			return "", err
		}
	}

	if !strings.HasPrefix(dir, absWorkingDir) {
		return "", PathNotPermittedError{listFilesInput.Path}
	}

	var files []string
	// walk the entire filepath from (and including) root (path var) and append to files accordingly
	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if relPath != "." {
			if info.IsDir() {
				files = append(files, relPath+"/")
			} else {
				files = append(files, relPath)
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

//
// EDITING FILES
//

var EditFileDefinition = ToolDefinition{
	Name: "edit_file",
	Description: `Make edits to a text file.

Replaces 'old_str' with 'new_str' in the given file. 'old_str' and 'new_str' MUST be different from each other.

If the file specified with path doesn't exist, it will be created.
`,
	InputSchema: EditFileInputSchema,
	Function:    EditFile,
}

type EditFileInput struct {
	Path   string `json:"path" jsonschema_description:"The path to the file"`
	OldStr string `json:"old_str" jsonschema_description:"Text to search for - must match exactly and must only have one match exactly"`
	NewStr string `json:"new_str" jsonschema_description:"Text to replace old_str with"`
}

var EditFileInputSchema = GenerateSchema[EditFileInput]()

// the magic happens in the model, it just inputs the surrogate new string into the tool function and you replace it programmatically
func EditFile(input json.RawMessage) (string, error) {
	editFileInput := EditFileInput{}
	err := json.Unmarshal(input, &editFileInput)
	if err != nil {
		return "", err
	}

	// check input validity
	if editFileInput.Path == "" || editFileInput.OldStr == editFileInput.NewStr {
		return "", fmt.Errorf("invalid input parameters")
	}

	absTargetPath, err := filepath.Abs(filepath.Join(absWorkingDir, editFileInput.Path))
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(absTargetPath, absWorkingDir) {
		return "", PathNotPermittedError{editFileInput.Path}
	}

	// read file or create it if it does not exist (new file case)
	content, err := os.ReadFile(editFileInput.Path)
	if err != nil {
		if os.IsNotExist(err) && editFileInput.OldStr == "" {
			return createNewFile(editFileInput.Path, editFileInput.NewStr)
		}
		return "", err
	}

	oldContent := string(content)
	newContent := strings.Replace(oldContent, editFileInput.OldStr, editFileInput.NewStr, -1)

	if oldContent == newContent && editFileInput.OldStr != "" {
		return "", fmt.Errorf("old_str not found in file")
	}

	err = os.WriteFile(editFileInput.Path, []byte(newContent), 0644)
	if err != nil {
		return "", err
	}

	return "OK", nil
}

func createNewFile(filePath, content string) (string, error) {
	dir := path.Dir(filePath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}

	return fmt.Sprintf("successfully created file %s", filePath), nil
}
