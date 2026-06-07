package anthropicadapter

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/my-scratch-agent/domain"
)

type Client struct {
	sdk anthropic.Client
	model anthropic.Model
}

func NewClient(apiKey string, model anthropic.Model) *Client {
	return &Client{
		sdk: anthropic.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL("https://api.anthropic.com")),
		model: model,
	}
}

func (c *Client) Complete(
	ctx context.Context,
	req domain.LLMRequest,
) (*domain.LLMResponse, error) {
	params := toSDKParams(c.model, req)

	msg, err := c.sdk.Messages.New(ctx, params)
	if err != nil {
		return nil, err
	}

	return fromSDKMessage(msg), nil
}

// domain → SDK

func toSDKParams(model anthropic.Model, req domain.LLMRequest) anthropic.MessageNewParams {
	return anthropic.MessageNewParams{
		Model: model,
		MaxTokens: int64(req.MaxTokens),
		System: []anthropic.TextBlockParam{{Text: req.System}},
		Messages: toSDKMessages(req.Messages),
	}
}

func toSDKMessages(msgs []domain.Message) []anthropic.MessageParam {
	out := make([]anthropic.MessageParam, len(msgs))
	for i, m := range msgs {
		out[i] = anthropic.MessageParam {
			Role: toSDKRole(m.Role),
			Content: toSDKBlocks(m.Blocks),
		}
	}
	return out
}

func toSDKRole(r domain.Role) anthropic.MessageParamRole {
	if r == domain.RoleUser {
		return anthropic.MessageParamRoleUser
	}
	return anthropic.MessageParamRoleAssistant
}

func toSDKBlocks(blocks []domain.Block) []anthropic.ContentBlockParamUnion {
	out := make([]anthropic.ContentBlockParamUnion, 0, len(blocks))
	for _, b := range blocks {
		switch v := b.(type) {
			case domain.TextBlock:
				out = append(out, anthropic.ContentBlockParamUnion{
					OfText: &anthropic.TextBlockParam{Text: v.Text},
			})
			default:
				panic(fmt.Sprintf("unknown block type: %T", v))
		}
	}
	return out
}

// SDK → domain

func fromSDKMessage(msg *anthropic.Message) *domain.LLMResponse {
	return &domain.LLMResponse{
		StopReason: fromSDKStopReason(msg.StopReason),
		Blocks: fromSDKBlocks(msg.Content),
	}
}


func fromSDKStopReason(r anthropic.StopReason) domain.StopReason {
	switch r {
		case anthropic.StopReasonMaxTokens:
			return domain.StopReasonMaxTokens
		default:
			return domain.StopReasonEndTurn
	}
}

func fromSDKBlocks(blocks []anthropic.ContentBlockUnion) []domain.Block {
	out := make([]domain.Block, 0, len(blocks))
	for _, b := range blocks {
		switch v := b.AsAny().(type) {
			case anthropic.TextBlock:
				out = append(out, domain.TextBlock{Text: v.Text})
		}
	}
	return out
}
