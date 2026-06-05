package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {

	agent := NewAgent("claude-opus-4.6", 10)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	answer, err := agent.Run(ctx, "Describe yourself")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Answer: ", answer)
}
