package ai

import (
	"context"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type CreateEmbeddResult struct {
	Text        string
	Embedd      []float64
	TotalTokens int
}

func ConvertF32toF64(f32 []float32) []float64 {
	f64 := make([]float64, len(f32))
	for i, v := range f32 {
		f64[i] = float64(v)
	}
	return f64
}

func CreateEmbedding(text string) (*CreateEmbeddResult, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: []string{strings.ToLower(text)},
		Model: openai.AdaEmbeddingV2,
	})

	if err != nil {
		return nil, err
	}
	return &CreateEmbeddResult{
		Text:        text,
		Embedd:      ConvertF32toF64(resp.Data[0].Embedding),
		TotalTokens: resp.Usage.TotalTokens,
	}, nil
}
