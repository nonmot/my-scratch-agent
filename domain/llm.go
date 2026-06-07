package domain

import (
	"context"
	"github.com/anthropics/anthropic-sdk-go"
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
	Blocks []ContentBlock
}

type LLMClient interface {
	// Complete(ctx context.Context, req LLMRequest) (*LLMResponse, error)
	Complete(ctx context.Context, params anthropic.MessageNewParams) (*anthropic.Message, error)
}
