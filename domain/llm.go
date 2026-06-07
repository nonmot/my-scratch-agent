package domain

import (
	"context"
)

type StopReason string

const (
	StopReasonEndTurn StopReason = "end_turn"
	StopReasonMaxTokens StopReason = "max_tokens"
)

type LLMRequest struct {
	System string
	Messages []Message
	MaxTokens int
}

type LLMResponse struct {
	StopReason StopReason
	Blocks []Block
}

type LLMClient interface {
	Complete(ctx context.Context, req LLMRequest) (*LLMResponse, error)
}
