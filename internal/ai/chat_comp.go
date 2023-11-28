package ai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
)

type ChatCompResult struct {
	Message     string
	Result      string
	TotalTokens int
}

func ChatComp(messages []openai.ChatCompletionMessage) (*ChatCompResult, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo1106,
			Messages: messages,
		},
	)

	if err != nil {
		return nil, err
	}

	message := ""
	for _, msg := range messages {
		message = message + fmt.Sprintf("ROLE:%s, CONTENT:%s\n", msg.Role, msg.Content)
	}

	return &ChatCompResult{
		Message:     message,
		Result:      resp.Choices[0].Message.Content,
		TotalTokens: resp.Usage.TotalTokens,
	}, nil
}

func ChatCompStream(messages []openai.ChatCompletionMessage) error {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo1106,
		Messages: messages,
		Stream:   true,
	}
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}
	defer stream.Close()

	fmt.Printf("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return nil
		}

		if err != nil {
			return err
		}

		fmt.Printf(response.Choices[0].Delta.Content)
	}
}
