package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	anthropicadapter "github.com/my-scratch-agent/adapter/anthropic"
)

func main() {

	llm := anthropicadapter.NewClient(
		os.Getenv("ANTHROPIC_API_KEY"),
		anthropic.ModelClaudeSonnet4_5,
	)
	agent := NewAgent(llm, "claude-opus-4-6", 10)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	answer, err := agent.Run(ctx, "あなたの使っているモデルは何？")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Answer: ", answer)
}
