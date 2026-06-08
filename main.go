package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	anthropicadapter "github.com/my-scratch-agent/adapter/anthropic"

	"github.com/my-scratch-agent/memory"
	"github.com/my-scratch-agent/tools"
)

func main() {

	llm := anthropicadapter.NewClient(
		os.Getenv("ANTHROPIC_API_KEY"),
		anthropic.ModelClaudeSonnet4_5,
	)
	mem := memory.NewInMemory()
	agent := NewAgent(llm, mem, 10, &tools.ReadFileTool{}, &tools.BashTool{})

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			// Ctrl + D で終了
			fmt.Println("\nBye!")
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" {
			fmt.Println("\nBye!")
			break
		}
		answer, err := agent.Run(context.Background(), input)
		if err != nil {
			log.Printf("error: %v\n", err)
			continue
		}

		fmt.Printf("\n%s\n\n", answer)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
