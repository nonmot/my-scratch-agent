package tools

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/my-scratch-agent/domain"
)

type BashTool struct {}

func (t *BashTool) Name() string { return "bash" }

func (t *BashTool) Definition() domain.ToolDefinition {
	return domain.ToolDefinition{
		Name: "bash",
		Description: "シェルコマンドを実行し、標準出力と標準エラーを返す",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"command": map[string]any{
					"type": "string",
					"description": "実行するシェルコマンド",
				},
			},
			"required": []string{"command"},
		},
	}
}

func (t *BashTool) Execute(input map[string]any) (string, error) {
	command, _ := input["command"].(string)

	if !confirm(fmt.Sprintf("bash %s", command)) {
		return "user cancelled", nil
	}

	cmd := exec.Command("sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Run()
	return out.String(), nil
}

func confirm(prompt string) bool {
	fmt.Printf("\n⚠ %s\n実行しますか？ [y/N]", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	answer := strings.TrimSpace(scanner.Text())
	return strings.EqualFold(answer, "y")
}
