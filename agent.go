package main

import (
	"context"
	"fmt"
	"log"

	"github.com/my-scratch-agent/domain"
)

type Agent struct {
	client   domain.LLMClient
	maxSteps int
}

func NewAgent(client domain.LLMClient, maxSteps int) *Agent {
	a := &Agent{
		client: client,
		maxSteps: maxSteps,
	}
	return a
}

func (a *Agent) Run(ctx context.Context, userInput string) (string, error) {

	messages := []domain.Message {
		{Role: domain.RoleUser, Blocks: []domain.Block{
			domain.TextBlock{Text: userInput},
		}},
	}

	for step := 1; step <= a.maxSteps; step++ {
		log.Printf("--- Step %d ---", step)

		resp, err := a.client.Complete(ctx, domain.LLMRequest {
			System: "あなたはエージェントです。ユーザーの質問に答えてください。",
			Messages: messages,
			MaxTokens: 1024,
		})

		if err != nil {
			return "", fmt.Errorf("llm at step %d: %w", step, err)
		}

		messages = append(messages, domain.Message{
			Role: domain.RoleAssistant,
			Blocks: resp.Blocks,
		})

		if resp.StopReason == domain.StopReasonEndTurn{
			var textOut string
			for _, b := range resp.Blocks {
				if v, ok := b.(domain.TextBlock); ok {
					textOut += v.Text
				}
			}
			return textOut, nil
		}
	}

	return "", fmt.Errorf("max steps (%d) reached without end_turn", a.maxSteps)
}
