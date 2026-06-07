package main

import (
	"context"
	"fmt"
	"log"

	"github.com/my-scratch-agent/domain"
	"github.com/my-scratch-agent/tools"
)

type Agent struct {
	client   domain.LLMClient
	tools map[string]tools.Tool
	maxSteps int
}

func NewAgent(client domain.LLMClient, maxSteps int, ts ...tools.Tool) *Agent {
	toolMap := make(map[string]tools.Tool, len(ts))
	for _, t := range ts {
		toolMap[t.Name()] = t
	}
	a := &Agent{
		client: client,
		tools: toolMap,
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

	toolDefs := make([]domain.ToolDefinition, 0, len(a.tools))
	for _, t := range a.tools {
		toolDefs = append(toolDefs, t.Definition())
	}

	for step := 1; step <= a.maxSteps; step++ {
		log.Printf("--- Step %d ---", step)

		resp, err := a.client.Complete(ctx, domain.LLMRequest {
			System: "あなたはエージェントです。ユーザーの質問に答えてください。",
			Messages: messages,
			MaxTokens: 1024,
			Tools: toolDefs,
		})

		if err != nil {
			return "", fmt.Errorf("llm at step %d: %w", step, err)
		}

		// Debug
		for _, b := range resp.Blocks {
			switch v := b.(type) {
				case domain.TextBlock:
					log.Printf("[LLM] %s", v.Text)
				case domain.ToolUseBlock:
					log.Printf("[ToolUse] %s %v", v.Name, v.Input)
			}
		}

		messages = append(messages, domain.Message{
			Role: domain.RoleAssistant,
			Blocks: resp.Blocks,
		})

		if resp.StopReason == domain.StopReasonToolUse {
			var resultBlocks []domain.Block
			for _, b := range resp.Blocks {
				tu, ok := b.(domain.ToolUseBlock)
				if !ok {
					continue
				}
				tool, exists := a.tools[tu.Name]
				if !exists {
					resultBlocks = append(resultBlocks, domain.ToolResultBlock{
						ToolUseID: tu.ID,
						Content: fmt.Sprintf("unknown tool: %s", tu.Name),
					})
					continue
				}
				result, err := tool.Execute(tu.Input)
				if err != nil {
					result = fmt.Sprintf("error: %v", err)
				}
				log.Printf("[ToolResult] %s", result)
				resultBlocks = append(resultBlocks, domain.ToolResultBlock{
					ToolUseID: tu.ID,
					Content: result,
				})
			}

			messages = append(messages, domain.Message{
				Role: domain.RoleUser,
				Blocks: resultBlocks,
			})
			continue
		}

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
