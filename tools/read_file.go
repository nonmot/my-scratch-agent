package tools

import (
	"os"

	"github.com/my-scratch-agent/domain"
)

type ReadFileTool struct {}

func (t *ReadFileTool) Name() string { return "read_file" }

func (t *ReadFileTool) Definition() domain.ToolDefinition {
	return domain.ToolDefinition{
		Name: "read_file",
		Description: "指定したパスのファイルを読み込み、内容を返す",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{"type": "string", "description": "ファイルパス"},
			},
			"required": []string{"path"},
		},
	}
}

func (t *ReadFileTool) Execute(input map[string]any) (string, error) {
	path, _ := input["path"].(string)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
