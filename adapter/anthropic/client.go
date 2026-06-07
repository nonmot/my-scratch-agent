package anthropicadapter

import (
	"context"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
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
	params anthropic.MessageNewParams,
) (*anthropic.Message, error) {
	return c.sdk.Messages.New(ctx, params)
}
