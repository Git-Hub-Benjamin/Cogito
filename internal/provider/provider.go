package provider

import "context"

type Provider interface {
	StreamChat(ctx context.Context, messages []ChatMessage, chunks chan<- string) error
	ListModels(ctx context.Context) ([]string, error)
}
