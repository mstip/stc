package ai

import (
	"context"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
)

func TextToSpeach(text string, filePath string) error {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	req := openai.CreateSpeechRequest{
		Model: openai.TTSModel1,
		Input: text,
		Voice: openai.VoiceNova,
	}

	resp, err := client.CreateSpeech(context.Background(), req)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(resp)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, b, 0644)
	if err != nil {
		return err
	}
	return nil
}
