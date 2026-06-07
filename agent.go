package main

import (
	"context"
	"fmt"
	"log"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/my-scratch-agent/domain"
)

type AnthropicClient interface {
	CreateMessage(ctx context.Context, params anthropic.MessageNewParams) (*anthropic.Message, error)
}

type Agent struct {
	client   domain.LLMClient
	model    string
	maxSteps int
}

func NewAgent(client domain.LLMClient, model string, maxSteps int) *Agent {
	a := &Agent{
		client: client,
		model:    model,
		maxSteps: maxSteps,
	}
	return a
}

func (a *Agent) Run(ctx context.Context, userInput string) (string, error) {

	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(userInput)),
	}

	for step := 1; step <= a.maxSteps; step++ {
		log.Printf("--- Step %d ---", step)
		resp, err := a.client.Complete(ctx, anthropic.MessageNewParams{
			Model: a.model,
			MaxTokens: 1024,
			System: []anthropic.TextBlockParam{{Text: "あなたはエージェントです。ユーザーの質問に答えてください。"}},
			Messages: messages,
		})
		if err != nil {
			return "", fmt.Errorf("llm at step %d: %w", step, err)
		}

		messages = append(messages, resp.ToParam())

		if resp.StopReason == anthropic.StopReasonEndTurn {
			if len(resp.Content) > 0 {
				return resp.Content[0].Text, nil
			}
			return "", nil
		}

	}

	return "", fmt.Errorf("Max steps (%d) reached without end_turn", a.maxSteps)
}
