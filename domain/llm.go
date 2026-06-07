package domain

import (
	"context"
)

type StopReason string

const (
	StopReasonEndTurn StopReason = "end_turn"
	StopReasonToolUse StopReason = "tool_use"
	StopReasonMaxTokens StopReason = "max_tokens"
)

type ToolDefinition struct {
	Name string
	Description string
	InputSchema map[string]any
}

type LLMRequest struct {
	System string
	Messages []Message
	MaxTokens int
	Tools []ToolDefinition
}

type LLMResponse struct {
	StopReason StopReason
	Blocks []Block
}

type LLMClient interface {
	Complete(ctx context.Context, req LLMRequest) (*LLMResponse, error)
}
