package anthropicadapter

import (
	"context"
	"encoding/json"
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
	params := anthropic.MessageNewParams{
		Model: model,
		MaxTokens: int64(req.MaxTokens),
		System: []anthropic.TextBlockParam{{Text: req.System}},
		Messages: toSDKMessages(req.Messages),
	}
	if len(req.Tools) > 0 {
		params.Tools = toSDKTools(req.Tools)
	}
	return params
}

func toSDKTools(tools []domain.ToolDefinition) []anthropic.ToolUnionParam {
	out := make([]anthropic.ToolUnionParam, len(tools))
	for i, t := range tools {
		out[i] = anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name: t.Name,
				Description: anthropic.String(t.Description),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: t.InputSchema["properties"],
				},
			},
		}
	}
	return out
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
			case domain.ToolResultBlock:
				out = append(out, anthropic.ContentBlockParamUnion{
					OfToolResult: &anthropic.ToolResultBlockParam{
					ToolUseID: v.ToolUseID,
					Content: []anthropic.ToolResultBlockParamContentUnion{
						{OfText: &anthropic.TextBlockParam{Text: v.Content}},
					},
				},
			})
			case domain.ToolUseBlock:
				out = append(out, anthropic.ContentBlockParamUnion{
				OfToolUse: &anthropic.ToolUseBlockParam{
					ID: v.ID,
					Name: v.Name,
					Input: v.Input,
				},
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
		case anthropic.StopReasonToolUse:
			return domain.StopReasonToolUse
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
			case anthropic.ToolUseBlock:
				var input map[string]any
				_ = json.Unmarshal([]byte(v.Input), &input)
				out = append(out, domain.ToolUseBlock{
					ID: v.ID,
					Name: v.Name,
					Input: input,
				})
		}
	}
	return out
}
