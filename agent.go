package main

import (
	"context"
)

type Agent struct {
	model string
	maxSteps int
}

func NewAgent(model string, maxSteps int) *Agent {
	a := &Agent{
		model: model,
		maxSteps: maxSteps,
	}
	return a
}

func (a *Agent) Run(ctx context.Context, userInput string) (string, error) {
	return "Hello, I'm my-scratch-agent", nil
}
