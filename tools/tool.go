package tools

import "github.com/my-scratch-agent/domain"

type Tool interface {
	Name() string
	Definition() domain.ToolDefinition
	Execute(input map[string]any) (string, error)
}
