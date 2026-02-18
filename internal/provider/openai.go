package provider

import (
	"context"
	"errors"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	client *openai.Client
	model  string
}

func NewOpenAI(apiKey, model, baseURL string) *OpenAIProvider {
	var client *openai.Client
	if baseURL != "" {
		cfg := openai.DefaultConfig(apiKey)
		cfg.BaseURL = baseURL
		client = openai.NewClientWithConfig(cfg)
	} else {
		client = openai.NewClient(apiKey)
	}
	return &OpenAIProvider{
		client: client,
		model:  model,
	}
}

func (p *OpenAIProvider) SetModel(model string) {
	p.model = model
}

func (p *OpenAIProvider) StreamChat(ctx context.Context, messages []ChatMessage, chunks chan<- string) error {
	defer close(chunks)

	msgs := make([]openai.ChatCompletionMessage, len(messages))
	for i, m := range messages {
		msgs[i] = openai.ChatCompletionMessage{
			Role:    string(m.Role),
			Content: m.Content,
		}
	}

	stream, err := p.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:    p.model,
		Messages: msgs,
		Stream:   true,
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if len(resp.Choices) > 0 && resp.Choices[0].Delta.Content != "" {
			select {
			case chunks <- resp.Choices[0].Delta.Content:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func (p *OpenAIProvider) ListModels(ctx context.Context) ([]string, error) {
	resp, err := p.client.ListModels(ctx)
	if err != nil {
		return nil, err
	}
	var models []string
	for _, m := range resp.Models {
		models = append(models, m.ID)
	}
	return models, nil
}
