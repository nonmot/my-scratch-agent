package domain

import (
	"encoding/json"
)

type Message struct {
	Role string
	Content []ContentBlock
}

type ContentBlock struct {
	Type string // "text" / "tool_use" / "tool_result"

	// Type == "text"
	Text string

	// Type == "tool_use"
	ID string
	Name string
	Input json.RawMessage

	// Type == "tool_result"
	ToolUseID string
	Content string
	IsError bool
}
